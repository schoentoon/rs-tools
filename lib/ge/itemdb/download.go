package itemdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

type itemWithMetadata struct {
	*ge.Item
	Category int  `json:"category"`
	Alpha    rune `json:"alpha"`
}

type items struct {
	Items []ge.Item `json:"items"`
}

type downloadTask struct {
	category int
	alpha    rune
}

type Progress struct {
	Tasks    int64
	Finished int64
}

func (p *Progress) Send(ch chan<- *Progress) {
	// send it non blocking as we have no clue about the underlaying channel
	if ch != nil {
		select {
		case ch <- p:
		default:
		}
	}
}

func (t *downloadTask) Process(client *http.Client, db *DB, progCh chan<- *Progress, progress *Progress) error {
	for page := 1; ; page++ {
		out := &items{}

		err := backoff.Retry(func() error {
			values := url.Values{}
			values.Set("category", fmt.Sprintf("%d", t.category))
			values.Set("alpha", string(t.alpha))
			values.Set("page", fmt.Sprintf("%d", page))
			url := fmt.Sprintf("https://secure.runescape.com/m=itemdb_rs/api/catalogue/items.json?%s", values.Encode())
			res, err := client.Get(url)
			if err != nil {
				return err
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				return fmt.Errorf("HTTP Status: %d %s", res.StatusCode, res.Status)
			}

			return json.NewDecoder(res.Body).Decode(&out)
		}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 3))
		if err != nil && err != io.EOF {
			return err
		}

		for _, item := range out.Items {
			err = db.add(&item)
			if err != nil {
				return err
			}
		}

		atomic.AddInt64(&progress.Finished, 1)
		progress.Send(progCh)

		// in case we have less than 12 items (that's the maximum amount returned)
		// or no items at all, this will be the last page and we cleanly exit
		if len(out.Items) < 12 || len(out.Items) == 0 {
			return nil
		}

		// looks like we have an extra task
		atomic.AddInt64(&progress.Tasks, 1)
		progress.Send(progCh)
	}
}

func Download(client *http.Client, concurrency int, progCh chan<- *Progress) (*DB, error) {
	db := New()
	err := db.update(client, concurrency, progCh)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) update(client *http.Client, concurrency int, progCh chan<- *Progress) (err error) {
	var wg sync.WaitGroup
	ch := make(chan downloadTask)
	errCh := make(chan error, 1)

	progress := &Progress{
		Tasks:    42 * 27,
		Finished: 0,
	}

	if progCh != nil {
		defer func() {
			close(progCh)
		}()
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, client *http.Client, errCh chan<- error) {
			defer wg.Done()

			for task := range ch {
				err := task.Process(client, db, progCh, progress)
				if err != nil {
					errCh <- err
					return
				}
			}
		}(&wg, client, errCh)
	}

	// we go off the assumption here that there are just 41 categories (and category 0, so actually 42)
	// source for this is https://runescape.wiki/w/Application_programming_interface#Catalogue
	for category := 0; category <= 41; category++ {
		for _, alpha := range "#abcdefghijklmnopqrstovwxyz" {
			select {
			case ch <- downloadTask{
				category: category,
				alpha:    alpha,
			}:
			case err = <-errCh:
				return err
			}
		}
	}

	close(ch)

	wg.Wait()
	return err
}

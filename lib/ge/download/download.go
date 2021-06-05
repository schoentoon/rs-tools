package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync/atomic"

	"github.com/cenkalti/backoff/v4"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
	"golang.org/x/net/publicsuffix"
)

type Progress struct {
	Tasks    int64
	Finished int64
}

func (p *Progress) send(ch chan<- *Progress) {
	// send it non blocking as we have no clue about the underlaying channel
	if ch != nil {
		select {
		case ch <- p:
		default:
		}
	}
}

func (m *meta) Download(client *http.Client, db ge.SearchItemInterface, w io.Writer, progCh chan<- *Progress) error {
	if client.Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return err
		}
		client.Jar = jar
	}

	out := make(chan *ge.Item)

	copy := &meta{
		Categories: m.Categories,
		Runedate:   m.Runedate,
		inserted:   make([]int64, 0, 1024),
	}

	progress := &Progress{
		Tasks:    0,
		Finished: 0,
	}

	for _, category := range copy.Categories {
		for _, count := range category.Count {
			if count > 0 {
				progress.Tasks++
			}
		}
	}

	if progCh != nil {
		defer close(progCh)
	}

	go func() {
		err := copy.download(client, db, out, progCh, progress)
		if err != nil {
			log.Println(err)
		}
	}()

	encoder := json.NewEncoder(w)

	for item := range out {
		err := encoder.Encode(item)
		if err != nil {
			return err
		}
	}

	return nil
}

type itemsAPI struct {
	Items []ge.Item `json:"items"`
}

func (m *meta) download(client *http.Client, db ge.SearchItemInterface, ch chan *ge.Item, progCh chan<- *Progress, progress *Progress) error {
	defer close(ch)

	for category := 0; category <= CATEGORY_COUNT; category++ {
		for alpha, count := range m.Categories[category].Count {
			err := m.downloadCategory(client, db, ch, category, count, alpha, progCh, progress)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *meta) alreadyInserted(id int64) bool {
	if m.inserted == nil {
		return false
	}

	for _, i := range m.inserted {
		if i == id {
			return true
		}
	}

	return false
}

func (m *meta) downloadCategory(client *http.Client, db ge.SearchItemInterface, ch chan *ge.Item, category, left int, alpha string, progCh chan<- *Progress, progress *Progress) error {
	if left == 0 {
		return nil
	}

	for page := 1; left > 0; page++ {
		var out *itemsAPI

		err := backoff.Retry(func() error {
			out = &itemsAPI{}
			values := url.Values{}
			values.Set("category", fmt.Sprintf("%d", category))
			values.Set("alpha", alpha)
			values.Set("page", fmt.Sprintf("%d", page))
			url := fmt.Sprintf("https://secure.runescape.com/m=itemdb_rs/api/catalogue/items.json?%s", values.Encode())
			resp, err := client.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
			}

			if resp.ContentLength == 0 {
				return errors.New("Empty content length??")
			}

			return json.NewDecoder(resp.Body).Decode(&out)
		}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
		if err != nil {
			return err
		}

		for _, item := range out.Items {
			left--

			if m.alreadyInserted(item.ID) {
				continue
			}

			exists, _ := db.GetItem(item.ID)
			if exists == nil {
				ch <- &item
			}

			m.inserted = append(m.inserted, item.ID)
		}

		atomic.AddInt64(&progress.Finished, 1)
		progress.send(progCh)

		if len(out.Items) < 12 || len(out.Items) == 0 {
			break
		}

		if left > 0 {
			// looks like we have an extra task
			atomic.AddInt64(&progress.Tasks, 1)
			progress.send(progCh)
		}
	}

	if left != 0 {
		log.Printf("Warning, still have %d items left in category: %d, alpha: %s\n", left, category, alpha)
	}

	return nil
}

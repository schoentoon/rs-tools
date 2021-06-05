package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/net/publicsuffix"
)

// this is the amount of categories currently in the database, this may need to be increased in the future
// in case there's ever a new skill or whatever, check this wiki page for it
// https://runescape.wiki/w/Application_programming_interface#items
const CATEGORY_COUNT = 41

type meta struct {
	Categories map[int]category `json:"categories"`
	Runedate   int              `json:"runedate"`

	inserted []int64
}

type category struct {
	Count map[string]int `json:"count"`
}

type categoryAPI struct {
	Alpha []struct {
		Letter string `json:"letter"`
		Items  int    `json:"items"`
	} `json:"alpha"`
}

type lastUpdateAPI struct {
	LastConfigUpdateRuneday int `json:"lastConfigUpdateRuneday"`
}

func DiffMetadataFromFile(client *http.Client, filename string) (*meta, error) {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return BuildMetadata(client)
	}
	defer f.Close()

	old, err := ReadMetadata(f)
	if err != nil {
		return nil, err
	}

	new, err := BuildMetadata(client)
	if err != nil {
		return nil, err
	}

	return new.Diff(old), nil
}

func fetchCategory(client *http.Client, c int) (*category, error) {
	var parsed categoryAPI

	err := backoff.Retry(func() error {
		resp, err := client.Get(fmt.Sprintf("https://services.runescape.com/m=itemdb_rs/api/catalogue/category.json?category=%d", c))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.ContentLength == 0 {
			return errors.New("Empty content length??")
		}

		return json.NewDecoder(resp.Body).Decode(&parsed)
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return nil, err
	}

	out := &category{
		Count: make(map[string]int),
	}

	for _, alpha := range parsed.Alpha {
		out.Count[alpha.Letter] = alpha.Items
	}

	return out, nil
}

// BuildMetadata this will build a new meta structure with live data
func BuildMetadata(client *http.Client) (*meta, error) {
	m := &meta{
		Categories: make(map[int]category),
	}

	if client.Jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			return nil, err
		}
		client.Jar = jar
	}

	// fetch all the categories and put them in our metadata
	for c := 0; c <= CATEGORY_COUNT; c++ {
		category, err := fetchCategory(client, c)
		if err != nil {
			return nil, err
		}

		m.Categories[c] = *category
	}

	// at the end we will check the api for the last time the database was updated and add it to our metadata
	var lastConfigUpdateRuneday lastUpdateAPI
	err := backoff.Retry(func() error {
		resp, err := client.Get("https://secure.runescape.com/m=itemdb_rs/api/info.json")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.ContentLength == 0 {
			return errors.New("Empty content length??")
		}

		return json.NewDecoder(resp.Body).Decode(&lastConfigUpdateRuneday)
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10))
	if err != nil {
		return nil, err
	}

	m.Runedate = lastConfigUpdateRuneday.LastConfigUpdateRuneday

	return m, nil
}

// ReadMetadata recreate a metadata structure from most likely a file
func ReadMetadata(r io.Reader) (*meta, error) {
	var out meta

	err := json.NewDecoder(r).Decode(&out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Serialize mostly used to save the metadata structure to disk or whatever
func (m *meta) Serialize(w io.Writer) error {
	return json.NewEncoder(w).Encode(m)
}

// Diff check the difference between 2 meta structures, makes it easier to figure out what should be updated
func (m *meta) Diff(m2 *meta) *meta {
	out := &meta{
		Categories: make(map[int]category),
	}

	for c := 0; c <= CATEGORY_COUNT; c++ {
		workmap := category{
			Count: make(map[string]int),
		}
		out.Categories[c] = workmap
		for alpha, count := range m.Categories[c].Count {
			workmap.Count[alpha] = count - m2.Categories[c].Count[alpha]
		}
	}

	out.Runedate = m.Runedate - m2.Runedate

	return out
}

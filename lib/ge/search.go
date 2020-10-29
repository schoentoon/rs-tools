package ge

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	ItemID int64
	Name   string
}

func SearchItems(query string, client *http.Client) ([]SearchResult, error) {
	resp, err := client.PostForm("https://secure.runescape.com/m=itemdb_rs/a=13/results", url.Values{"query": {query}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	out := []SearchResult{}
	doc.Find(".table-item-link").Each(func(i int, s *goquery.Selection) {
		href, h := s.Attr("href")
		title, t := s.Attr("title")
		if h && t {
			uri, err := url.Parse(href)
			if err != nil {
				return
			}
			id, err := strconv.ParseInt(uri.Query().Get("obj"), 10, 64)
			if err != nil {
				return
			}
			out = append(out, SearchResult{
				ItemID: id,
				Name:   title,
			})
		}
	})

	return out, nil
}

func (s *SearchResult) Graph(client *http.Client) (*Graph, error) {
	return PriceGraph(s.ItemID, client)
}

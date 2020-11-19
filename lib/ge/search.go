package ge

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	ItemID int64
	Name   string
}

func (g *Ge) SearchItems(query string) ([]SearchResult, error) {
	req, err := http.NewRequest("POST", "https://secure.runescape.com/m=itemdb_rs/a=13/results", strings.NewReader(url.Values{"query": {query}}.Encode()))
	if err != nil {
		return nil, err
	}
	if g.UserAgent != "" {
		req.Header.Set("User-Agent", g.UserAgent)
	}
	resp, err := g.Client.Do(req)
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

func (s *SearchResult) Graph(ge GeInterface) (*Graph, error) {
	return ge.PriceGraph(s.ItemID)
}

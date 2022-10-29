package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ViswaxOption struct {
	Rune string
	Cost int64
}

type ViswaxCombination struct {
	Date      time.Time
	Primary   ViswaxOption
	Secondary []ViswaxOption
}

func Viswax(client *http.Client) (out *ViswaxCombination, err error) {
	params := url.Values{
		"action":             {"parse"},
		"format":             {"json"},
		"text":               {"{{Rune Goldberg Machine/Current combinations}}"},
		"contentmodel":       {"wikitext"},
		"prop":               {"text"},
		"disablelimitreport": {"1"},
	}
	resp, err := client.Get(fmt.Sprintf("https://runescape.wiki/api.php?%s", params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	wrapper := struct {
		Parse struct {
			Text struct {
				Data string `json:"*"`
			} `json:"text"`
		} `json:"parse"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrapper.Parse.Text.Data))
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()

	out = &ViswaxCombination{
		Date:      time.Time{},
		Secondary: make([]ViswaxOption, 3),
	}

	secondaryPos := 0
	doc.Find("table.wikitable tr").Each(func(i int, s *goquery.Selection) {

		if out.Secondary[secondaryPos].Cost != 0 && secondaryPos < (len(out.Secondary)-1) {
			secondaryPos++
		}

		s.Find("td").Each(func(i int, s *goquery.Selection) {
			fmt.Printf("%s\n", s.Text())
			if attr, ok := s.Attr("rowspan"); ok && attr == "3" {
				if out.Date.IsZero() {
					out.Date, err = time.Parse("2 Jan 2006", s.Text())
					if err != nil {
						panic(err)
					}
				} else if out.Primary.Rune == "" {
					out.Primary.Rune = strings.Trim(s.Text(), " ")
				} else if out.Primary.Cost == 0 {
					raw := strings.ReplaceAll(s.Text(), ",", "")
					out.Primary.Cost, err = strconv.ParseInt(raw, 10, 64)
					if err != nil {
						panic(err)
					}
				}
			} else {
				if out.Secondary[secondaryPos].Rune == "" {
					out.Secondary[secondaryPos].Rune = strings.Trim(s.Text(), " ")
				} else if out.Secondary[secondaryPos].Cost == 0 {
					raw := strings.ReplaceAll(s.Text(), ",", "")
					out.Secondary[secondaryPos].Cost, err = strconv.ParseInt(raw, 10, 64)
					if err != nil {
						panic(err)
					}
				}
			}
		})
	})

	return
}

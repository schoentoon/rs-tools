package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AraxxorPathInfo struct {
	Minions     bool
	Acid        bool
	Darkness    bool
	Description string
	DaysLeft    int
}

var araxxorDaysLeftRegex = regexp.MustCompile(`Days until next rotation: (\d)`)

func AraxxorPath(client *http.Client) (*AraxxorPathInfo, error) {
	params := url.Values{
		"action":             {"parse"},
		"format":             {"json"},
		"text":               {"{{Araxxor rotation}}"},
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

	out := &AraxxorPathInfo{}

	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			out.Minions = s.HasClass("table-bg-green")
		case 1:
			out.Acid = s.HasClass("table-bg-green")
		case 2:
			out.Darkness = s.HasClass("table-bg-green")
		case 3:
			out.Description = s.Text()
		}
	})

	doc.Find("th").Each(func(i int, s *goquery.Selection) {
		if _, ok := s.Attr("colspan"); ok {
			results := araxxorDaysLeftRegex.FindStringSubmatch(s.Text())
			if len(results) == 2 {
				out.DaysLeft, err = strconv.Atoi(results[1])
			}
		}
	})

	return out, err
}

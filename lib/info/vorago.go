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

type VoragoRotationInfo struct {
	Rotation string
	DaysLeft int
}

var voragoDaysLeftRegex = regexp.MustCompile(`Next: (\d) days`)

func VoragoRotation(client *http.Client) (*VoragoRotationInfo, error) {
	params := url.Values{
		"action":             {"parse"},
		"format":             {"json"},
		"text":               {"{{Vorago rotations}}"},
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

	out := &VoragoRotationInfo{}

	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("table-bg-green") {
			out.Rotation = s.Text()
		} else {
			results := voragoDaysLeftRegex.FindStringSubmatch(s.Text())
			if len(results) == 2 {
				out.DaysLeft, err = strconv.Atoi(results[1])
			}
		}
	})

	return out, err
}

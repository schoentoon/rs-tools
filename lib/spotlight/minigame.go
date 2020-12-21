package spotlight

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type MinigamesInfo struct {
	Current  string
	Schedule map[time.Time]string

	scheduleOrder []time.Time
}

func (m *MinigamesInfo) Iterate(f func(when time.Time, minigame string) error) error {
	for _, t := range m.scheduleOrder {
		err := f(t, m.Schedule[t])
		if err != nil {
			return err
		}
	}

	return nil
}

func Minigames(client *http.Client) (*MinigamesInfo, error) {
	params := url.Values{
		"action":             {"parse"},
		"format":             {"json"},
		"text":               {"{{Minigame spotlight}}"},
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

	out := &MinigamesInfo{
		Schedule:      make(map[time.Time]string),
		scheduleOrder: make([]time.Time, 0),
	}

	firstMonth := time.Month(-1) // this is invalid on purpose so we can easily set it in the first loop
	currentYear := time.Now().Year()
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("table-bg-green") {
			out.Current = s.Find("a").Text()
		} else if s.HasClass("table-bg-grey") {
			minigame := s.Find("a").Text()
			whenRaw := s.Find("span").Text()

			when, err := time.Parse(`2 Jan`, string(whenRaw))
			if err != nil {
				return
			}

			if firstMonth == time.Month(-1) {
				firstMonth = when.Month()
			}

			when = when.AddDate(currentYear, 0, 0)

			// if the first month seen in this table is larger than the current month we crossed a year boundary
			// so first of all happy new year, second of all.. let's add an extra year
			if when.Month() < firstMonth {
				when = when.AddDate(1, 0, 0)
			}

			out.Schedule[when] = minigame
			out.scheduleOrder = append(out.scheduleOrder, when)
		}
	})

	sort.SliceStable(out.scheduleOrder, func(i, j int) bool { return out.scheduleOrder[i].Unix() < out.scheduleOrder[j].Unix() })

	return out, nil
}

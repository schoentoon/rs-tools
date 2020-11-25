//go:generate stringer -type=ElvenClan
package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ElvenClan int8

const (
	Invalid ElvenClan = iota
	Amlodd
	Cadarn
	Crwys
	Hefin
	Iorwerth
	Ithell
	Meilyr
	Trahaearn
)

func VoiceOfSeren(client *http.Client) ([]ElvenClan, error) {
	params := url.Values{
		"action":    {"query"},
		"format":    {"json"},
		"meta":      {"allmessages"},
		"ammesages": {"VoS"},
	}
	resp, err := client.Get(fmt.Sprintf("https://runescape.wiki/api.php?%s", params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	wrapper := struct {
		Query struct {
			AllMessages []struct {
				Name string `json:"name"`
				Data string `json:"*"`
			} `json:"allmessages"`
		} `json:"query"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}

	if len(wrapper.Query.AllMessages) != 1 {
		return nil, fmt.Errorf("Unexpected amount of parsed messages? %#+v", wrapper)
	}

	split := strings.Split(wrapper.Query.AllMessages[0].Data, ",")

	out := make([]ElvenClan, len(split)) // should always be 2 lol
	for clan := ElvenClan(0); clan <= Trahaearn; clan++ {
		for i, current := range split {
			if clan.String() == current {
				out[i] = clan
			}
		}
	}

	for _, clan := range out {
		if clan == Invalid {
			return nil, fmt.Errorf("Something went wrong parsing the VoS: %#+v", wrapper)
		}
	}

	return out, nil
}

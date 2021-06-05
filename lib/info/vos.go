//go:generate go run golang.org/x/tools/cmd/stringer -type=ElvenClan
package info

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func parseElvenClan(in string) (ElvenClan, error) {
	for clan := Amlodd; clan <= Trahaearn; clan++ {
		if clan.String() == in {
			return clan, nil
		}
	}

	return Invalid, fmt.Errorf("Invalid elven clan: %s", in)
}

func VoiceOfSeren(client *http.Client) ([]ElvenClan, error) {
	resp, err := client.Get("https://runescape.wiki/api.php?action=query&format=json&meta=allmessages&ammessages=VoS")
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
	for i := 0; i < 2; i++ {
		clan, err := parseElvenClan(split[i])
		if err != nil {
			return nil, err
		}
		out[i] = clan
	}

	return out, nil
}

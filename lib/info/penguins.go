package info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PenguinInfo struct {
	ActivePenguins []ActivePenguin `json:"Activepenguin"`
	Bear           []ActiveBear    `json:"Bear"`
}

type ActivePenguin struct {
	Name         string   `json:"name"`
	LastLocation string   `json:"last_location"`
	Disguise     string   `json:"disguise"`
	ConfinedTo   string   `json:"confined_to"`
	Warning      string   `json:"warning"`
	Requirements string   `json:"requirements"`
	LastSeen     LastSeen `json:"time_seen"`
}

type ActiveBear struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type LastSeen struct {
	time.Time
}

func (t *LastSeen) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(bytes.Trim(b, "\""), &timestamp)
	if err != nil {
		return err
	}
	t.Time = time.Unix(timestamp, 0)
	return nil
}

func Penguins(client *http.Client) (*PenguinInfo, error) {
	resp, err := client.Get("https://jq.world60pengs.com/rest/cache/actives.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	out := &PenguinInfo{}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

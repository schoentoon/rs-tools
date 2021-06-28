package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RunesphereStatus struct {
	Next    time.Time `json:"next"`
	Prev    time.Time `json:"previous"`
	Unknown bool      `json:"unknown"`
	Active  bool      `json:"active"`
}

func Runesphere(client *http.Client) (*RunesphereStatus, error) {
	resp, err := client.Get("https://www.redstormi.com/runeapps/runesphere/api.php")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	out := &RunesphereStatus{}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

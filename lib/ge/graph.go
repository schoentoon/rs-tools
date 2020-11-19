package ge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/*
{"daily":{"1588377600000":18095,"1588464000000":18201,"1588550400000":18316,"1588636800000":18939,"1588723200000":19295,"1588809600000":19906,
...
}}
*/

type Graph struct {
	ItemID int64
	Graph  map[time.Time]int
}

func (g *Ge) PriceGraph(itemID int64) (*Graph, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://secure.runescape.com/m=itemdb_rs/api/graph/%d.json", itemID), nil)
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

	wrapper := struct {
		Daily map[string]int `json:"daily"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}

	out := &Graph{ItemID: itemID, Graph: make(map[time.Time]int)}

	for epoch, price := range wrapper.Daily {
		parsed, err := strconv.ParseInt(epoch, 10, 64)
		if err != nil {
			return nil, err
		}
		out.Graph[time.Unix(parsed/1000, 0)] = price
	}

	return out, nil
}

func (g *Graph) LatestPrice() (when time.Time, price int) {
	for w, p := range g.Graph {
		if w.After(when) {
			when = w
			price = p
		}
	}

	return
}

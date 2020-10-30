package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

// {"app":"explore","dashboardId":0,"timezone":"browser","startTime":1604025095900,"interval":"2s","intervalMs":2000,"panelId":"Q-4fd19502-9294-45db-b54c-52afbfae11a0-0",
//"targets":[{"refId":"A","key":"Q-4fd19502-9294-45db-b54c-52afbfae11a0-0","data":"","target":245,"type":"timeseries"}],
//"range":{"from":"2020-10-30T01:31:35.900Z","to":"2020-10-30T02:31:35.900Z","raw":{"from":"now-1h","to":"now"}},
//"requestId":"explore","rangeRaw":{"from":"now-1h","to":"now"},"scopedVars":{"__interval":{"text":"2s","value":"2s"},
//"__interval_ms":{"text":2000,"value":2000}},"maxDataPoints":1860,"liveStreaming":false,"showingGraph":true,"showingTable":true,"adhocFilters":[]}

type queryRequest struct {
	Targets []target `json:"targets"`
}

type target struct {
	Target int64  `json:"target"`
	Type   string `json:"type"`
}

type queryResponse struct {
	Target     string    `json:"target"`
	Datapoints [][]int64 `json:"datapoints"`
}

func (s *server) query(w http.ResponseWriter, r *http.Request) {
	req := queryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// setup channels to receive responses
	outArray := make([]queryResponse, 0, len(req.Targets))
	ch := make(chan queryResponse, len(req.Targets))
	errCh := make(chan error, 1)

	// fire up a goroutine for every target, fetch the price graph and prepare it to the desired output
	for _, target := range req.Targets {
		go func(target int64) {
			graph, err := ge.PriceGraph(target, s.Client)
			if err != nil {
				errCh <- err
				return
			}
			name, err := s.itemIDToItem(target)
			if err != nil {
				errCh <- err
				return
			}
			out := queryResponse{
				Target:     name,
				Datapoints: make([][]int64, 0, len(graph.Graph)),
			}
			for when, price := range graph.Graph {
				// TODO filter on the actually specified times
				out.Datapoints = append(out.Datapoints, []int64{int64(price), when.Unix() * 1000})
			}

			sort.SliceStable(out.Datapoints, func(i, j int) bool { return out.Datapoints[i][1] < out.Datapoints[j][1] })

			ch <- out
		}(target.Target)
	}

	// on our main thread, listen for responses and append them to an end result
	for range req.Targets {
		select {
		case data := <-ch:
			outArray = append(outArray, data)
		case err := <-errCh:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// sort the endresults by name, to prevent the plugin from shuffling around results/graph colors upon refreshes
	sort.SliceStable(outArray, func(i, j int) bool { return outArray[i].Target < outArray[j].Target })

	if err := json.NewEncoder(w).Encode(outArray); err != nil {
		log.Printf("json enc: %+v", err)
	}
}

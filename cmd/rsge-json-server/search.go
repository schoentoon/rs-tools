package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

type searchRequest struct {
	Target string `json:"target"`
}

type searchResponse struct {
	Text  string `json:"text"`
	Value int64  `json:"value"`
}

func (s *server) search(w http.ResponseWriter, r *http.Request) {
	req := searchRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// gracefully return when not looking for anything
	if req.Target == "" {
		if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
			log.Printf("json enc: %+v", err)
		}
	}
	lower := strings.ToLower(req.Target)
	for id, name := range s.ItemCache {
		if strings.ToLower(name) == lower {
			if err := json.NewEncoder(w).Encode([]searchResponse{{Text: name, Value: id}}); err != nil {
				log.Printf("json enc: %+v", err)
			}
		}
	}

	results, err := ge.SearchItems(req.Target, s.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]searchResponse, len(results))
	for i, item := range results {
		out[i].Text = item.Name
		out[i].Value = item.ItemID
		s.ItemCache[item.ItemID] = item.Name
	}

	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Printf("json enc: %+v", err)
	}
}

func (s *server) itemIDToItem(itemID int64) string {
	out, ok := s.ItemCache[itemID]
	if ok {
		return out
	}

	res, err := ge.GetItem(itemID, s.Client)
	if err != nil {
		panic(err)
	}
	s.ItemCache[res.ID] = res.Name

	return res.Name
}

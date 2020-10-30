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

	// we first go look through our local item cache, perhaps we have an item with the same name
	s.itemCacheMutex.RLock()
	for id, name := range s.itemCache {
		if strings.ToLower(name) == lower {
			if err := json.NewEncoder(w).Encode([]searchResponse{{Text: name, Value: id}}); err != nil {
				log.Printf("json enc: %+v", err)
			}
			s.itemCacheMutex.RUnlock()
			return
		}
	}
	s.itemCacheMutex.RUnlock()

	// otherwise we go online and search for it
	results, err := ge.SearchItems(req.Target, s.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// afterwards we will of course fill our cache and return the results to the client
	s.itemCacheMutex.Lock()
	defer s.itemCacheMutex.Unlock()
	out := make([]searchResponse, len(results))
	for i, item := range results {
		out[i].Text = item.Name
		out[i].Value = item.ItemID
		s.itemCache[item.ItemID] = item.Name
	}

	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Printf("json enc: %+v", err)
	}
}

func (s *server) itemIDToItem(itemID int64) (string, error) {
	s.itemCacheMutex.RLock()
	out, ok := s.itemCache[itemID]
	s.itemCacheMutex.RUnlock() // we don't RUnlock this with a defer as otherwise it could deadlock later on
	if ok {
		return out, nil
	}

	res, err := ge.GetItem(itemID, s.Client)
	if err != nil {
		return "", err
	}
	s.itemCacheMutex.Lock()
	s.itemCache[res.ID] = res.Name
	s.itemCacheMutex.Unlock()

	return res.Name, nil
}

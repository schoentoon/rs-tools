package ge

import "net/http"

type GeInterface interface {
	PriceGraph(itemID int64) (*Graph, error)
	GetItem(itemID int64) (*Item, error)
	SearchItems(query string) ([]SearchResult, error)
}

type Ge struct {
	Client    *http.Client
	UserAgent string
}

package ge

import "net/http"

type GeInterface interface {
	PriceGraph(itemID int64) (*Graph, error)
}

type SearchItemInterface interface {
	SearchItems(query string) ([]Item, error)
	GetItem(itemID int64) (*Item, error)
}

type Ge struct {
	Client    *http.Client
	UserAgent string
}

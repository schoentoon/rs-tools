package ge

import (
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestPriceGraph(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 200, `{"daily":{"1588377600000":18095,"1588464000000":18201,"1588550400000":18316}}`
	})
	graph, err := PriceGraph(245, client)
	assert.Nil(t, err)
	assert.Equal(t, graph.ItemID, int64(245))
	assert.Len(t, graph.Graph, 3)
	assert.Equal(t, graph.Graph[time.Unix(1588377600, 0)], 18095)
	assert.Equal(t, graph.Graph[time.Unix(1588464000, 0)], 18201)
	assert.Equal(t, graph.Graph[time.Unix(1588550400, 0)], 18316)

	when, price := graph.LatestPrice()
	assert.Equal(t, when, time.Unix(1588550400, 0))
	assert.Equal(t, price, 18316)
}

func TestPriceGraphNotFound(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 404, `{}`
	})
	graph, err := PriceGraph(245, client)
	assert.NotNil(t, err)
	assert.Nil(t, graph)
	assert.EqualError(t, err, "HTTP Status: 404 ")
}

func TestPriceGraphInvalidJSON(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 200, `{`
	})
	graph, err := PriceGraph(245, client)
	assert.NotNil(t, err)
	assert.Nil(t, graph)
}

func TestPriceGraphInvalidTimestamps(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 200, `{"daily":{"0xdeadbeef":18095,"1588464000000":18201,"1588550400000":18316}}`
	})
	graph, err := PriceGraph(245, client)
	assert.NotNil(t, err)
	assert.Nil(t, graph)
}

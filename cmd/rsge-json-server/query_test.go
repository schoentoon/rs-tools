package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
	testge "gitlab.com/schoentoon/rs-tools/lib/ge/test"
)

func TestQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/query",
		bytes.NewBufferString(`{"targets":[{"target":245,"type":"timeseries"},{"target":28256,"type":"timeseries"}],
		"range":{"from":"2020-10-01T01:31:35.900Z","to":"2020-10-30T01:31:35.900Z"}}`))

	rec := httptest.NewRecorder()

	mock := testge.MockGe(t)
	defer mock.Close()
	mock.NextGraph <- &ge.Graph{
		ItemID: 245,
		Graph: map[time.Time]int{
			time.Unix(1602638294, 0): 18317,
			time.Unix(1602638293, 0): 18095,
			time.Unix(1588464000, 0): 18316,
		},
	}
	mock.NextGraph <- &ge.Graph{
		ItemID: 28256,
		Graph: map[time.Time]int{
			time.Unix(1602638294, 0): 18317,
			time.Unix(1602638293, 0): 18095,
			time.Unix(1588464000, 0): 18316,
		},
	}

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}
	s.itemCache[245] = "Wine of Zamorak"
	s.itemCache[28256] = "Wine of Saradomin"

	s.query(rec, req)
	assert.Equal(t, 200, rec.Result().StatusCode)

	out := []queryResponse{}
	err := json.NewDecoder(rec.Body).Decode(&out)
	assert.Nil(t, err)

	// check if the correct range is filtered out and whether it is sorted by time
	assert.Len(t, out, 2)
	assert.Equal(t, int64(1602638293000), out[0].Datapoints[0][1])
	assert.Equal(t, int64(1602638294000), out[0].Datapoints[1][1])
}

func TestQuery_EmptyQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/query",
		bytes.NewBufferString(`{"targets":[{"target":"","type":"timeseries"}]}`))

	rec := httptest.NewRecorder()

	mock := testge.MockGe(t)
	defer mock.Close()

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	s.query(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestQuery_FailPriceGraph(t *testing.T) {
	req := httptest.NewRequest("GET", "/query",
		bytes.NewBufferString(`{"targets":[{"target":245,"type":"timeseries"}]}`))

	rec := httptest.NewRecorder()

	mock := testge.MockGe(t)
	defer mock.Close()
	mock.NextError <- errors.New("PriceGraph API is down")

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	s.query(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
}

func TestQuery_FailGetItem(t *testing.T) {
	req := httptest.NewRequest("GET", "/query",
		bytes.NewBufferString(`{"targets":[{"target":245,"type":"timeseries"}]}`))

	rec := httptest.NewRecorder()

	mock := testge.MockGe(t)
	defer mock.Close()
	mock.SkipErrors = 1
	mock.NextGraph <- &ge.Graph{
		ItemID: 245,
		Graph: map[time.Time]int{
			time.Unix(1602638294, 0): 18317,
			time.Unix(1602638293, 0): 18095,
			time.Unix(1588464000, 0): 18316,
		},
	}
	mock.NextError <- errors.New("GetItem API is down")

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	s.query(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
}

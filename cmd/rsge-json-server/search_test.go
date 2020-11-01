package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
	testge "gitlab.com/schoentoon/rs-tools/lib/ge/test"
)

func TestSearch(t *testing.T) {
	mock := testge.MockGe(t)
	defer mock.Close()
	mock.NextSearchResult <- []ge.SearchResult{
		ge.SearchResult{
			ItemID: 245,
			Name:   "Wine of Zamorak",
		},
	}

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	// we do 2 calls, the second one should be cached
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/search", bytes.NewBufferString(`{"target":"Wine of Zamorak"}`))

		rec := httptest.NewRecorder()

		s.search(rec, req)
		assert.Equal(t, 200, rec.Result().StatusCode)

		out := []searchResponse{}
		err := json.NewDecoder(rec.Body).Decode(&out)
		assert.Nil(t, err)

		assert.Len(t, out, 1)
		assert.Equal(t, "Wine of Zamorak", out[0].Text)
		assert.Equal(t, int64(245), out[0].Value)
	}
}

func TestSearch_Empty(t *testing.T) {
	mock := testge.MockGe(t)
	defer mock.Close()

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	req := httptest.NewRequest("GET", "/search", bytes.NewBufferString(`{"target":""}`))

	rec := httptest.NewRecorder()

	s.search(rec, req)
	assert.Equal(t, 200, rec.Result().StatusCode)

	var out interface{}
	err := json.NewDecoder(rec.Body).Decode(&out)
	assert.Nil(t, err)
	assert.Len(t, out, 0)
}

func TestSearch_FailSearchItems(t *testing.T) {
	mock := testge.MockGe(t)
	defer mock.Close()
	mock.NextError <- errors.New("SearchItems API is down")

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	req := httptest.NewRequest("GET", "/search", bytes.NewBufferString(`{"target":"Wine of Zamorak"}`))

	rec := httptest.NewRecorder()

	s.search(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
}

func TestSearch_BadInput(t *testing.T) {
	mock := testge.MockGe(t)
	defer mock.Close()

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	req := httptest.NewRequest("GET", "/search", bytes.NewBufferString(`NOT JSON`))

	rec := httptest.NewRecorder()

	s.search(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestItemIDToItem(t *testing.T) {
	mock := testge.MockGe(t)
	defer mock.Close()
	mock.NextItem <- &ge.Item{
		ID:   245,
		Name: "Wine of Zamorak",
	}

	s := server{
		itemCache: make(map[int64]string),
		ge:        mock,
	}

	name, err := s.itemIDToItem(245)
	assert.Nil(t, err)
	assert.Equal(t, name, "Wine of Zamorak")

	assert.Len(t, s.itemCache, 1)

	_, exists := s.itemCache[245]
	assert.True(t, exists)

	name, err = s.itemIDToItem(245)
	assert.Nil(t, err)
	assert.Equal(t, name, "Wine of Zamorak")
}

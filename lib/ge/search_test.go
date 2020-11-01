package ge

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestSearch(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/search_wine.html")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})
	ge := Ge{Client: client}
	res, err := ge.SearchItems("wine")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 6)
}

func TestSearchNotFound(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 404, ""
	})
	ge := Ge{Client: client}
	res, err := ge.SearchItems("wine")
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestSearchInvalidItems(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/search_invalid.html")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})
	ge := Ge{Client: client}
	res, err := ge.SearchItems("wine")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 4)
}

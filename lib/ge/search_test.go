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

	assert.Equal(t, int64(1993), res[0].ID)
	assert.Equal(t, "Jug of wine", res[0].Name)
	assert.Equal(t, "https://secure.runescape.com/m=itemdb_rs/a=13/1603720907702_obj_sprite.gif?id=1993", res[0].Icon)
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

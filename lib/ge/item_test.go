package ge

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestGetItemValid(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 200, `{
			"item": {
			  "icon": "https://secure.runescape.com/m=itemdb_rs/1603720907702_obj_sprite.gif?id=245",
			  "icon_large": "https://secure.runescape.com/m=itemdb_rs/1603720907702_obj_big.gif?id=245",
			  "id": 245,
			  "type": "Herblore materials",
			  "typeIcon": "https://www.runescape.com/img/categories/Herblore materials",
			  "name": "Wine of Zamorak",
			  "description": "A jug full of Wine of Zamorak."
			}
		  }`
	})
	item, err := GetItem(245, client)
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.ID, int64(245))
	assert.Equal(t, item.Name, "Wine of Zamorak")
}

func TestGetItemNotFound(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 404, `{}`
	})
	item, err := GetItem(245, client)
	assert.NotNil(t, err)
	assert.Nil(t, item)
	assert.EqualError(t, err, "HTTP Status: 404 ")
}

func TestGetItemInvalidJSON(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 200, `{`
	})
	item, err := GetItem(245, client)
	assert.NotNil(t, err)
	assert.Nil(t, item)
}

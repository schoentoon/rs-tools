package info

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestViswax(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := os.ReadFile("testdata/viswax.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := Viswax(client)
	assert.Nil(t, err)

	assert.Equal(t, "Air rune", res.Primary.Rune)
	assert.Equal(t, int64(25000), res.Primary.Cost)
	assert.Equal(t, "Chaos rune", res.Secondary[0].Rune)
	assert.Equal(t, int64(71500), res.Secondary[0].Cost)
	assert.Equal(t, "Steam rune", res.Secondary[1].Rune)
	assert.Equal(t, int64(781000), res.Secondary[1].Cost)
	assert.Equal(t, "Blood rune", res.Secondary[2].Rune)
	assert.Equal(t, int64(371000), res.Secondary[2].Cost)
}

func TestViswaxOnline(t *testing.T) {
	lib.TestOnline(t)

	_, err := Viswax(http.DefaultClient)

	assert.Nil(t, err)
}

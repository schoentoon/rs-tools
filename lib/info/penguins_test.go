package info

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestPenguins(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := os.ReadFile("testdata/penguins.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := Penguins(client)

	assert.Nil(t, err)
	assert.Len(t, res.ActivePenguins, 12)
	assert.Len(t, res.Bear, 1)

	assert.Equal(t, "Lunar Isle", res.ActivePenguins[0].Name)
	assert.Equal(t, int64(1623426566), res.ActivePenguins[0].LastSeen.Time.Unix())
}

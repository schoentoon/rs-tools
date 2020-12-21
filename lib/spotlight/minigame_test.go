package spotlight

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestMinigames(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/minigame.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := Minigames(client)

	assert.Nil(t, err)

	assert.Equal(t, "The Great Orb Project", res.Current)
	assert.Len(t, res.Schedule, 26)
	assert.Equal(t, res.Schedule[time.Date(2020, time.December, 23, 0, 0, 0, 0, time.UTC)], "Flash Powder Factory", res.Schedule)
	assert.Equal(t, res.Schedule[time.Date(2020, time.December, 26, 0, 0, 0, 0, time.UTC)], "Castle Wars", res.Schedule)
	assert.Equal(t, res.Schedule[time.Date(2020, time.December, 29, 0, 0, 0, 0, time.UTC)], "Stealing Creation", res.Schedule)
	assert.Equal(t, res.Schedule[time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)], "Cabbage Facepunch Bonanza", res.Schedule)
	assert.Equal(t, res.Schedule[time.Date(2021, time.January, 4, 0, 0, 0, 0, time.UTC)], "Heist", res.Schedule)

	prev := time.Time{}

	assert.Nil(t, res.Iterate(func(w time.Time, m string) error {
		assert.Greater(t, w.Unix(), prev.Unix())

		prev = w
		return nil
	}))
}

func TestMinigamesOnline(t *testing.T) {
	lib.TestOnline(t)

	_, err := Minigames(http.DefaultClient)

	assert.Nil(t, err)
}

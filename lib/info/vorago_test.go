package info

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestVorago(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/vorago.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := VoragoRotation(client)

	assert.Nil(t, err)
	assert.Equal(t, "Team Split", res.Rotation)
	assert.Equal(t, 6, res.DaysLeft)
}

func TestVoragoOnline(t *testing.T) {
	lib.TestOnline(t)

	_, err := VoragoRotation(http.DefaultClient)

	assert.Nil(t, err)
}

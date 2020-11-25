package info

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestVoiceOfSeren(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/vos.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := VoiceOfSeren(client)
	assert.Nil(t, err)

	assert.Equal(t, res[0], Trahaearn)
	assert.Equal(t, res[1], Meilyr)
}

func TestVoiceOfSerenEmptyJSON(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 200, `{}`
	})

	res, err := VoiceOfSeren(client)
	assert.Error(t, err)
	assert.Nil(t, res)
}

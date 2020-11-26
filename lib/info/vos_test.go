package info

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestParseElvenClan(t *testing.T) {
	cases := []struct {
		in           string
		expectedClan ElvenClan
		expectError  bool
	}{
		{"Amlodd", Amlodd, false},
		{"Not a clan", Invalid, true},
	}

	for _, c := range cases {
		clan, err := parseElvenClan(c.in)
		if c.expectError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, c.expectedClan, clan)
		}
	}
}

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

func TestPossibleNextVoiceOfSeren(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/vos_history.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := PossibleNextVoiceOfSeren(client)
	assert.Nil(t, err)

	assert.Len(t, res, 4)

	assert.Contains(t, res, Amlodd)
	assert.Contains(t, res, Cadarn)
	assert.Contains(t, res, Meilyr)
	assert.Contains(t, res, Trahaearn)
}

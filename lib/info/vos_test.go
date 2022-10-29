package info

import (
	"net/http"
	"os"
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
		data, err := os.ReadFile("testdata/vos.json")
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

func TestVoiceOfSerenOnline(t *testing.T) {
	lib.TestOnline(t)

	_, err := VoiceOfSeren(http.DefaultClient)

	assert.Nil(t, err)
}

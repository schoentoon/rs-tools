package runemetrics

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestFetchProfile(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		assert.Contains(t, req.URL.String(), "Schoentoon")

		data, err := ioutil.ReadFile("testdata/profile.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := FetchProfile(client, "Schoentoon")

	assert.Nil(t, err)
	assert.Equal(t, "Schoentoon", res.Name)
	assert.Len(t, res.Activities, 20)

	arch, err := res.GetSkill(Archaeology)
	assert.Nil(t, err)
	if assert.NotNil(t, arch) {
		assert.Equal(t, arch.Level, 120)
		assert.Equal(t, arch.XP, 1263439689)
	}

	// to test the custom json marshallers we simply marshal the profile to a temporary buffer
	// to then read it from there again and compare it
	var buf bytes.Buffer
	if assert.Nil(t, json.NewEncoder(&buf).Encode(res)) {
		res1, err := ParseProfile(&buf)

		assert.Nil(t, err)
		assert.Equal(t, "Schoentoon", res1.Name)
		assert.Len(t, res1.Activities, 20)

		assert.Equal(t, res, res1)
	}
}

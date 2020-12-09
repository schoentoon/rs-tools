package runemetrics

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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

func TestFetchProfileOnline(t *testing.T) {
	lib.TestOnline(t)

	_, err := FetchProfile(http.DefaultClient, "Schoentoon")

	assert.Nil(t, err)
}

func TestNewAchievementsSince(t *testing.T) {
	f1, err1 := os.Open("testdata/profile.json")
	f2, err2 := os.Open("testdata/profile2.json")
	if !(assert.Nil(t, err1) || assert.Nil(t, err2)) {
		t.FailNow()
	}
	defer f1.Close()
	defer f2.Close()

	p1, err1 := ParseProfile(f1)
	p2, err2 := ParseProfile(f2)
	if !(assert.Nil(t, err1) || assert.Nil(t, err2)) {
		t.FailNow()
	}

	new := NewAchievementsSince(p1.Activities, p2.Activities)

	assert.Len(t, new, 6)

	var prev *Activity = nil
	for _, a := range new {
		assert.Greater(t, a.Date.Unix(), p1.Activities[0].Date.Unix())
		if prev != nil {
			assert.GreaterOrEqual(t, a.Date.Unix(), prev.Date.Unix())
		}
		prev = &a
	}

	newWrongOrder := NewAchievementsSince(p2.Activities, p1.Activities)

	assert.Equal(t, new, newWrongOrder)
}

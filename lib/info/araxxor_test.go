package info

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestAraxxor(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/araxxor.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := AraxxorPath(client)

	assert.Nil(t, err)
	assert.True(t, res.Minions)
	assert.False(t, res.Acid)
	assert.True(t, res.Darkness)
	assert.Equal(t, 3, res.DaysLeft)
	assert.Equal(t, "I died in the dark, covered in spiders.", res.Description)
}

package info

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestRunesphere(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/runesphere.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	res, err := Runesphere(client)

	assert.Nil(t, err)
	assert.False(t, res.Unknown)
	assert.False(t, res.Active)
	assert.Equal(t, time.Date(2021, time.June, 28, 20, 34, 21, 0, time.UTC), res.Next.UTC())
}

package info

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestViswax(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		data, err := ioutil.ReadFile("testdata/viswax.json")
		if err != nil {
			t.Fatal(err)
		}
		return 200, string(data)
	})

	_, err := Viswax(client)
	assert.Nil(t, err)
}

func TestViswaxOnline(t *testing.T) {
	//lib.TestOnline(t)

	res, err := Viswax(http.DefaultClient)

	fmt.Printf("%+v", res)

	assert.Nil(t, err)
}

package download

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestDownload(t *testing.T) {
	itemsReply, err := ioutil.ReadFile("testdata/items.api.json")
	if err != nil {
		t.Error(err)
	}

	meta := &meta{
		Categories: map[int]category{
			0: {
				Count: map[string]int{
					"a": 13,
				},
			},
		},
	}

	var count int32
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		atomic.AddInt32(&count, 1)
		if req.URL.Path == "/m=itemdb_rs/api/catalogue/items.json" {
			return 200, string(itemsReply)
		}
		t.Error("Unexpected http call", req)
		return 500, "" // just making the compiler happy
	})

	var buffer bytes.Buffer

	err = meta.Download(client, NewEmptyDB(), &buffer)
	assert.Nil(t, err)

	assert.Equal(t, int32(2), count)
}

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
					"a": 12,
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

	ch := make(chan *Progress, 1024*64) // we make this buffered channel large enough to just keep everything in memory
	err = meta.Download(client, NewEmptyDB(), &buffer, ch)
	assert.Nil(t, err)

	assert.Equal(t, int32(1), count)

	last := &Progress{
		Tasks:    0,
		Finished: 0,
	}
	for progress := range ch {
		assert.GreaterOrEqual(t, progress.Finished, last.Finished)
		last = progress
	}

	assert.Equal(t, last.Finished, last.Tasks)
}

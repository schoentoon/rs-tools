package itemdb

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestDownload(t *testing.T) {
	var count int32
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		switch atomic.AddInt32(&count, 1) {
		case 1:
			return 404, "lol we're down or something"
		case 2:
			return 200, "200, but no json I guess"
		}
		return 200, `{"items":[{"id":42,"name":"Test"}]}`
	})

	db, err := Download(client, 2)
	assert.Nil(t, err)
	assert.Len(t, db.idToItems, 1)

	// the 42 is the amount of categories here, so this should be updated when new categories are added
	// the +2 is because we give 2 false responses
	assert.Equal(t, int32(42*27)+2, count)
}

func TestDownloadFail(t *testing.T) {
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		return 404, "we're always gonna give a false response"
	})

	db, err := Download(client, 2)
	assert.Nil(t, db)
	assert.Error(t, err)
}

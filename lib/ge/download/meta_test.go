package download

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib"
)

func TestDiff(t *testing.T) {
	f1, err := os.Open("testdata/meta.0.json")
	assert.Nil(t, err)
	defer f1.Close()

	f2, err := os.Open("testdata/meta.changed.json")
	assert.Nil(t, err)
	defer f2.Close()

	m1, err := ReadMetadata(f1)
	assert.Nil(t, err)

	m2, err := ReadMetadata(f2)
	assert.Nil(t, err)

	diff := m2.Diff(m1)

	assert.Equal(t, 5, diff.Categories[0].Count["a"])
}

func TestBuildMetadata(t *testing.T) {
	metaReply, err := ioutil.ReadFile("testdata/meta.api.json")
	if err != nil {
		t.Error(err)
	}

	var count int32
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		atomic.AddInt32(&count, 1)
		if req.URL.Path == "/m=itemdb_rs/api/info.json" {
			return 200, `{"lastConfigUpdateRuneday":7034}`
		} else if req.URL.Path == "/m=itemdb_rs/api/catalogue/category.json" {
			return 200, string(metaReply)
		}
		t.Error("Unexpected http call", req)
		return 500, "" // just making the compiler happy
	})

	meta, err := BuildMetadata(client)
	assert.Nil(t, err)

	assert.Equal(t, 7034, meta.Runedate)

	// we count the amount of requests that were made, this should be the amount of categories
	// again we do plus 1 due to zero indexing.. and 1 extra call for the runedate
	assert.Equal(t, int32(CATEGORY_COUNT+1+1), count)
}

func TestDiffMetadataFromFile(t *testing.T) {
	metaReply, err := ioutil.ReadFile("testdata/meta.api.json")
	if err != nil {
		t.Error(err)
	}

	var count int32
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		atomic.AddInt32(&count, 1)
		if req.URL.Path == "/m=itemdb_rs/api/info.json" {
			if atomic.LoadInt32(&count) == 1 {
				// first request will be to test up to date, any later calls will be out of date
				return 200, `{"lastConfigUpdateRuneday":7034}`
			}
			return 200, `{"lastConfigUpdateRuneday":7040}`
		} else if req.URL.Path == "/m=itemdb_rs/api/catalogue/category.json" {
			return 200, string(metaReply)
		}
		t.Error("Unexpected http call", req)
		return 500, "" // just making the compiler happy
	})

	_, _, err = DiffMetadataFromFile(client, "testdata/meta.0.json")
	assert.Equal(t, ErrNotOutdated, err)

	meta, _, err := DiffMetadataFromFile(client, "testdata/meta.0.json")
	assert.Nil(t, err)

	// we really just check if this meta is filled in properly, +1 due to zero indexing
	assert.Len(t, meta.Categories, CATEGORY_COUNT+1)

	// we count the amount of requests that were made, this should be the amount of categories
	// again we do plus 1 due to zero indexing.. and 3 extra calls for the runedate from the outdated tests
	assert.Equal(t, int32(CATEGORY_COUNT+1+3), count)
}

func TestIsOutdated(t *testing.T) {
	var count int32
	client := lib.NewTestClient(func(req *http.Request) (int, string) {
		atomic.AddInt32(&count, 1)
		if req.URL.Path == "/m=itemdb_rs/api/info.json" {
			return 200, `{"lastConfigUpdateRuneday":7040}`
		}
		t.Error("Unexpected http call", req)
		return 500, "" // just making the compiler happy
	})

	f, err := os.Open("testdata/meta.0.json")
	assert.Nil(t, err)
	defer f.Close()

	meta, err := ReadMetadata(f)
	assert.Nil(t, err)

	outdated, err := meta.IsOutdated(client)
	assert.Nil(t, err)
	assert.True(t, outdated)

	assert.Equal(t, int32(1), count)
}

func TestIsEmpty(t *testing.T) {
	empty := &meta{
		Categories: make(map[int]category),
	}

	assert.True(t, empty.IsEmpty())

	notEmpty := &meta{
		Categories: map[int]category{
			0: {
				Count: map[string]int{
					"a": 12,
				},
			},
		},
	}

	assert.False(t, notEmpty.IsEmpty())
}

package itemdb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

func TestSerialize(t *testing.T) {
	writer := bytes.Buffer{}
	db := New()
	db.Writer = &writer
	assert.Nil(t, db.add(&ge.Item{
		Name: "Test",
		ID:   42,
	}))

	assert.Equal(t, "{\"icon\":\"\",\"icon_large\":\"\",\"id\":42,\"type\":\"\",\"name\":\"Test\",\"description\":\"\"}\n", writer.String())

	out := bytes.Buffer{}
	assert.Nil(t, db.Serialize(&out))

	assert.Equal(t, "{\"icon\":\"\",\"icon_large\":\"\",\"id\":42,\"type\":\"\",\"name\":\"Test\",\"description\":\"\"}\n", out.String())
}

func TestNewFromReader(t *testing.T) {
	in := bytes.NewBufferString("{\"icon\":\"\",\"icon_large\":\"\",\"id\":42,\"type\":\"\",\"name\":\"Test\",\"description\":\"\"}\n")

	db, err := NewFromReader(in)
	assert.Nil(t, err)

	assert.Len(t, db.idToItems, 1)
	assert.Equal(t, "Test", db.idToItems[42].Name)
}

func TestSearch(t *testing.T) {
	db := New()
	assert.Nil(t, db.add(&ge.Item{
		Name: "Test",
		ID:   42,
	}))

	res, err := db.SearchItems("Test")
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, int64(42), res[0].ID)

	res, err = db.SearchItems("t")
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, int64(42), res[0].ID)
}

func TestGetItem(t *testing.T) {
	db := New()
	assert.Nil(t, db.add(&ge.Item{
		Name: "Test",
		ID:   42,
	}))

	res, err := db.GetItem(42)
	assert.Nil(t, err)
	assert.Equal(t, int64(42), res.ID)

	res, err = db.GetItem(1337)
	assert.Error(t, err)
	assert.Nil(t, res)
}

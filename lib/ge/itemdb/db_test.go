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

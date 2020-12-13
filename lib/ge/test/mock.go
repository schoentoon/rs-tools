package ge

// this is purely in its own folder so it'll be ignored by coverage, as this is simply
// a helper for testing anyway

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

type TestingGe struct {
	T          *testing.T
	NextError  chan error
	NextGraph  chan *ge.Graph
	NextItem   chan *ge.Item
	SkipErrors int
}

func MockGe(t *testing.T) *TestingGe {
	return &TestingGe{
		T:          t,
		NextError:  make(chan error, 8),
		NextGraph:  make(chan *ge.Graph, 8),
		NextItem:   make(chan *ge.Item, 8),
		SkipErrors: 0,
	}
}

func (t *TestingGe) Close() {
	if len(t.NextError) > 0 {
		assert.Fail(t.T, "Still errors planned that didn't happen")
	}
	if len(t.NextGraph) > 0 {
		assert.Fail(t.T, "Still PriceGraph calls planned that didn't happen")
	}
	if len(t.NextItem) > 0 {
		assert.Fail(t.T, "Still GetItem calls planned that didn't happen")
	}
}

func (t *TestingGe) PriceGraph(itemID int64) (*ge.Graph, error) {
	if len(t.NextError) > 0 {
		if t.SkipErrors > 0 {
			t.SkipErrors--
		} else {
			return nil, <-t.NextError
		}
	}
	if len(t.NextGraph) == 0 {
		assert.FailNow(t.T, fmt.Sprintf("Unexpected call: PriceGraph(%d)", itemID))
	}
	return <-t.NextGraph, nil
}

func (t *TestingGe) GetItem(itemID int64) (*ge.Item, error) {
	if len(t.NextError) > 0 {
		if t.SkipErrors > 0 {
			t.SkipErrors--
		} else {
			return nil, <-t.NextError
		}
	}
	if len(t.NextItem) == 0 {
		assert.FailNow(t.T, fmt.Sprintf("Unexpected call: NextItem(%d)", itemID))
	}
	return <-t.NextItem, nil
}

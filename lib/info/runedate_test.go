package info

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunedate(t *testing.T) {
	now := time.Date(2021, time.June, 4, 0, 0, 0, 0, time.UTC)

	runedate := GetRunedate(now)

	assert.Equal(t, 7037, runedate)
}

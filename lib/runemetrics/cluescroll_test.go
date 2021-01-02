package runemetrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClueScroll(t *testing.T) {
	cases := []struct {
		Activity    Activity
		ClueScroll  ClueScroll
		ExpectError bool
	}{
		{
			Activity:    Activity{Details: "I have completed a master treasure trail.   ", Text: "Master treasure trail completed."},
			ClueScroll:  ClueScroll{Difficulty: Master},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I have completed a master treasure trail. I got an Ice dye out of it.", Text: "Master treasure trail completed."},
			ClueScroll:  ClueScroll{Difficulty: Master, Loot: "an Ice dye"},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I killed 9 Magisters, the unkillable holder of the Crossing.", Text: "I killed 9 Magisters."},
			ExpectError: true,
		},
		{
			Activity:    Activity{Details: "I have completed an elite treasure trail.   ", Text: "Elite treasure trail completed."},
			ClueScroll:  ClueScroll{Difficulty: Elite},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I have completed an impossible treasure trail.   ", Text: "Elite treasure trail completed."},
			ExpectError: true,
		},
	}

	for _, c := range cases {
		out, err := c.Activity.ClueScroll()
		if err != nil {
			if !c.ExpectError {
				assert.Fail(t, "Unexpected error", err)
			}
		} else if assert.NotNil(t, out, c) {
			assert.Equal(t, c.ClueScroll.Difficulty, out.Difficulty, c)
			assert.Equal(t, c.ClueScroll.Loot, out.Loot, c)
		}
	}
}

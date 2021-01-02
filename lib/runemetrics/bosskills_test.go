package runemetrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBossKillsParsing(t *testing.T) {
	cases := []struct {
		Activity    Activity
		BossKills   BossKills
		ExpectError bool
	}{
		{
			Activity:    Activity{Details: "I killed 4 Helwyrs, limb tearing hunters.", Text: "I killed 4 Helwyrs."},
			BossKills:   BossKills{Boss: "Helwyr", Amount: 4},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I killed 3 Gregorovics, all blade wielding terrors.", Text: "I killed 3 Gregorovics."},
			BossKills:   BossKills{Boss: "Gregorovic", Amount: 3},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I defeated the Twin Furies 3 times, all servants of Zamorak.", Text: "I defeated the Twin Furies 3 times."},
			BossKills:   BossKills{Boss: "Twin Furies", Amount: 3},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I killed  a spinner of death, Araxxi.", Text: "I killed  Araxxi."},
			BossKills:   BossKills{Boss: "Araxxi", Amount: 1},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "I killed 31 servants of the god Zamorak, all called K'ril Tsutsaroth. (Hard mode)", Text: "I killed 31 (Hard mode) K'ril Tsutsaroths."},
			BossKills:   BossKills{Boss: "K'ril Tsutsaroth", Amount: 31, Hardmode: true},
			ExpectError: false,
		},
		{
			Activity:    Activity{Details: "After killing a Gorvek and Vindicta, it dropped a Crest of Zaros.", Text: "I found a Crest of Zaros"},
			ExpectError: true,
		},
		{
			Activity:    Activity{Details: "I defeated many waves of TokHaar, before vanquishing the mighty Har'Aken and conquering the Fight Kiln.", Text: "Completed the Fight Kiln"},
			BossKills:   BossKills{Boss: "Har'Aken", Amount: 1},
			ExpectError: false,
		},
	}

	for _, c := range cases {
		out, err := c.Activity.BossKills()
		if err != nil {
			if !c.ExpectError {
				assert.Fail(t, "Unexpected error", err)
			}
		} else if assert.NotNil(t, out, c) {
			assert.Equal(t, c.BossKills.Boss, out.Boss, c)
			assert.Equal(t, c.BossKills.Amount, out.Amount, c)
			assert.Equal(t, c.BossKills.Hardmode, out.Hardmode, c)
		}
	}
}

package runemetrics

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type BossKills struct {
	Boss     string
	Amount   int
	Hardmode bool
}

var iKilledRegex = regexp.MustCompile(`I killed (\d*) ([^\.]+)\.`)
var iDefeatedRegex = regexp.MustCompile(`I defeated ([^\d]+) (\d) times\.`)
var fightKiln = "Completed the Fight Kiln"

func ParseBossKills(activity Activity) (out *BossKills, err error) {
	// Har'Aken being a sort of minigame means it has it's own activity log text
	// so we have a special case checking for that.
	if activity.Text == fightKiln {
		return &BossKills{
			Boss:     "Har'Aken",
			Amount:   1,
			Hardmode: false,
		}, nil
	}

	// some post processing to always strip off whitespaces and some other details
	defer func() {
		if out != nil {
			out.Boss = strings.Trim(out.Boss, " ")
			out.Boss = strings.TrimPrefix(out.Boss, "the ")
		}
	}()

	results := iKilledRegex.FindStringSubmatch(activity.Text)
	if len(results) > 0 {
		out = &BossKills{
			Boss: strings.TrimRight(results[2], "s"),
		}
		if results[1] == "" {
			out.Amount = 1
		} else {
			out.Amount, err = strconv.Atoi(results[1])
			if err != nil {
				// this should be impossible to reach
				return nil, err
			}
		}

		// All the GWD1 bosses that use the tag (Hard mode) at least are listed as "I killed"
		// Need to investigate whether this is the same for other bosses
		if strings.Contains(out.Boss, "(Hard mode)") {
			out.Hardmode = true
			out.Boss = strings.Replace(out.Boss, "(Hard mode)", "", 1)
		}

		return
	}

	results = iDefeatedRegex.FindStringSubmatch(activity.Text)
	if len(results) > 0 {
		out = &BossKills{
			Boss: results[1],
		}
		out.Amount, err = strconv.Atoi(results[2])
		if err != nil {
			// this should be impossible to reach
			return nil, err
		}
		return
	}

	return nil, errors.New("Doesn't seem like a boss kill?")
}

func (a *Activity) BossKills() (*BossKills, error) {
	return ParseBossKills(*a)
}

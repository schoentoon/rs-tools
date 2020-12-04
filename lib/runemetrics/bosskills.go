package runemetrics

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type BossKills struct {
	Boss   string
	Amount int
}

var iKilledRegex = regexp.MustCompile(`I killed (\d*) ([^\.]+)\.`)
var iDefeatedRegex = regexp.MustCompile(`I defeated ([^\d]+) (\d) times\.`)

func ParseBossKills(activity Activity) (out *BossKills, err error) {
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

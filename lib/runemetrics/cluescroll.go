//go:generate stringer -type=ClueDifficulty
package runemetrics

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ClueDifficulty int

const (
	Invalid ClueDifficulty = iota
	Easy
	Medium
	Hard
	Elite
	Master
)

func parseClueDifficulty(in string) (ClueDifficulty, error) {
	for d := Easy; d <= Master; d++ {
		if strings.EqualFold(d.String(), in) {
			return d, nil
		}
	}

	return Invalid, fmt.Errorf("Invalid clue difficult: %s", in)
}

type ClueScroll struct {
	Difficulty ClueDifficulty
	Loot       string // This is going to be empty the majority of the time
}

var iCompletedClue = regexp.MustCompile(`I have completed an? ([^ ]*) treasure trail\.[ +](I got ([a-zA-Z ]+) out of it\.|)`)

func ParseClueScroll(activity Activity) (*ClueScroll, error) {
	results := iCompletedClue.FindStringSubmatch(activity.Details)
	if len(results) > 0 {
		d, err := parseClueDifficulty(results[1])
		if err != nil {
			return nil, err
		}
		return &ClueScroll{
			Difficulty: d,
			Loot:       results[3],
		}, nil
	}

	return nil, errors.New("Doesn't seem like a clue scroll?")
}

func (a *Activity) ClueScroll() (*ClueScroll, error) {
	return ParseClueScroll(*a)
}

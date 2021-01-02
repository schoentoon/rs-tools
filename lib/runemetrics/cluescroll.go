package runemetrics

import (
	"errors"
	"regexp"
)

type ClueScroll struct {
	Difficulty string
	Loot       string // This is going to be empty the majority of the time
}

var iCompletedClue = regexp.MustCompile(`I have completed a ([^ ]*) treasure trail\.[ +](I got ([a-zA-Z ]+) out of it\.|)`)

func ParseClueScroll(activity Activity) (*ClueScroll, error) {
	results := iCompletedClue.FindStringSubmatch(activity.Details)
	if len(results) > 0 {
		return &ClueScroll{
			Difficulty: results[1],
			Loot:       results[3],
		}, nil
	}

	return nil, errors.New("Doesn't seem like a clue scroll?")
}

func (a *Activity) ClueScroll() (*ClueScroll, error) {
	return ParseClueScroll(*a)
}

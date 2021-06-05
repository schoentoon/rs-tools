package info

import "time"

// https://runescape.wiki/w/Runedate
func GetRunedate(t time.Time) int {
	start := time.Date(2002, time.February, 27, 0, 0, 0, 0, time.UTC)

	now := t.UTC()

	return int(now.Sub(start) / (time.Hour * 24))
}

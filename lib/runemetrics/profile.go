package runemetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"
)

func FetchProfile(client *http.Client, username string) (*Profile, error) {
	params := url.Values{
		"user":       {username},
		"activities": {"20"},
	}
	resp, err := client.Get(fmt.Sprintf("https://apps.runescape.com/runemetrics/profile/profile?%s", params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	return ParseProfile(resp.Body)
}

func ParseProfile(r io.Reader) (*Profile, error) {
	out := &Profile{}

	err := json.NewDecoder(r).Decode(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (p *Profile) GetSkill(skill SkillID) (*SkillValue, error) {
	for _, s := range p.Skills {
		if s.ID == skill {
			return &s, nil
		}
	}
	return nil, errors.New("Skill not found? Are you not in the hiscores or something?")
}

func NewAchievementsSince(old, new []Activity) []Activity {
	// pretty pointless if either of the arrays are empty
	if len(old) == 0 || len(new) == 0 {
		return []Activity{}
	}

	// they should be sorted, but let's make sure of that
	sort.Slice(old, func(i, j int) bool { return time.Time(old[i].Date).Unix() > time.Time(old[j].Date).Unix() })
	sort.Slice(new, func(i, j int) bool { return time.Time(new[i].Date).Unix() > time.Time(new[j].Date).Unix() })

	// let's make sure that old is actually the older one
	if time.Time(old[0].Date).Unix() > time.Time(new[0].Date).Unix() {
		old, new = new, old
	}

	out := []Activity{}

	for _, a := range new {
		if old[0] != a {
			out = append(out, a)
		} else {
			break
		}
	}

	return out
}

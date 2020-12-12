package runemetrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	if out.Err != "" {
		return nil, errors.New(out.Err)
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

// NewAchievementsSince returns all the new achievements in new compared to old
// this does assume that old and new are both sorted by date already, which is default
// when it comes fresh from the API. The returned slice will only have entries newer than
// the newest from the other slice. The order of it will not be touched.
func NewAchievementsSince(old, new []Activity) []Activity {
	// pretty pointless if either of the arrays are empty
	if len(old) == 0 || len(new) == 0 {
		return []Activity{}
	}

	newest := func(activities []Activity) Activity {
		highest := activities[0]
		for _, a := range activities {
			if a.Date.Unix() >= highest.Date.Unix() {
				highest = a
			}
		}
		return highest
	}

	latestOld := newest(old)
	latestNew := newest(new)
	latest := latestOld

	// let's make sure that old is actually the older one
	if latestOld.Date.Unix() > latestNew.Date.Unix() {
		new = old
		latest = latestNew
	}

	out := []Activity{}

	for _, a := range new {
		if latest.Date.Unix() == a.Date.Unix() && latest.Details == a.Details {
			break
		}
		out = append(out, a)
	}

	return out
}

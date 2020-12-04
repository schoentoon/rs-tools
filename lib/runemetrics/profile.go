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

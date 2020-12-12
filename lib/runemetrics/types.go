package runemetrics

import "time"

type Profile struct {
	MagicXP          int64  `json:"magic"`
	RangedXP         int64  `json:"ranged"`
	MeleeXP          int64  `json:"melee"`
	Combat           int    `json:"combat"`
	TotalLevel       int    `json:"totalskill"`
	TotalXP          int64  `json:"totalxp"`
	Name             string `json:"name"`
	QuestsComplete   int    `json:"questscomplete"`
	QuestsNotStarted int    `json:"questsnotstarted"`
	QuestsStarted    int    `json:"questsstarted"`
	Rank             string `json:"rank"`

	Activities []Activity   `json:"activities"`
	Skills     []SkillValue `json:"skillvalues"`

	Err string `json:"error"`
}

type ActivityTimeFormat struct {
	*time.Time
}

func (at *ActivityTimeFormat) UnmarshalJSON(in []byte) error {
	t, err := time.Parse(`"02-Jan-2006 15:04"`, string(in))
	if err != nil {
		return err
	}
	*at = ActivityTimeFormat{&t}
	return nil
}

func (at ActivityTimeFormat) MarshalJSON() ([]byte, error) {
	if at.Time == nil {
		return []byte(time.Now().Format(`"02-Jan-2006 15:04"`)), nil
	}
	return []byte(at.Format(`"02-Jan-2006 15:04"`)), nil
}

type Activity struct {
	Date    ActivityTimeFormat `json:"date"`
	Details string             `json:"details"`
	Text    string             `json:"text"`
}

type SkillID int8

const (
	Attack  SkillID = 0
	Defence         = iota
	Strength
	Constitution
	Ranged
	Prayer
	Magic
	Cooking
	Woodcutting
	Fletching
	Fishing
	Firemaking
	Crafting
	Smithing
	Mining
	Herblore
	Agility
	Thieving
	Slayer
	Farming
	Runecrafting
	Hunter
	Construction
	Summoning
	Dungeoneering
	Divination
	Invention
	Archaeology
)

type SkillValue struct {
	Level int     `json:"level"`
	XP    int     `json:"xp"`
	Rank  int     `json:"rank"`
	ID    SkillID `json:"id"`
}

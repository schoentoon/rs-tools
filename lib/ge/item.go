package ge

import (
	"encoding/json"
	"fmt"
)

/*
{
  "item": {
    "icon": "https://secure.runescape.com/m=itemdb_rs/1603720907702_obj_sprite.gif?id=245",
    "icon_large": "https://secure.runescape.com/m=itemdb_rs/1603720907702_obj_big.gif?id=245",
    "id": 245,
    "type": "Herblore materials",
    "name": "Wine of Zamorak",
    "description": "A jug full of Wine of Zamorak."
  }
}
*/

type Item struct {
	Icon        string `json:"icon"`
	IconLarge   string `json:"icon_large"`
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *Ge) GetItem(itemID int64) (*Item, error) {
	resp, err := g.Client.Get(fmt.Sprintf("https://secure.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item=%d", itemID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Status: %d %s", resp.StatusCode, resp.Status)
	}

	out := struct {
		Item Item `json:"item"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out.Item, nil
}

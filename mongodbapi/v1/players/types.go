package mongodbapi

import (
	"image/color"
	"time"
)

// DBPlayerPofile - Player data db entry struct
type DBPlayerPofile struct {
	Realm      string    `json:"realm" bson:"realm,omitempty"`
	ID         int       `json:"player_id" bson:"_id,omitempty"`
	ClanID     int       `json:"clan_id" bson:"clan_id,omitempty"`
	ClanTag    string    `json:"clan_tag" bson:"clan_tag,omitempty"`
	Nickname   string    `json:"nickname" bson:"nickname,omitempty"`
	ClanName   string    `json:"clan_name" bson:"clan_name,omitempty"`
	CareerWN8  int       `json:"career_wn8" bson:"career_wn8,omitempty"`
	PlayerPins []UserPin `json:"player_pins" bson:"player_pins,omitempty"`
	LastBattle time.Time `json:"last_battle_time" bson:"last_battle_time,omitempty"`
}

// UserPin -
type UserPin struct {
	URL    string `bson:"url"`
	Label  string `bson:"label"`
	Weight int    `bson:"weight"`

	Size      int        `bson:"-"`
	Glow      bool       `bson:"glow"`
	GlowColor color.RGBA `bson:"glow_color"`
}

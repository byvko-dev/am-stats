package mongodbapi

import (
	"image/color"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DBPlayerProfile - Player data db entry struct
type DBPlayerProfile struct {
	Realm      string               `json:"realm" bson:"realm,omitempty"`
	ID         int                  `json:"player_id" bson:"_id,omitempty"`
	ClanID     int                  `json:"clan_id" bson:"clan_id,omitempty"`
	ClanTag    string               `json:"clan_tag" bson:"clan_tag,omitempty"`
	Nickname   string               `json:"nickname" bson:"nickname,omitempty"`
	ClanName   string               `json:"clan_name" bson:"clan_name,omitempty"`
	CareerWN8  int                  `json:"career_wn8" bson:"career_wn8,omitempty"`
	UniquePins []UserPin            `json:"unique_pins" bson:"unique_pins,omitempty"`
	PlayerPins []primitive.ObjectID `json:"player_pins" bson:"player_pins,omitempty"`
	LastBattle time.Time            `json:"last_battle_time" bson:"last_battle_time,omitempty"`
}

// UserPin -
type UserPin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	URL         string             `bson:"url"`
	Name        string             `bson:"name"`
	Label       string             `bson:"label"`
	Weight      int                `bson:"weight"`
	Hidden      bool               `bson:"hidden"`
	Description string             `bson:"description"`

	Size      int        `bson:"-"`
	Glow      bool       `bson:"glow"`
	GlowColor color.RGBA `bson:"glow_color"`
}

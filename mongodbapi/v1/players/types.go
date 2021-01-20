package mongodbapi

import "time"

// DBPlayerPofile - Player data db entry struct
type DBPlayerPofile struct {
	ID         int       `json:"player_id" bson:"_id,omitempty"`
	ClanID     int       `json:"clan_id" bson:"clan_id,omitempty"`
	ClanName   string    `json:"clan_name" bson:"clan_name,omitempty"`
	ClanTag    string    `json:"clan_tag" bson:"clan_tag,omitempty"`
	LastBattle time.Time `json:"last_battle_time" bson:"last_battle_time,omitempty"`
	Nickname   string    `json:"nickname" bson:"nickname,omitempty"`
	Realm      string    `json:"realm" bson:"realm,omitempty"`
	CareerWN8  int       `json:"career_wn8" bson:"career_wn8,omitempty"`
}

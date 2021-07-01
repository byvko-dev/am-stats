package mongodbapi

import (
	"time"

	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// MedalWeight - Object for calculating per medal scores
type MedalWeight struct {
	Name    string `json:"medal" bson:"medal"`
	Weight  int    `json:"weight" bson:"weight"`
	IconURL string `json:"-" bson:"-"`
	Score   int    `json:"-" bson:"-"`
}

// AchievementsPlayerData -
type AchievementsPlayerData struct {
	Timestamp time.Time               `json:"timestamp,omitempty" bson:"timestamp"`
	Nickname  string                  `json:"nickname,omitempty" bson:"nickname"`
	ClanTag   string                  `json:"clan_tag,omitempty" bson:"clan_tag"`
	ClanID    int                     `json:"clan_id,omitempty" bson:"clan_id"`
	Realm     string                  `json:"realm,omitempty" bson:"realm"`
	PID       int                     `json:"_id,omitempty" bson:"_id"`
	Score     int                     `json:"score" bson:"score"`
	Data      wgapi.AchievementsFrame `json:"data" bson:"data"`
	Medals    []MedalWeight           `json:"-" bson:"-"`
}

// AchievementsMap -
type AchievementsMap struct {
	Timestamp time.Time               `json:"timestamp,omitempty" bson:"timestamp"`
	PID       int                     `json:"_id,omitempty" bson:"_id"`
	Nickname  string                  `json:"nickname,omitempty" bson:"-"`
	ClanTag   string                  `json:"clan_tag,omitempty" bson:"-"`
	ClanID    int                     `json:"clan_id,omitempty" bson:"-"`
	Realm     string                  `json:"realm,omitempty" bson:"realm"`
	Data      wgapi.AchievementsFrame `json:"data" bson:"data"`
	Score     int                     `json:"score" bson:"score"`
	Medals    []MedalWeight           `json:"-" bson:"-"`
}

// ClanAchievements -
type ClanAchievements struct {
	ClanID    int                     `json:"_id,omitempty" bson:"_id"`
	ClanTag   string                  `json:"clan_tag,omitempty" bson:"clan_tag"`
	Realm     string                  `json:"realm,omitempty" bson:"realm"`
	Members   int                     `json:"members,omitempty" bson:"members"`
	Timestamp time.Time               `json:"timestamp,omitempty" bson:"timestamp"`
	Data      wgapi.AchievementsFrame `json:"data,omitempty" bson:"data,omitempty"`
	Score     int                     `json:"score,omitempty" bson:"score,omitempty"`
	Position  int                     `json:"position,omitempty" bson:"position,omitempty"`
	Medals    []MedalWeight           `json:"-" bson:"-"`
}

// CachedMedalsRequest -
type CachedMedalsRequest struct {
	Request struct {
		Realm  string        `bson:"realm"`
		Days   int           `bson:"days"`
		Medals []MedalWeight `bson:"medals"`
	} `bson:"request"`
	Result struct {
		TotalScore    int                      `bson:"total_score"`
		SortedPlayers []AchievementsPlayerData `bson:"sorted_players"`
		UpdatedAt     time.Time                `bson:"updated_timestamp"`
	}
	LastRequested time.Time `bson:"requested_timestamp"`
}

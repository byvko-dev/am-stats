package mongodbapi

import (
	"time"

	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// AchievementsPlayerData -
type AchievementsPlayerData struct {
	Timestamp time.Time               `json:"timestamp,omitempty" bson:"timestamp"`
	PID       int                     `json:"_id,omitempty" bson:"_id"`
	Nickname  string                  `json:"nickname,omitempty" bson:"-"`
	ClanTag   string                  `json:"clan_tag,omitempty" bson:"-"`
	ClanID    int                     `json:"clan_id,omitempty" bson:"-"`
	Realm     string                  `json:"realm,omitempty" bson:"realm"`
	Data      wgapi.AchievementsFrame `json:"data" bson:"data"`
	Score     int                     `json:"score" bson:"score"`
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
}

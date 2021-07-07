package dataprep

import (
	"time"

	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
)

// AchievementsPIDPos -
type AchievementsPIDPos struct {
	Position int `json:"position"`
	dbAch.AchievementsPlayerData
}

// AchievementsPlayerSettings -
type AchievementsPlayerSettings struct {
	PersonalCutoff time.Time
}

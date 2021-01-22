package dataprep

import dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"

type AchievementsPIDPos struct {
	Position int `json:"position"`
	dbAch.AchievementsPlayerData
}

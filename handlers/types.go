package handlers

import dataprep "github.com/cufee/am-stats/dataprep/achievements"

// StatsRequest - Request for stats image
type StatsRequest struct {
	PlayerID  int    `json:"player_id"`
	Premium   bool   `json:"premium"`
	Verified  bool   `json:"verified"`
	Realm     string `json:"realm"`
	Days      int    `json:"days"`
	Sort      string `json:"sort_key"`
	TankLimit int    `json:"detailed_limit"`
	BgURL     string `json:"bg_url"`
}

// AchievementsRequest - Request for achievements data
type AchievementsRequest struct {
	BgURL    string                 `json:"bg_url"`
	Premium  bool                   `json:"premium"`
	Verified bool                   `json:"verified"`
	ClanTag  string                 `json:"clan_tag"`
	PlayerID int                    `json:"player_id"`
	Days     int                    `json:"days"`
	Limit    int                    `json:"limit"`
	Realm    string                 `json:"realm"`
	Medals   []dataprep.MedalWeight `json:"medals"`
}

package dataprep

import (
	dbPlayers "github.com/cufee/am-stats/mongodbapi/v1/players"
	dbStats "github.com/cufee/am-stats/mongodbapi/v1/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ExportData - Struct to export final data for use in Python bot
type ExportData struct {
	PlayerCache   dbPlayers.DBPlayerPofile `json:"player_cache"`
	PlayerDetails wgapi.PlayerProfile      `json:"player_details"`
	SessionStats  dbStats.Session          `json:"session"`
	LastSession   dbStats.RetroSession     `json:"last_session"`
	Analytics     `json:"analytics"`
}

// Analytics data for a request
type Analytics struct {
	TimeToComplete float64 `json:"request_time_sec"`
}

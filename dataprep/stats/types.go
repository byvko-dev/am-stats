package dataprep

import (
	db "github.com/cufee/am-stats/mongodbapi"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ExportData - Struct to export final data for use in Python bot
type ExportData struct {
	PlayerDetails wgapi.PlayerProfile `json:"player_details"`
	SessionStats  db.Session          `json:"session"`
	LastSession   db.RetroSession     `json:"last_session"`
	Analytics     `json:"analytics"`
}

// Analytics data for a request
type Analytics struct {
	TimeToComplete float64 `json:"request_time_sec"`
}

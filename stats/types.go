package stats

import (
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ExportData - Struct to export final data for use in Python bot
type ExportData struct {
	Vehicles		[]wgapi.VehicleStats	`json:"session_detailed"`
	PlayerDetails	wgapi.PlayerProfile		`json:"player_details"`
	SessionStats	wgapi.StatsFrame		`json:"session_all"`
}
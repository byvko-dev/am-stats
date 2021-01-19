package dataprep

import (
	stats "github.com/cufee/am-stats/dataprep/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ExportAchievementsSession - Export achievements from a session
func ExportAchievementsSession(pid int, realm string, days int) (wgapi.AchievementsFrame, error) {
	// Get session
	data, err := stats.ExportSessionAsStruct(pid, realm, days)

	// Return Achievemenets
	return data.SessionStats.Achievements, err
}

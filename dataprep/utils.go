package dataprep

import (
	"time"

	db "github.com/cufee/am-stats/mongodbapi/v1/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// LiveToSession - Convert Live data to Session
func LiveToSession(profile wgapi.PlayerProfile, vehicles []wgapi.VehicleStats, achievements wgapi.AchievementsFrame) (liveSession db.Session) {
	liveSession.Vehicles = vehicles
	liveSession.Achievements = achievements
	liveSession.PlayerID = profile.ID
	liveSession.LastBattle = time.Unix(int64(profile.LastBattle), 0)
	liveSession.BattlesAll = profile.Stats.All.Battles
	liveSession.StatsAll = profile.Stats.All
	liveSession.BattlesRating = profile.Stats.Rating.Battles
	liveSession.StatsRating = profile.Stats.Rating
	return liveSession
}

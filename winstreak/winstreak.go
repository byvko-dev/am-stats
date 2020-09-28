package winstreak

import (
	"fmt"
	"log"

	db "github.com/cufee/am-stats/mongodbapi"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// CheckStreak - Check player streak and update db
func CheckStreak(pid int, stats wgapi.StatsFrame) (streakData db.PlayerStreak, err error) {
	streakData, err = db.GetStreak(pid)
	if err != nil {
		switch err.Error() {
		case "mongo: no documents in result":
			// New user
			streakData.PlayerID = &pid
			streakData.Battles = &stats.Battles
			streakData.Losses = &stats.Losses
			streakData.Streak = 0
			// Update DB
			err := db.UpdateStreak(streakData)
			return streakData, err
		default:
			log.Print(err)
			return streakData, err
		}
	}
	if stats.Battles >= *streakData.Battles && stats.Losses == *streakData.Losses {
		// Streak increased or did not change
		newStreak := streakData.Streak + stats.Battles - *streakData.Battles
		// Update DB
		if newStreak > streakData.BestStreak {
			streakData.BestStreak = newStreak
		}
		streakData.Streak = newStreak
		streakData.Battles = &stats.Battles
		streakData.Losses = &stats.Losses
		err := db.UpdateStreak(streakData)
		return streakData, err
	}
	if stats.Battles >= *streakData.Battles && stats.Losses > *streakData.Losses {
		// Streak broken
		// Update DB
		streakData.Streak = 0
		streakData.Battles = &stats.Battles
		streakData.Losses = &stats.Losses
		err := db.UpdateStreak(streakData)
		return streakData, err
	}
	return streakData, fmt.Errorf("invalid data")
}

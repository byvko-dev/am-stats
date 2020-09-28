package winstreak

import (
	"fmt"
	"log"

	db "github.com/cufee/am-stats/mongodbapi"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// CheckStreak - Check player streak and update db
func CheckStreak(pid int, stats wgapi.StatsFrame) (streak int, err error) {
	streakData, err := db.GetStreak(pid)
	if err != nil {
		switch err.Error() {
		case "mongo: no documents in result":
			// New user
			streak = 0
			streakData.PlayerID = &pid
			streakData.Battles = &stats.Battles
			streakData.Losses = &stats.Losses
			streakData.Streak = &streak
			// Update DB
			err := db.UpdateStreak(streakData)
			return streak, err
		default:
			log.Print(err)
			return streak, err
		}
	}

	log.Print(stats.Battles, stats.Losses)
	log.Print(*streakData.Battles, *streakData.Losses)

	if stats.Battles >= *streakData.Battles && stats.Losses == *streakData.Losses {
		// Streak increased or did not change
		newStreak := *streakData.Streak + stats.Battles - *streakData.Battles
		streak = newStreak
		// Update DB
		streakData.Streak = &newStreak
		streakData.Battles = &stats.Battles
		streakData.Losses = &stats.Losses
		err := db.UpdateStreak(streakData)
		return streak, err
	}
	if stats.Battles >= *streakData.Battles && stats.Losses > *streakData.Losses {
		// Streak broken
		streak = 0
		// Update DB
		streakData.Streak = &streak
		streakData.Battles = &stats.Battles
		streakData.Losses = &stats.Losses
		err := db.UpdateStreak(streakData)
		return streak, err
	}
	return streak, fmt.Errorf("invalid data")
}

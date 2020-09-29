package winstreak

import (
	"fmt"
	"log"
	"math"

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
			streakData.MinStreak = int(math.Ceil(float64(stats.Battles) / (float64(stats.Losses) + 1)))
			streakData.MaxStreak = stats.Battles - stats.Losses
			streakData.BestStreak = streakData.MinStreak
			streakData.Streak = 0
			// Update DB
			err := db.UpdateStreak(streakData)
			return streakData, err
		default:
			log.Print(err)
			return streakData, err
		}
	}
	if stats.Battles == *streakData.Battles {
		// No battles played
		return streakData, err
	}
	if stats.Battles < *streakData.Battles {
		// There is an error in the DB record, fixing
		streakData.Battles = &stats.Battles
		streakData.Losses = &stats.Losses
		streakData.Streak = 0
		err := db.UpdateStreak(streakData)
		return streakData, err
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
	if stats.Battles >= *streakData.Battles && stats.Losses != *streakData.Losses {
		// Calc minimum possible streak
		newStreak := int(math.Ceil(float64(stats.Battles-*streakData.Battles) / (float64(stats.Losses-*streakData.Losses) + 1)))
		if newStreak > streakData.BestStreak {
			streakData.BestStreak = newStreak
		}
		// Update DB
		streakData.Streak = newStreak
		streakData.Battles = &stats.Battles
		streakData.Losses = &stats.Losses
		err := db.UpdateStreak(streakData)
		return streakData, err
	}
	return streakData, fmt.Errorf("invalid data")
}

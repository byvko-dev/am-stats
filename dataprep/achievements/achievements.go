package dataprep

import (
	"log"
	"reflect"
	"strings"
	"sync"

	stats "github.com/cufee/am-stats/dataprep/stats"
	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbPlayers "github.com/cufee/am-stats/mongodbapi/v1/players"
	dbStats "github.com/cufee/am-stats/mongodbapi/v1/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ExportAchievementsSession - Export achievements from a session
func ExportAchievementsSession(pid int, realm string, days int) (wgapi.AchievementsFrame, error) {
	// Get session
	data, err := stats.ExportSessionAsStruct(pid, realm, days)

	// Return Achievemenets
	return data.SessionStats.Achievements, err
}

// ExportAchievementsLeaderboard - Export achievements from a session
func ExportAchievementsLeaderboard(realm string) (export []dbAch.AchievementsPlayerData, err error) {
	// Get leaderboard
	// fields := []string{}

	medals := []MedalWeight{{"MarkOfMastery", 4}, {"MarkOfMasteryI", 3}, {"MarkOfMasteryII", 2}, {"MarkOfMasteryIII", 1}}

	// Generate fields
	fields := []string{}
	for _, m := range medals {
		fields = append(fields, strings.ToLower(m.Name))
	}

	// Get data
	data, err := dbAch.GetPlayerAchievementsLb(realm, fields...)
	if err != nil {
		return []dbAch.AchievementsPlayerData{}, err
	}

	dataChan := make(chan dbAch.AchievementsPlayerData, len(data))
	var wg sync.WaitGroup
	// Fill nicknames and clan tags
	for _, player := range data {
		wg.Add(1)

		go func(player dbAch.AchievementsPlayerData) {
			defer wg.Done()

			// Get player profile
			playerData, err := dbPlayers.GetPlayerProfile(player.PID)
			if err != nil {
				return
			}

			// Get player cached achievements
			achCache, err := dbStats.GetPlayerSessionAchievements(player.PID, 1, fields...)
			if achCache == (dbAch.AchievementsPlayerData{}).Data {
				return
			}

			// Get diff
			newData := player.Data.Diff(achCache)
			if newData == (dbAch.AchievementsPlayerData{}.Data) {
				return
			}

			// Fill name and clan tag
			player.Nickname = playerData.Nickname
			player.ClanTag = playerData.ClanTag
			player.Data = newData

			for _, m := range medals {
				cnt := getField(&player.Data.Achievements, m.Name)
				if cnt > 0 {
					player.Score += cnt * m.Weight
					log.Print(player.Score)
				}
			}

			// Send to chan
			dataChan <- player
		}(player)
	}
	wg.Wait()
	close(dataChan)

	// Quicksort
	// TODO

	// Export
	for d := range dataChan {
		export = append(export, d)
	}
	return export, err
}

type playerDataSorted struct {
	index int
	data  dbAch.AchievementsPlayerData
}

func getField(v *wgapi.Achievements, field string) int {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	if f == (reflect.Value{}) {
		return 0
	}
	return int(f.Int())
}

package dataprep

import (
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
func ExportAchievementsLeaderboard(realm string, limit int, checkPid int, medals ...MedalWeight) (export []dbAch.AchievementsPlayerData, checkPos int, err error) {
	// Generate fields
	fields := []string{}
	for _, m := range medals {
		fields = append(fields, strings.ToLower(m.Name))
	}

	// Get data
	data, err := dbAch.GetPlayerAchievementsLb(realm, fields...)
	if err != nil {
		return []dbAch.AchievementsPlayerData{}, checkPos, err
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
			achCache, err := dbStats.GetPlayerSessionAchievements(player.PID, 0, fields...)
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
				cnt := getField(&player.Data, m.Name)
				if cnt > 0 {
					player.Score += cnt * m.Weight
				}
			}

			// Send to chan
			dataChan <- player
		}(player)
	}
	wg.Wait()
	close(dataChan)

	// Export
	for d := range dataChan {
		export = append(export, d)
	}

	// Quicksort
	sorted := quickSortPlayers(export)

	// Check Pid position
	if checkPid != 0 {
		for i, d := range sorted {
			if d.PID == checkPid {
				checkPos = i + 1
				break
			}
		}
	}

	// Check limit
	if len(sorted) > limit {
		return sorted[:limit], checkPos, err
	}
	return sorted, checkPos, err

}

// ExportClanAchievementsLB - Export clan achievements LB
func ExportClanAchievementsLB(realm string) (export []dbAch.AchievementsPlayerData, err error) {
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
			achCache, err := dbStats.GetPlayerSessionAchievements(player.PID, 0, fields...)
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
				cnt := getField(&player.Data, m.Name)
				if cnt > 0 {
					player.Score += cnt * m.Weight
				}
			}

			// Send to chan
			dataChan <- player
		}(player)
	}
	wg.Wait()
	close(dataChan)

	// Export
	for d := range dataChan {
		export = append(export, d)
	}
	return quickSortPlayers(export), err
}

func getField(v *wgapi.AchievementsFrame, field string) int {
	r := reflect.ValueOf(v.Achievements)
	f := reflect.Indirect(r).FieldByName(field)
	if f == (reflect.Value{}) {
		return 0
	}
	return int(f.Int())
}

// QuickSort is a quick sort algorithm
func quickSortPlayers(arr []dbAch.AchievementsPlayerData) []dbAch.AchievementsPlayerData {
	// clone arr to keep immutability
	newArr := make([]dbAch.AchievementsPlayerData, len(arr))

	for i, v := range arr {
		newArr[i] = v
	}

	// call recursive funciton with initial values
	recursiveSort(newArr, 0, len(newArr)-1)

	// at this point newArr is sorted
	return newArr
}

func recursiveSort(arr []dbAch.AchievementsPlayerData, start, end int) {
	if (end - start) < 1 {
		return
	}

	pivot := arr[end]
	splitIndex := start

	// Iterate sub array to find values less than pivot
	//   and move them to the beginning of the array
	//   keeping splitIndex denoting less-value array size
	for i := start; i < end; i++ {
		if arr[i].Score > pivot.Score {
			if splitIndex != i {
				temp := arr[splitIndex]

				arr[splitIndex] = arr[i]
				arr[i] = temp
			}

			splitIndex++
		}
	}

	arr[end] = arr[splitIndex]
	arr[splitIndex] = pivot

	recursiveSort(arr, start, splitIndex-1)
	recursiveSort(arr, splitIndex+1, end)
}

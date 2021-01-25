package dataprep

import (
	"reflect"
	"strings"
	"sync"

	stats "github.com/cufee/am-stats/dataprep/stats"
	dbAch "github.com/cufee/am-stats/mongodbapi/v1/achievements"
	dbPlayers "github.com/cufee/am-stats/mongodbapi/v1/players"
	dbStats "github.com/cufee/am-stats/mongodbapi/v1/stats"
	"github.com/cufee/am-stats/utils"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ExportAchievementsSession - Export achievements from a session
func ExportAchievementsSession(pid int, realm string, days int) (wgapi.AchievementsFrame, error) {
	// Get session
	data, err := stats.ExportSessionAsStruct(pid, realm, days)

	// Return Achievemenets
	return data.SessionStats.Achievements, err
}

// ExportClanAchievementsByID - Export clan achievements LB by clan ID
func ExportClanAchievementsByID(clanID int, realm string, days int, medals ...dbAch.MedalWeight) (export []dbAch.AchievementsPlayerData, clanTotalScore int, err error) {
	// Get clan members from WG
	ClanProfile, err := wgapi.ClanDataByID(clanID, realm)
	if err != nil {
		return export, clanTotalScore, err
	}

	// Get clan leaderboard
	return exportAchievementsByPIDs(ClanProfile.MembersIds, days, medals...)
}

// ExportClanAchievementsByTag - Export clan achievements by clan tag
func ExportClanAchievementsByTag(clanTag string, realm string, days int, medals ...dbAch.MedalWeight) (export []dbAch.AchievementsPlayerData, clanTotalScore int, err error) {
	// Get clan members from WG
	ClanProfile, err := wgapi.ClanDataByTag(clanTag, realm)
	if err != nil {
		return export, clanTotalScore, err
	}

	// Get clan leaderboard
	return exportAchievementsByPIDs(ClanProfile.MembersIds, days, medals...)
}

// ExportClanAchievementsLbByRealm - Export clan achievements LB by realm
func ExportClanAchievementsLbByRealm(realm string, checkPID int, days int, limit int, medals ...dbAch.MedalWeight) (export []dbAch.ClanAchievements, checkData dbAch.ClanAchievements, err error) {
	// Timer
	timer := utils.Timer{Name: "get players on realm", FunctionName: "ExportClanAchievementsLbByRealm", Enabled: false}
	timer.Start()

	// Get realm players
	pidSlice, err := dbPlayers.GetRealmPlayers(realm)
	if err != nil {
		return export, checkData, err
	}

	// Timer
	timer.Reset("get leaderboard")

	// Get Leaderboard
	leaderboard, _, err := exportAchievementsByPIDs(pidSlice, days, medals...)
	if err != nil {
		return export, checkData, err
	}

	// Timer
	timer.Reset("sort players by clan")

	// Sort by clan
	clanMap := make(map[int]dbAch.ClanAchievements)
	for _, p := range leaderboard {
		clanData := clanMap[p.ClanID]
		if p.ClanID == 0 {
			continue
		}

		for _, m := range medals {
			oldVal := getField(clanData.Data, m.Name)
			pScore := getField(p.Data, m.Name)
			clanData.Data = setField(clanData.Data, m.Name, (oldVal + pScore))
		}

		if clanData.Timestamp.Before(p.Timestamp) {
			clanData.Timestamp = p.Timestamp
		}

		clanData.ClanID = p.ClanID
		clanData.ClanTag = p.ClanTag
		clanData.Score += p.Score
		clanData.Members++
		clanMap[p.ClanID] = clanData

		if checkPID != 0 && p.ClanID != 0 && p.PID == checkPID {
			checkData = clanData
		}
	}
	// Create a slice
	for _, clan := range clanMap {
		export = append(export, clan)
	}

	// Timer
	timer.Reset("sort clans by score")

	// Sort
	export = quickSortClans(export)

	// Timer
	timer.End()

	// Get clan check position
	if checkPID != 0 {
		for i, c := range export {
			if c.ClanID == checkData.ClanID {
				checkData.Position = i + 1
			}
		}
	}

	if len(export) > limit {
		return export[:limit], checkData, err
	}
	return export, checkData, err
}

// ExportAchievementsLeaderboard - Export achievements from a session
func ExportAchievementsLeaderboard(realm string, days int, limit int, checkPid int, medals ...dbAch.MedalWeight) (export []dbAch.AchievementsPlayerData, chackData AchievementsPIDPos, err error) {
	// Get realm players
	pidSlice, err := dbPlayers.GetRealmPlayers(realm)
	if err != nil {
		return export, chackData, err
	}

	// Get Leaderboard
	export, _, err = exportAchievementsByPIDs(pidSlice, days, medals...)
	if err != nil {
		return export, chackData, err
	}

	// Check Pid position
	if checkPid != 0 {
		for i, d := range export {
			if d.PID == checkPid {
				chackData.Position = i + 1
				chackData.AchievementsPlayerData = export[i]
				break
			}
		}
	}

	// Check limit
	if len(export) > limit {
		return export[:limit], chackData, err
	}
	return export, chackData, err
}

// ExportAchievementsByPIDs - Export achievements from a slice of player IDs
func exportAchievementsByPIDs(pidSlice []int, days int, medals ...dbAch.MedalWeight) (export []dbAch.AchievementsPlayerData, totalScore int, err error) {
	// Timer
	timer := utils.Timer{Name: "prep", FunctionName: "exportAchievementsByPIDs", Enabled: false}
	timer.Start()

	// Generate fields
	fields := []string{}
	for _, m := range medals {
		fields = append(fields, strings.ToLower(m.Name))
	}

	dataChan := make(chan dbAch.AchievementsPlayerData, len(pidSlice))
	totalChan := make(chan int, len(pidSlice))
	var wg sync.WaitGroup

	// Timer
	timer.Reset("fill player data")

	// Fill nicknames and clan tags
	for _, pid := range pidSlice {
		wg.Add(1)

		go func(pid int) {
			defer wg.Done()

			player, err := dbAch.GetPlayerAchievements(pid, medals...)
			if err != nil {
				return
			}

			// Get player profile
			playerData, err := dbPlayers.GetPlayerProfile(player.PID)
			if err != nil {
				return
			}

			// Get player cached achievements
			achCache, err := dbStats.GetPlayerSessionAchievements(pid, days, fields...)
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
			player.ClanID = playerData.ClanID
			player.Data = newData

			for _, m := range medals {
				cnt := getField(player.Data, m.Name)
				if cnt > 0 {
					player.Score += cnt * m.Weight
				}
			}

			// Send total score
			totalChan <- player.Score

			// Send to chan
			dataChan <- player
		}(pid)
	}
	wg.Wait()
	close(dataChan)
	close(totalChan)

	// Timer
	timer.Reset("fill clan total scores")

	// Export
	for d := range dataChan {
		export = append(export, d)
	}

	// Clan Score
	for s := range totalChan {
		totalScore += s
	}

	// Timer
	timer.Reset("sort")

	// Quicksort
	sorted := quickSortPlayers(export)

	// Timer
	timer.End()

	return sorted, totalScore, err

}

func getField(data wgapi.AchievementsFrame, field string) int {
	v := reflect.ValueOf(&data.Achievements)
	f := reflect.Indirect(v).FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(field) })
	if f == (reflect.Value{}) {
		return 0
	}
	return int(f.Int())
}

func setField(data wgapi.AchievementsFrame, field string, value int) wgapi.AchievementsFrame {
	v := reflect.ValueOf(&data.Achievements)
	f := reflect.Indirect(v).FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == strings.ToLower(field) })
	if f != (reflect.Value{}) {
		f.SetInt(int64(value))
		return data
	}
	return data
}

// QuickSort is a quick sort algorithm
func quickSortPlayers(arr []dbAch.AchievementsPlayerData) []dbAch.AchievementsPlayerData {
	// clone arr to keep immutability
	newArr := make([]dbAch.AchievementsPlayerData, len(arr))

	for i, v := range arr {
		newArr[i] = v
	}

	// call recursive funciton with initial values
	recursivePlayerSort(newArr, 0, len(newArr)-1)

	// at this point newArr is sorted
	return newArr
}

func recursivePlayerSort(arr []dbAch.AchievementsPlayerData, start, end int) {
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

	recursivePlayerSort(arr, start, splitIndex-1)
	recursivePlayerSort(arr, splitIndex+1, end)
}

// QuickSort is a quick sort algorithm
func quickSortClans(arr []dbAch.ClanAchievements) []dbAch.ClanAchievements {
	// clone arr to keep immutability
	newArr := make([]dbAch.ClanAchievements, len(arr))

	for i, v := range arr {
		newArr[i] = v
	}

	// call recursive funciton with initial values
	recursiveClanSort(newArr, 0, len(newArr)-1)

	// at this point newArr is sorted
	return newArr
}

func recursiveClanSort(arr []dbAch.ClanAchievements, start, end int) {
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

	recursiveClanSort(arr, start, splitIndex-1)
	recursiveClanSort(arr, splitIndex+1, end)
}

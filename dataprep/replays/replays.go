package dataprep

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/dataprep"
	stats "github.com/cufee/am-stats/dataprep/stats"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ProcessReplay Get replay data from url
func ProcessReplay(replayURL string) (summary ReplaySummary, err error) {
	// Get replay data
	summary, err = getReplayDetails(replayURL)
	if err != nil {
		return summary, err
	}

	// Detect realm
	realm := realmFromID(summary.Protagonist)

	// Get player profile data
	wgPlayerData, err := wgapi.PlayerSliceProfileData(realm, append(summary.Allies, summary.Enemies...))
	if err != nil {
		return summary, err
	}

	var badProfiles int
	// Fill player profiles and tanks
	var wg sync.WaitGroup
	players := make(chan ReplayPlayerData, len(summary.Details))
	for _, player := range summary.Details {
		wg.Add(1)
		go func(p ReplayPlayerData, winner int, protagonist int) {
			defer wg.Done()

			// Set profile
			p.Profile = wgPlayerData[fmt.Sprint(p.ID)]

			// Set protagonist
			if p.ID == protagonist {
				p.IsProtagonist = true
			}

			// Set tank profile
			var t wgapi.VehicleStats
			t.TankID = p.VehicleDescr
			detailsToStatsFrame(&p, &t.StatsFrame)
			if p.Team == winner {
				t.StatsFrame.Wins = 1
			} else {
				t.StatsFrame.Losses = 1
			}
			p.TankProfile, _ = stats.CalcVehicleWN8(t)

			// Get mark of mastery
			// TODO

			// Check if a profile is blank
			if p.Profile == (wgapi.PlayerProfile{}) {
				badProfiles++
			}

			// Set team
			if dataprep.IntInSlice(p.ID, summary.Allies) {
				p.Team = 1
			} else {
				p.Team = 2
			}

			players <- p
		}(player, summary.WinnerTeam, summary.Protagonist)
	}
	wg.Wait()
	close(players)

	// Check for bad profiles threshold
	if badProfiles > 3 {
		return summary, fmt.Errorf("replays api error: too many bad profiles")
	}

	// Set summary details
	summary.Details = []ReplayPlayerData{}
	for p := range players {
		summary.Details = append(summary.Details, p)
	}

	// Quicksort by WN8
	summary.Details = quickSortPlayers(summary.Details)
	return summary, err
}

func getReplayDetails(replayURL string) (summary ReplaySummary, err error) {
	// Make URL
	requestURL, err := url.Parse(config.WotInspectorAPI + replayURL)
	if err != nil {
		return summary, fmt.Errorf("replays api error: %s", err.Error())
	}

	// Send request
	var replayRes ReplayDetailsRes
	err = dataprep.DecodeHTTPResponse("GET", make(map[string]string), requestURL, nil, &replayRes)
	if err != nil {
		return summary, fmt.Errorf("replays api error: %s", err.Error())
	}
	if replayRes.Error.Message != "" {
		return summary, fmt.Errorf("replays api error: %s", replayRes.Error.Message)
	}

	// Set Download and File URLs
	replayRes.Data.Summary.DownloadURL = replayRes.Data.DownloadURL
	replayRes.Data.Summary.FileURL = replayURL

	return replayRes.Data.Summary, nil
}

// detailsToStatsFrame - Convert replay details to stats frame fow WN8 calculations
func detailsToStatsFrame(player *ReplayPlayerData, frame *wgapi.StatsFrame) {
	frame.Battles = 1
	frame.DroppedCapturePoints = player.BaseDefendPoints
	frame.CapturePoints = player.BaseCapturePoints
	frame.DamageReceived = player.DamageReceived
	frame.DamageDealt = player.DamageMade
	frame.Spotted = player.EnemiesSpotted
	frame.Frags = player.EnemiesDestroyed
	frame.Shots = player.ShotsMade
	frame.Hits = player.ShotsHit
	frame.Xp = player.Exp
	if player.HitpointsLeft > 0 {
		frame.SurvivedBattles = 1
	}
}

func realmFromID(pidInt int) string {
	switch {
	case pidInt < 500000000:
		return "RU"
	case pidInt < 1000000000:
		return "EU"
	case pidInt < 2000000000:
		return "NA"
	default:
		return "ASIA"
	}
}

// QuickSort is a quick sort algorithm
func quickSortPlayers(arr []ReplayPlayerData) []ReplayPlayerData {
	// clone arr to keep immutability
	newArr := make([]ReplayPlayerData, len(arr))

	for i, v := range arr {
		newArr[i] = v
	}

	// call recursive funciton with initial values
	recursivePlayerSort(newArr, 0, len(newArr)-1)

	// at this point newArr is sorted
	return newArr
}

func recursivePlayerSort(arr []ReplayPlayerData, start, end int) {
	if (end - start) < 1 {
		return
	}

	pivot := arr[end]
	splitIndex := start

	// Iterate sub array to find values less than pivot
	//   and move them to the beginning of the array
	//   keeping splitIndex denoting less-value array size
	for i := start; i < end; i++ {
		if arr[i].TankProfile.TankWN8 > pivot.TankProfile.TankWN8 {
			if splitIndex != i {
				temp := arr[splitIndex]

				arr[splitIndex] = arr[i]
				arr[i] = temp
			}
			splitIndex++
		} else if arr[i].TankProfile.TankWN8 == 0 && pivot.TankProfile.TankWN8 == 0 && arr[i].TankProfile.DamageDealt > pivot.TankProfile.DamageDealt {
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

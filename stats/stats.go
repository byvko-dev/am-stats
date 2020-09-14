package stats

import (
	"sync"
	"math"
	"strconv"
	"time"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	db "github.com/cufee/am-stats/mongodbapi"
	"encoding/json"
)

// CalcVehicleWN8 - Calculate WN8 for a VehicleStats struct
func calcVehicleWN8(tank wgapi.VehicleStats) (battles int, rawRating int, err error) {
	// Get tank averages
	tankAvgData, err := db.GetTankAverages(tank.TankID)
	if err != nil {
		return battles, rawRating, err
	}

	battles = tank.Battles
	// Expected values for WN8
	expDef := tankAvgData.All.DroppedCapturePoints / tankAvgData.All.Battles
	expFrag := tankAvgData.Special.KillsPerBattle
	expSpot := tankAvgData.Special.SpotsPerBattle
	expDmg := tankAvgData.Special.DamagePerBattle
	expWr := tankAvgData.Special.Winrate

	// Actual performance
	pDef := float64(tank.DroppedCapturePoints) / float64(battles)
	pFrag := float64(tank.Frags) / float64(battles)
	pSpot := float64(tank.Spotted) / float64(battles)
	pDmg := float64(tank.DamageDealt) / float64(battles)
	pWr := float64(tank.Wins) / float64(battles) * 100

	// Calculate WN8 metrics
	rDef := pDef / expDef
	rFrag := pFrag / expFrag
	rSpot := pSpot / expSpot
	rDmg := pDmg / expDmg
	rWr := pWr / expWr

	adjustedWr := math.Max(0, ((rWr - 0.71) / (1 - 0.71)))
	adjustedDmg := math.Max(0, ((rDmg - 0.22) / (1 - 0.22)))
	adjustedDef := math.Max(0, (math.Min(adjustedDmg+0.1, (rDef-0.10)/(1-0.10))))
	adjustedSpot := math.Max(0, (math.Min(adjustedDmg+0.1, (rSpot-0.38)/(1-0.38))))
	adjustedFrag := math.Max(0, (math.Min(adjustedDmg+0.2, (rFrag-0.12)/(1-0.12))))

	rating := int(math.Round(((980 * adjustedDmg) + (210 * adjustedDmg * adjustedFrag) + (155 * adjustedFrag * adjustedSpot) + (75 * adjustedDef * adjustedFrag) + (145 * math.Min(1.8, adjustedWr)))))
	rawRating = rating * battles

	return battles, rawRating, nil
}

// Convert Live data to Session
func liveToSession(profile wgapi.PlayerProfile, vehicles []wgapi.VehicleStats) (liveSession db.Session) {
	liveSession.Vehicles 		= vehicles
	liveSession.PlayerID		= profile.ID
	liveSession.LastBattle		= time.Unix(int64(profile.LastBattle), 0)
	liveSession.BattlesAll		= profile.Stats.All.Battles
	liveSession.StatsAll		= profile.Stats.All
	liveSession.BattlesRating	= profile.Stats.Rating.Battles
	liveSession.StatsRating		= profile.Stats.Rating
	return liveSession
}

// SliceDiff - Calculate the difference in two VehicleStats slices
func sessionDiff(oldStats db.Session, liveStats db.Session) (session db.Session) {
	// Convert to RetroSession
	var sessionConv db.Convert = oldStats
	retroSession := sessionConv.ToRetro()
	
	vahiclesChan := make(chan wgapi.VehicleStats, len(liveStats.Vehicles))
	var wg sync.WaitGroup
	for _, newData := range liveStats.Vehicles {
		if newData.Battles == retroSession.Vehicles[strconv.Itoa(newData.TankID)].Battles {
			// Skop if no battles were played
			continue
		}
		// Start go routines
		wg.Add(1)
		go func(newData wgapi.VehicleStats){
			defer wg.Done()
			vahiclesChan <- wgapi.Diff(retroSession.Vehicles[strconv.Itoa(newData.TankID)], newData)
		}(newData)
	}
	wg.Wait()
	close(vahiclesChan)

	for v := range vahiclesChan {
		session.Vehicles = append(session.Vehicles, v)
	}
	session.PlayerID		= liveStats.PlayerID
	session.LastBattle		= liveStats.LastBattle
	session.StatsAll		= wgapi.FrameDiff(oldStats.StatsAll, liveStats.StatsAll)
	session.StatsRating		= wgapi.FrameDiff(oldStats.StatsRating, liveStats.StatsRating)
	session.BattlesAll		= session.StatsAll.Battles
	session.BattlesRating	= session.StatsRating.Battles

	return session
}

// CalcSession - Calculate a new session
func  calcSession(pid int, realm string, days int) (session db.Session, playerProfile wgapi.PlayerProfile, err error) {
	// Get live session
	playerProfile, err = wgapi.PlayerProfileData(pid, realm)
	if err != nil {
		return session, playerProfile,  err
	}
	// Get previous session
	oldSession, err := db.GetPlayerSession(pid, days, playerProfile.Stats.All.Battles)
	if err != nil {
		// No previous session
		return session, playerProfile,  err
	}
	playerVehicles, err := wgapi.PlayerVehicleStats(pid, realm)
	if err != nil {
		return session, playerProfile,  err
	}
	// Calculate session differance and return
	return sessionDiff(oldSession, liveToSession(playerProfile, playerVehicles)), playerProfile,  nil
}

// ExportSessionAsJSON - Export a full player session as a JSON byte slice
func ExportSessionAsJSON(pid int, realm string, days int) (JSONdata []byte, err error) {
	session, playerProfile,  err := calcSession(pid , realm, days)
	if err != nil {
		return JSONdata, err
	}
	var export ExportData
	export.Vehicles 		= session.Vehicles
	export.SessionStats 	= session.StatsAll
	export.PlayerDetails	= playerProfile

	JSONdata, err = json.MarshalIndent(export, "", "  ")
	if err != nil {
		return JSONdata, err
	}
	return JSONdata, nil
}

// ExportSessionAsStruct - Export a full player session as a struct
func ExportSessionAsStruct(pid int, realm string, days int) (export ExportData, err error) {
	session, playerProfile,  err := calcSession(pid , realm, days)
	if err != nil {
		return export, err
	}
	export.Vehicles 		= session.Vehicles
	export.SessionStats 	= session.StatsAll
	export.PlayerDetails	= playerProfile

	return export, nil
}
package dataprep

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"sync"
	"sync/atomic"

	"math"
	"time"

	"github.com/cufee/am-stats/dataprep"
	db "github.com/cufee/am-stats/mongodbapi"
	wgapi "github.com/cufee/am-stats/wargamingapi"
	"go.mongodb.org/mongo-driver/bson"
)

// CalcVehicleWN8 - Calculate WN8 for a VehicleStats struct
func calcVehicleWN8(tank wgapi.VehicleStats) (wgapi.VehicleStats, error) {
	// Get tank info
	tankInfo, err := db.GetTankGlossary(tank.TankID)
	tank.TankTier = tankInfo.Tier
	tank.TankName = tankInfo.Name

	if err != nil || tankInfo.Name == "" {
		// Refresh Glossary cache -  Disabled, cache updates are done every 24 hours
		// go refreshGlossary()

		log.Print("no tank glossary data (", err, ")")
		tank.TankTier = 0
		tank.TankName = "Unknown"
	}
	// Get tank averages
	tankAvgData, err := db.GetTankAverages(tank.TankID)
	if err != nil {
		log.Print("no tank avg data, but name and tier found")
		tank.TankRawWN8 = 0
		tank.TankWN8 = -1
		return tank, nil
	}
	battles := tank.Battles
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
	rawRating := rating * battles

	tank.TankRawWN8 = rawRating
	tank.TankWN8 = rawRating / battles

	return tank, nil
}

// SliceDiff - Calculate the difference in two VehicleStats slices
func sessionDiff(oldStats db.Session, liveStats db.Session) (session db.Session) {
	// Convert to RetroSession
	var sessionConv db.Convert = oldStats
	retroSession := sessionConv.ToRetro()

	vahiclesChan := make(chan wgapi.VehicleStats, len(liveStats.Vehicles))
	var wg sync.WaitGroup
	var totalRawRating uint64
	var totalRawBattles uint64

	for _, newData := range liveStats.Vehicles {
		if newData.Battles == retroSession.Vehicles[strconv.Itoa(newData.TankID)].Battles {
			// Skop if no battles were played
			continue
		}
		// Start go routines
		wg.Add(1)
		go func(newData wgapi.VehicleStats) {
			defer wg.Done()
			// Get session diff and add vehicle WN8
			finalVehicle, err := calcVehicleWN8(wgapi.Diff(retroSession.Vehicles[strconv.Itoa(newData.TankID)], newData))
			if err != nil {
				log.Println(err)
				return
			}
			// Add raw WN8 to total
			if finalVehicle.TankWN8 > -1 {
				atomic.AddUint64(&totalRawRating, uint64(finalVehicle.TankRawWN8))
				atomic.AddUint64(&totalRawBattles, uint64(finalVehicle.Battles))
			}
			vahiclesChan <- finalVehicle
		}(newData)
	}
	wg.Wait()
	close(vahiclesChan)

	for v := range vahiclesChan {
		session.Vehicles = append(session.Vehicles, v)
	}
	session.Achievements = liveStats.Achievements.Diff(oldStats.Achievements)
	session.PlayerID = liveStats.PlayerID
	session.LastBattle = liveStats.LastBattle
	session.StatsAll = wgapi.FrameDiff(oldStats.StatsAll, liveStats.StatsAll)
	session.StatsRating = wgapi.FrameDiff(oldStats.StatsRating, liveStats.StatsRating)
	session.BattlesAll = session.StatsAll.Battles
	session.BattlesRating = session.StatsRating.Battles
	session.Timestamp = oldStats.Timestamp
	session.SessionRating = -1
	if totalRawBattles > 0 {
		session.SessionRating = int(totalRawRating) / int(totalRawBattles)
	}

	return session
}

// CalcSession - Calculate a new session
func calcSession(pid int, realm string, days int) (session db.Session, oldSession db.Session, playerProfile wgapi.PlayerProfile, err error) {
	// Get live profile
	playerProfile, err = wgapi.PlayerProfileData(pid, realm)
	if err != nil {
		return session, oldSession, playerProfile, err
	}

	// Get live achievements
	liveAchievements, err := wgapi.PlayerAchievements(pid, realm)
	if err != nil {
		return session, oldSession, playerProfile, err
	}

	// Get cached profile
	newCache := convWGtoDBprofile(playerProfile)
	cachedPlayerProfile, err := db.GetPlayerProfile(pid)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			newCache.CareerWN8 = -1
			newCache.Realm = strings.ToUpper(realm)
			err = db.AddPlayer(newCache)
		}
		if err != nil {
			return session, oldSession, playerProfile, err
		}
	}

	// Update profile cache
	newCache.CareerWN8 = cachedPlayerProfile.CareerWN8
	newCache.Realm = strings.ToUpper(realm)

	// Fix WN8
	if cachedPlayerProfile.CareerWN8 == 0 {
		newCache.CareerWN8 = -1
	}

	// Commit update
	_, err = db.UpdatePlayer(bson.M{"_id": playerProfile.ID}, newCache)
	if err != nil {
		log.Printf("Failed to update player profile cache for %v, error: %s", playerProfile.ID, err.Error())
	}

	playerProfile.CareerWN8 = cachedPlayerProfile.CareerWN8
	playerVehicles, err := wgapi.PlayerVehicleStats(pid, realm)
	if err != nil {
		return session, oldSession, playerProfile, err
	}
	// Get previous session
	oldSession, err = db.GetPlayerSession(pid, days, playerProfile.Stats.All.Battles)
	if err != nil {
		if err.Error() == "mongo: no documents in result" && days == 0 {
			// Check if session exists
			s, _ := db.GetSession(bson.M{"player_id": pid})
			// Add a new session if one does not exist
			if s.PlayerID == 0 {
				sessionData := dataprep.LiveToSession(playerProfile, playerVehicles, liveAchievements)
				sessionData.SessionRating = -1
				err = db.AddSession(sessionData)
				if err == nil {
					err = fmt.Errorf("stats: new player, started tracking")
				}
			}
		}
		return session, oldSession, playerProfile, err
	}

	// Calculate session differance and return
	return sessionDiff(oldSession, dataprep.LiveToSession(playerProfile, playerVehicles, liveAchievements)), oldSession, playerProfile, nil
}

// ExportSessionAsStruct - Export a full player session as a struct
func ExportSessionAsStruct(pid int, realm string, days int) (export ExportData, err error) {
	timerStart := time.Now()
	session, lastSession, playerProfile, err := calcSession(pid, realm, days)
	if err != nil {
		return export, err
	}
	lastRetro := lastSession.ToRetro()
	vehicleMap := make(map[string]wgapi.VehicleStats)
	for _, v := range session.Vehicles {
		vehicleMap[strconv.Itoa(v.TankID)] = lastRetro.Vehicles[strconv.Itoa(v.TankID)]
	}
	export.PlayerDetails = playerProfile
	export.SessionStats = session
	export.LastSession.Vehicles = vehicleMap
	export.TimeToComplete = time.Now().Sub(timerStart).Seconds()

	return export, nil
}

func convWGtoDBprofile(wgData wgapi.PlayerProfile) (dbData db.DBPlayerPofile) {
	dbData.ID = wgData.ID
	dbData.LastBattle = time.Unix(int64(wgData.LastBattle), 0)
	dbData.Nickname = wgData.Name
	dbData.ClanID = wgData.ClanID
	dbData.ClanName = wgData.ClanName
	dbData.ClanTag = wgData.ClanTag
	return dbData
}
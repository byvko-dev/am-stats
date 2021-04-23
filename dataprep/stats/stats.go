package dataprep

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"sync"
	"sync/atomic"

	"math"
	"time"

	dbGlossary "github.com/cufee/am-stats/mongodbapi/v1/glossary"
	dbPlayers "github.com/cufee/am-stats/mongodbapi/v1/players"
	dbStats "github.com/cufee/am-stats/mongodbapi/v1/stats"
	"go.mongodb.org/mongo-driver/bson"

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

// CalcVehicleWN8 - Calculate WN8 for a VehicleStats struct
func CalcVehicleWN8(tank wgapi.VehicleStats) (wgapi.VehicleStats, error) {
	// Get tank info
	tankInfo, err := dbGlossary.GetTankGlossary(tank.TankID)
	tank.TankTier = tankInfo.Tier
	tank.TankName = tankInfo.Name

	if err != nil || tankInfo.Name == "" {
		log.Print("no tank glossary data (", err, ")")
		tank.TankTier = 0
		tank.TankName = fmt.Sprintf("Unknown (%v)", tank.TankID)
	}

	// Get tank averages
	tankAvgData, err := dbGlossary.GetTankAverages(tank.TankID)
	if err != nil {
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
func sessionDiff(oldStats dbStats.Session, liveStats dbStats.Session) (session dbStats.Session) {
	// Convert to RetroSession
	var sessionConv dbStats.Convert = oldStats
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
			finalVehicle, err := CalcVehicleWN8(wgapi.Diff(retroSession.Vehicles[strconv.Itoa(newData.TankID)], newData))
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
func calcSession(pid int, tankID int, realm string, days int, special, includeRating bool) (session dbStats.Session, oldSession dbStats.Session, playerProfile wgapi.PlayerProfile, err error) {
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
	cachedPlayerProfile, err := dbPlayers.GetPlayerProfile(pid)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			newCache.CareerWN8 = -1
			newCache.Realm = strings.ToUpper(realm)
			err = dbPlayers.AddPlayer(newCache)
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
	_, err = dbPlayers.UpdatePlayer(bson.M{"_id": playerProfile.ID}, newCache)
	if err != nil {
		log.Printf("Failed to update player profile cache for %v, error: %s", playerProfile.ID, err.Error())
	}

	playerProfile.CareerWN8 = cachedPlayerProfile.CareerWN8
	playerVehicles, err := wgapi.PlayerVehicleStats(pid, tankID, realm)
	if err != nil {
		return session, oldSession, playerProfile, err
	}

	// -1 rating will never find anything valid
	var ratingBattles int = -1
	if includeRating {
		ratingBattles = playerProfile.Stats.Rating.Battles
	}

	// Get previous session
	switch special {
	case true:
		oldSession, err = dbStats.GetPlayerSpecialSession(pid, playerProfile.Stats.All.Battles, ratingBattles)
		if err == nil {
			break
		}
		fallthrough
	default:
		oldSession, err = dbStats.GetPlayerSession(pid, days, playerProfile.Stats.All.Battles, ratingBattles)
	}
	// Check errors
	if err != nil {
		if err.Error() == "mongo: no documents in result" && days == 0 {
			// Check if session exists
			s, _ := dbStats.GetSession(bson.M{"player_id": pid})
			// Add a new session if one does not exist
			if s.PlayerID == 0 {
				sessionData := LiveToSession(playerProfile, playerVehicles, liveAchievements)
				sessionData.SessionRating = -1
				err = dbStats.AddSession(sessionData)
				if err == nil {
					err = fmt.Errorf("stats: new player, started tracking")
				}
			}
		}
		return session, oldSession, playerProfile, err
	}

	// Calculate session differance and return
	return sessionDiff(oldSession, LiveToSession(playerProfile, playerVehicles, liveAchievements)), oldSession, playerProfile, nil
}

// ExportSessionAsStruct - Export a full player session as a struct
func ExportSessionAsStruct(pid int, tankID int, realm string, days int, limit int, sort string, special, includeRating bool) (export ExportData, err error) {
	timerStart := time.Now()
	session, lastSession, playerProfile, err := calcSession(pid, tankID, realm, days, special, includeRating)
	if err != nil {
		return export, err
	}
	lastRetro := lastSession.ToRetro()

	// Sort
	if sort != "" {
		session.Vehicles = SortTanks(session.Vehicles, sort)
	}

	// Limit
	if limit > 0 && len(session.Vehicles) > limit {
		limitedLastSession := make(map[string]wgapi.VehicleStats)
		session.Vehicles = session.Vehicles[0:limit]
		for _, v := range session.Vehicles {
			limitedLastSession[strconv.Itoa(v.TankID)] = lastRetro.Vehicles[strconv.Itoa(v.TankID)]
		}
		lastRetro.Vehicles = limitedLastSession
	}

	export.PlayerDetails = playerProfile
	export.PlayerDetails.Realm = realm
	export.SessionStats = session
	export.LastSession = lastRetro
	export.TimeToComplete = time.Now().Sub(timerStart).Seconds()

	return export, nil
}

func convWGtoDBprofile(wgData wgapi.PlayerProfile) (dbData dbPlayers.DBPlayerPofile) {
	dbData.ID = wgData.ID
	dbData.LastBattle = time.Unix(int64(wgData.LastBattle), 0)
	dbData.Nickname = wgData.Name
	dbData.ClanID = wgData.ClanID
	dbData.ClanName = wgData.ClanName
	dbData.ClanTag = wgData.ClanTag
	return dbData
}

// SortTanks - Sorting of vehicles
func SortTanks(vehicles []wgapi.VehicleStats, sortKey string) []wgapi.VehicleStats {
	// Sort based on passed key
	switch sortKey {
	case "+battles":
		sort.Slice(vehicles, func(i, j int) bool {
			return vehicles[i].Battles < vehicles[j].Battles
		})
	case "-battles":
		sort.Slice(vehicles, func(i, j int) bool {
			return vehicles[i].Battles > vehicles[j].Battles
		})
	case "+winrate":
		sort.Slice(vehicles, func(i, j int) bool {
			return (float64(vehicles[i].Wins) / float64(vehicles[i].Battles)) < (float64(vehicles[j].Wins) / float64(vehicles[j].Battles))
		})
	case "-winrate":
		sort.Slice(vehicles, func(i, j int) bool {
			return (float64(vehicles[i].Wins) / float64(vehicles[i].Battles)) > (float64(vehicles[j].Wins) / float64(vehicles[j].Battles))
		})
	case "+wn8":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].TankWN8) < absInt(vehicles[j].TankWN8)
		})
	case "-wn8":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].TankWN8) > absInt(vehicles[j].TankWN8)
		})
	case "+last_battle":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].LastBattleTime) < absInt(vehicles[j].LastBattleTime)
		})
	case "-last_battle":
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].LastBattleTime) > absInt(vehicles[j].LastBattleTime)
		})
	case "+damage":
		sort.Slice(vehicles, func(i, j int) bool {
			return int(float64(vehicles[i].DamageDealt)/float64(vehicles[i].Battles)) < int(float64(vehicles[j].DamageDealt)/float64(vehicles[j].Battles))
		})
	case "-damage":
		sort.Slice(vehicles, func(i, j int) bool {
			return int(float64(vehicles[i].DamageDealt)/float64(vehicles[i].Battles)) > int(float64(vehicles[j].DamageDealt)/float64(vehicles[j].Battles))
		})
	case "relevance":
		sort.Slice(vehicles, func(i, j int) bool {
			return (absInt(vehicles[i].TankRawWN8) * vehicles[i].LastBattleTime * vehicles[i].Battles) > (absInt(vehicles[j].TankRawWN8) * vehicles[j].LastBattleTime * vehicles[j].Battles)
		})
	default:
		sort.Slice(vehicles, func(i, j int) bool {
			return absInt(vehicles[i].LastBattleTime) > absInt(vehicles[j].LastBattleTime)
		})
	}
	return vehicles
}

// absInt - Absolute value of an integer
func absInt(val int) int {
	if val >= 0 {
		return val
	}
	return -val
}

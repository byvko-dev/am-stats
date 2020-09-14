package mongodbapi

import (
	"strconv"
	"time"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)


// Session - Will be switching to this format soon
type Session struct {
	Vehicles		[]wgapi.VehicleStats	`json:"vehicles" bson:"vehicles"`
	PlayerID		int						`json:"player_id" bson:"player_id"`
	Timestamp		time.Time				`json:"timestamp" bson:"timestamp"`
	LastBattle		time.Time				`json:"last_battle_time" bson:"last_battle_time"`
	BattlesAll		int						`json:"battles_random" bson:"battles_random"`
	StatsAll		wgapi.StatsFrame		`json:"stats_random" bson:"stats_random"`
	BattlesRating	int						`json:"battles_rating" bson:"battles_rating"`
	StatsRating		wgapi.StatsFrame		`json:"stats_rating" bson:"stats_rating"`
	Convert
}

// RetroSession - Session using old data structure
type RetroSession struct {
	Vehicles		map[string]wgapi.VehicleStats	`json:"vehicles" bson:"vehicles"`
	PlayerID		int								`json:"player_id" bson:"player_id"`
	Timestamp		time.Time						`json:"timestamp" bson:"timestamp"`
	LastBattle		time.Time						`json:"last_battle_time" bson:"last_battle_time"`
	BattlesAll		int								`json:"battles_random" bson:"battles_random"`
	StatsAll		wgapi.StatsFrame				`json:"stats_random" bson:"stats_random"`
	BattlesRating	int								`json:"battles_rating" bson:"battles_rating"`
	StatsRating		wgapi.StatsFrame				`json:"stats_rating" bson:"stats_rating"`
	Convert
}
// Convert - Convert between Session and RetroSession
type Convert interface {
	ToSession() Session
	ToRetro() 	RetroSession
}
// ToSession - Covert RetroSession to Session Struct, Session is easier to work with in Go
func (s RetroSession) ToSession() (sessionNew Session) {
	sessionNew.PlayerID			= s.PlayerID
	sessionNew.Timestamp		= s.Timestamp
	sessionNew.LastBattle		= s.LastBattle
	sessionNew.BattlesAll		= s.BattlesAll
	sessionNew.StatsAll			= s.StatsAll
	sessionNew.BattlesRating	= s.BattlesRating
	sessionNew.StatsRating		= s.StatsRating
	// Convert Vehicle Stats
	for _, v := range s.Vehicles {
		sessionNew.Vehicles = append(sessionNew.Vehicles, v)
	}
	return sessionNew
}
// ToRetro - Covert RetroSession to Session Struct, RetroSession is the format used by Aftermath rendering script.
func (s Session) ToRetro() (sessionNew RetroSession) {
	sessionNew.PlayerID			= s.PlayerID
	sessionNew.Timestamp		= s.Timestamp
	sessionNew.LastBattle		= s.LastBattle
	sessionNew.BattlesAll		= s.BattlesAll
	sessionNew.StatsAll			= s.StatsAll
	sessionNew.BattlesRating	= s.BattlesRating
	sessionNew.StatsRating		= s.StatsRating
	// Convert Vehicle Stats
	vehicleMap := make(map[string]wgapi.VehicleStats)
	for _, v := range s.Vehicles {
		vehicleMap[strconv.Itoa(v.TankID)] = v
	}
	sessionNew.Vehicles = vehicleMap
	return sessionNew
}


// DBPlayerPofile - Player data db entry struct
type DBPlayerPofile struct {
	ID				int			`json:"player_id" bson:"_id"`
	ClanID			int			`json:"clan_id" bson:"clan_id"`
	ClanName		int			`json:"clan_name" bson:"clan_name"`
	ClanRole		int			`json:"clan_role" bson:"clan_role"`
	ClanTag			int			`json:"clan_tag" bson:"clan_tag"`
	LastBattle		time.Time	`json:"last_battle_time" bson:"last_battle_time"`
	Nickname		int			`json:"nickname" bson:"nickname"`
	CareerWN8		int			`json:"career_wn8" bson:"career_wn8"`
}

// FilterPair - Used to make BSON filters
type  FilterPair struct {
	Key 	string
	Value 	interface{}
}


// TankAverages - Averages data for a tank
type TankAverages struct {
	All struct {
		Battles              float64 `bson:"battles,omitempty"`
		DroppedCapturePoints float64 `bson:"dropped_capture_points,omitempty"`
	} 								 `bson:"all"`
	Special struct {
		Winrate         float64 `bson:"winrate,omitempty"`
		DamageRatio     float64 `bson:"damageRatio,omitempty"`
		Kdr             float64 `bson:"kdr,omitempty"`
		DamagePerBattle float64 `bson:"damagePerBattle,omitempty"`
		KillsPerBattle  float64 `bson:"killsPerBattle,omitempty"`
		HitsPerBattle   float64 `bson:"hitsPerBattle,omitempty"`
		SpotsPerBattle  float64 `bson:"spotsPerBattle,omitempty"`
		Wpm             float64 `bson:"wpm,omitempty"`
		Dpm             float64 `bson:"dpm,omitempty"`
		Kpm             float64 `bson:"kpm,omitempty"`
		HitRate         float64 `bson:"hitRate,omitempty"`
		SurvivalRate    float64 `bson:"survivalRate,omitempty"`
	} 							`bson:"special"`
	Name   string `bson:"name"`
	Tier   int    `bson:"tier"`
	Nation string `bson:"nation"`
}
package mongodbapi

import (
	"strconv"
	"time"

	wgapi "github.com/cufee/am-stats/wargamingapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlayerStreak - Player win streak data from DB
type PlayerStreak struct {
	PlayerID   *int      `bson:"_id" json:"_id"`
	Battles    *int      `bson:"battles" json:"battles"`
	Losses     *int      `bson:"losses" json:"losses"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
	Streak     int       `bson:"streak" json:"streak"`
	BestStreak int       `bson:"best_streak" json:"best_streak"`
	MinStreak  int       `bson:"min_streak" json:"min_streak"`
	MaxStreak  int       `bson:"max_streak" json:"max_streak"`
}

// Session - Will be switching to this format soon
type Session struct {
	Vehicles      []wgapi.VehicleStats    `json:"vehicles" bson:"vehicles"`
	Achievements  wgapi.AchievementsFrame `json:"achievements" bson:"achievements"`
	PlayerID      int                     `json:"player_id" bson:"player_id"`
	Timestamp     time.Time               `json:"timestamp" bson:"timestamp"`
	LastBattle    time.Time               `json:"last_battle_time" bson:"last_battle_time"`
	BattlesAll    int                     `json:"battles_random" bson:"battles_random"`
	StatsAll      wgapi.StatsFrame        `json:"stats_random" bson:"stats_random"`
	BattlesRating int                     `json:"battles_rating" bson:"battles_rating"`
	StatsRating   wgapi.StatsFrame        `json:"stats_rating" bson:"stats_rating"`
	SessionRating int                     `json:"session_wn8" bson:"session_wn8"`
	Convert       `json:"-" bson:"-"`
}

// RetroSession - Session using old data structure
type RetroSession struct {
	ID            primitive.ObjectID            `json:"_id" bson:"_id"`
	Vehicles      map[string]wgapi.VehicleStats `json:"vehicles" bson:"vehicles"`
	Achievements  wgapi.AchievementsFrame       `json:"achievements" bson:"achievements"`
	PlayerID      int                           `json:"player_id" bson:"player_id"`
	Timestamp     time.Time                     `json:"timestamp" bson:"timestamp"`
	LastBattle    time.Time                     `json:"last_battle_time" bson:"last_battle_time"`
	BattlesAll    int                           `json:"battles_random" bson:"battles_random"`
	StatsAll      wgapi.StatsFrame              `json:"stats_random" bson:"stats_random"`
	BattlesRating int                           `json:"battles_rating" bson:"battles_rating"`
	StatsRating   wgapi.StatsFrame              `json:"stats_rating" bson:"stats_rating"`
	SessionRating int                           `json:"session_wn8" bson:"session_wn8"`
	Convert       `json:"-" bson:"-"`
}

// Convert - Convert between Session and RetroSession
type Convert interface {
	ToSession() Session
	ToRetro() RetroSession
}

// ToSession - Covert RetroSession to Session Struct, Session is easier to work with in Go
func (s RetroSession) ToSession() (sessionNew Session) {
	sessionNew.Achievements = s.Achievements
	sessionNew.PlayerID = s.PlayerID
	sessionNew.Timestamp = s.Timestamp
	sessionNew.LastBattle = s.LastBattle
	sessionNew.BattlesAll = s.BattlesAll
	sessionNew.StatsAll = s.StatsAll
	sessionNew.BattlesRating = s.BattlesRating
	sessionNew.StatsRating = s.StatsRating
	sessionNew.SessionRating = s.SessionRating
	// Convert Vehicle Stats
	for _, v := range s.Vehicles {
		sessionNew.Vehicles = append(sessionNew.Vehicles, v)
	}
	return sessionNew
}

// ToRetro - Covert RetroSession to Session Struct, RetroSession is the format used by Aftermath rendering script.
func (s Session) ToRetro() (sessionNew RetroSession) {
	sessionNew.Achievements = s.Achievements
	sessionNew.PlayerID = s.PlayerID
	sessionNew.Timestamp = s.Timestamp
	sessionNew.LastBattle = s.LastBattle
	sessionNew.BattlesAll = s.BattlesAll
	sessionNew.StatsAll = s.StatsAll
	sessionNew.BattlesRating = s.BattlesRating
	sessionNew.StatsRating = s.StatsRating
	sessionNew.SessionRating = s.SessionRating
	// Convert Vehicle Stats
	vehicleMap := make(map[string]wgapi.VehicleStats)
	for _, v := range s.Vehicles {
		vehicleMap[strconv.Itoa(v.TankID)] = v
	}
	sessionNew.Vehicles = vehicleMap
	return sessionNew
}

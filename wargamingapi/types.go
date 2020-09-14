package externalapis


// StatsFrame - Stats frame struct to unpack json and bson
type StatsFrame struct {
	Spotted              int `json:"spotted" bson:"spotted"`
	Hits                 int `json:"hits" bson:"hits"`
	Frags                int `json:"frags" bson:"frags"`
	MaxXp                int `json:"max_xp" bson:"max_xp"`
	Wins                 int `json:"wins" bson:"wins"`
	Losses               int `json:"losses" bson:"losses"`
	CapturePoints        int `json:"capture_points" bson:"capture_points"`
	Battles              int `json:"battles" bson:"battles"`
	DamageDealt          int `json:"damage_dealt" bson:"damage_dealt"`
	DamageReceived       int `json:"damage_received" bson:"damage_received"`
	MaxFrags             int `json:"max_frags" bson:"max_frags"`
	Shots                int `json:"shots" bson:"shots"`
	Xp                   int `json:"xp" bson:"xp"`
	SurvivedBattles      int `json:"survived_battles" bson:"survived_battles"`
	DroppedCapturePoints int `json:"dropped_capture_points" bson:"dropped_capture_points"`
}
// Vehicle Stats
// dataToPIDres - JSON response from WG API
type vehiclesDataToPIDres struct {
	Data map[string][]VehicleStats `json:"data"`
}
// VehicleStats - Player Vehicle stats struct, used to return final data
type VehicleStats struct {
	StatsFrame				`json:"all" bson:"all"`
	LastBattleTime int 		`json:"last_battle_time" bson:"last_battle_time"`
	MarkOfMastery  int 		`json:"mark_of_mastery" bson:"mark_of_mastery"`
	TankID         int		`json:"tank_id" bson:"tank_id"`
	TankWN8        int		`json:"tank_wn8,omitempty" bson:"tank_wn8,omitempty"`
	TankRawWN8     int
}
// Diff - Calculate the difference in two VehicleStats structs
func Diff(oldStats VehicleStats, newStats VehicleStats) (diffStats VehicleStats) {
	// Set stats fields
	diffStats.StatsFrame = FrameDiff(oldStats.StatsFrame, newStats.StatsFrame)
	// Set other fields
	diffStats.LastBattleTime	= newStats.LastBattleTime
	diffStats.TankID			= newStats.TankID
	diffStats.TankWN8			= 0
	diffStats.TankRawWN8		= 0
	return diffStats
}
// FrameDiff - Calculate the difference in two StatsFrame structs
func FrameDiff(oldStats StatsFrame, newStats StatsFrame) (diffStats StatsFrame) {
	// Set stats fields
	diffStats.Spotted 				= newStats.Spotted - oldStats.Spotted
	diffStats.Hits 					= newStats.Hits - oldStats.Hits
	diffStats.Frags 				= newStats.Frags - oldStats.Frags
	diffStats.MaxXp 				= newStats.MaxXp - oldStats.MaxXp
	diffStats.Wins 					= newStats.Wins - oldStats.Wins
	diffStats.Losses 				= newStats.Losses - oldStats.Losses
	diffStats.CapturePoints			= newStats.CapturePoints - oldStats.CapturePoints
	diffStats.Battles 				= newStats.Battles - oldStats.Battles
	diffStats.DamageDealt 			= newStats.DamageDealt - oldStats.DamageDealt
	diffStats.DamageReceived 		= newStats.DamageReceived - oldStats.DamageReceived
	diffStats.MaxFrags 				= newStats.MaxFrags - oldStats.MaxFrags
	diffStats.Shots 				= newStats.Shots - oldStats.Shots
	diffStats.Xp 					= newStats.Xp - oldStats.Xp
	diffStats.SurvivedBattles 		= newStats.SurvivedBattles - oldStats.SurvivedBattles
	diffStats.DroppedCapturePoints 	= newStats.DroppedCapturePoints - oldStats.DroppedCapturePoints

	return diffStats
}

// Player profile data
// playerDataToPIDres - JSON response from WG API
type playerDataToPIDres struct {
	Data map[string]PlayerProfile `json:"data"`
}
// PlayerProfile - Player profile struct, newer format
type PlayerProfile struct {
	playerClanData				`json:"clan"`
	ID			int				`json:"account_id"`
	Name		string			`json:"nickname"`
	LastBattle	int				`json:"last_battle_time"`
	Stats		playerStatsRes	`json:"statistics"`
}
// Player stats response
type playerStatsRes struct {
	Rating StatsFrame	`json:"rating"`
	All StatsFrame		`json:"all"`
}
// PlayerProfile - 
type playerClanData struct {
	ClanTag		string	`json:"tag,omitempty"`
	ClanName	string	`json:"name,omitempty"`
	ClanID		int		`json:"clan_id,omitempty"`
}


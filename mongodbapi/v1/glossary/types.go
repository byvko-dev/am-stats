package mongodbapi

// TankAverages - Averages data for a tank
type TankAverages struct {
	All struct {
		Battles              float64 `bson:"battles,omitempty"`
		DroppedCapturePoints float64 `bson:"dropped_capture_points,omitempty"`
	} `bson:"all"`
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
	} `bson:"special"`
	Name   string `bson:"name"`
	Tier   int    `bson:"tier"`
	Nation string `bson:"nation"`
}

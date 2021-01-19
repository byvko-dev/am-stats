package handlers

// StatsRequest - Request for stats image
type StatsRequest struct {
	PlayerID  int    `json:"player_id"`
	Premium   bool   `json:"premium"`
	Verified  bool   `json:"verified"`
	Realm     string `json:"realm"`
	Days      int    `json:"days"`
	Sort      string `json:"sort_key"`
	TankLimit int    `json:"detailed_limit"`
	BgURL     string `json:"bg_url"`
}

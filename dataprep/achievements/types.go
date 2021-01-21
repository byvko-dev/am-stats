package dataprep

// MedalWeight - Object for calculating per medal scores
type MedalWeight struct {
	Name    string `json:"medal"`
	Weight  int    `json:"weight"`
	IconURL string `json:"-"`
}

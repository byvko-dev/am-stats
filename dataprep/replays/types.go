package dataprep

import (
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ReplayDetailsRes - Replay details response from WI
type ReplayDetailsRes struct {
	Status string `json:"status"`
	Data   struct {
		ViewURL     string        `json:"view_url"`
		DownloadURL string        `json:"download_url"`
		Summary     ReplaySummary `json:"summary"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// ReplaySummary - Replay summary
type ReplaySummary struct {
	Realm                string             `json:"realm"`
	DownloadURL          string             `json:"download_url"`
	FileURL              string             `json:"file_url"`
	WinnerTeam           int                `json:"winner_team"`
	UploadedBy           int                `json:"uploaded_by"`
	CreditsTotal         int                `json:"credits_total"`
	ExpBase              int                `json:"exp_base"`
	PlayerName           string             `json:"player_name"`
	Title                string             `json:"title"`
	Details              []ReplayPlayerData `json:"details"`
	Vehicle              string             `json:"vehicle"`
	Enemies              []int              `json:"enemies"`
	Description          string             `json:"description"`
	BattleDuration       float64            `json:"battle_duration"`
	ArenaUniqueID        int64              `json:"arena_unique_id"`
	VehicleTier          int                `json:"vehicle_tier"`
	BattleStartTime      string             `json:"battle_start_time"`
	MasteryBadge         int                `json:"mastery_badge"`
	Protagonist          int                `json:"protagonist"`
	BattleType           int                `json:"battle_type"`
	ExpTotal             int                `json:"exp_total"`
	Allies               []int              `json:"allies"`
	VehicleType          int                `json:"vehicle_type"`
	BattleStartTimestamp float64            `json:"battle_start_timestamp"`
	CreditsBase          int                `json:"credits_base"`
	ProtagonistTeam      int                `json:"protagonist_team"`
	MapName              string             `json:"map_name"`
	RoomType             int                `json:"room_type"`
	BattleResult         int                `json:"battle_result"`
}

// ReplayPlayerData - Player performance in replay
type ReplayPlayerData struct {
	DamageAssistedTrack int                 `json:"damage_assisted_track"`
	BaseCapturePoints   int                 `json:"base_capture_points"`
	WpPointsEarned      int                 `json:"wp_points_earned"`
	TimeAlive           int                 `json:"time_alive"`
	ChassisID           int                 `json:"chassis_id"`
	HitsReceived        int                 `json:"hits_received"`
	ShotsSplash         int                 `json:"shots_splash"`
	GunID               int                 `json:"gun_id"`
	HitsPen             int                 `json:"hits_pen"`
	HeroBonusCredits    int                 `json:"hero_bonus_credits"`
	HitpointsLeft       int                 `json:"hitpoints_left"`
	ID                  int                 `json:"dbid"`
	ShotsPen            int                 `json:"shots_pen"`
	ExpForAssist        int                 `json:"exp_for_assist"`
	DamageReceived      int                 `json:"damage_received"`
	HitsBounced         int                 `json:"hits_bounced"`
	HeroBonusExp        int                 `json:"hero_bonus_exp"`
	EnemiesDamaged      int                 `json:"enemies_damaged"`
	Achievements        []ReplayAchievement `json:"achievements"`
	ExpForDamage        int                 `json:"exp_for_damage"`
	DamageBlocked       int                 `json:"damage_blocked"`
	DistanceTravelled   int                 `json:"distance_travelled"`
	HitsSplash          int                 `json:"hits_splash"`
	Credits             int                 `json:"credits"`
	SquadIndex          int                 `json:"squad_index"`
	WpPointsStolen      int                 `json:"wp_points_stolen"`
	DamageMade          int                 `json:"damage_made"`
	VehicleDescr        int                 `json:"vehicle_descr"`
	ExpTeamBonus        int                 `json:"exp_team_bonus"`
	ClanTag             string              `json:"clan_tag"`
	EnemiesSpotted      int                 `json:"enemies_spotted"`
	ShotsHit            int                 `json:"shots_hit"`
	Clanid              int                 `json:"clanid"`
	TurretID            int                 `json:"turret_id"`
	EnemiesDestroyed    int                 `json:"enemies_destroyed"`
	KilledBy            int                 `json:"killed_by"`
	BaseDefendPoints    int                 `json:"base_defend_points"`
	Exp                 int                 `json:"exp"`
	DamageAssisted      int                 `json:"damage_assisted"`
	DeathReason         int                 `json:"death_reason"`
	ShotsMade           int                 `json:"shots_made"`
	Profile             wgapi.PlayerProfile `json:"profile"`
	TankProfile         wgapi.VehicleStats  `json:"tank_profile"`
	Team                int                 `json:"team"`
	IsProtagonist       bool                `json:"is_protagonist"`
}

// ReplayAchievement - Achivement value from replay
type ReplayAchievement struct {
	T int `json:"t"`
	V int `json:"v"`
}

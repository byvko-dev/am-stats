package config

import "os"

// WgAPIAppID - WG Application ID for Wargaming API
var WgAPIAppID string = os.Getenv("WG_APP_ID") // wargmaing application ID

// MongoURI - URI for connecting to MongoDB
var MongoURI string = os.Getenv("MONGO_CONN_STRING") // full mongo uri liek mongodb://user:passwd@host:port

// AMAPIKey - Aftermath Aftermath API key
var AMAPIKey string = os.Getenv("LEGACY_API_KEY") // internal API key, forgot why I put it here but it's required for some things

// UsersAPIURL - API URL for user checks
var UsersAPIURL string = os.Getenv("USERS_API_URL") // https://domain/users/v1

// AssetsPath - Assets path
var AssetsPath string = "assets/"

// OutRPSlimit - Outgoing request limiter for Wargaming API
var OutRPSlimit int = 20 // max is 20 per IP

// WGProxyURL - Proxy for outgoing requests
var WGProxyURL string = os.Getenv("WG_PROXY_URL")

// DefaultBG - default bg image
var DefaultBG string = os.Getenv("DEFAULT_BG_NAME")

// WotInspectorAPI - API endpoint for uploading replays to WoT Inspector
var WotInspectorAPI string = os.Getenv("WOT_INSPECTOR_API")

// If set to true, checking stats for a player will capture sessions for the whole clan
var GreedySessions bool = false

// UsersAPIURL - API URL for user checks
var CacheAPIURL string = os.Getenv("CACHE_API_URL")

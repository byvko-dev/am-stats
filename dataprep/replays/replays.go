package dataprep

import (
	"fmt"
	"net/url"

	"github.com/cufee/am-stats/config"
	"github.com/cufee/am-stats/dataprep"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// ProcessReplay Get replay data from url
func ProcessReplay(replayURL string) (summary ReplaySummary, err error) {
	// Get replay data
	summary, err = getReplayDetails(replayURL)
	if err != nil {
		return summary, err
	}

	// Detect realm
	realm := realmFromID(summary.Protagonist)

	// Get player profile data
	wgPlayerData, err := wgapi.PlayerSliceProfileData(realm, append(summary.Allies, summary.Enemies...))
	if err != nil {
		return summary, err
	}

	// Fill player profiles and team IDs
	var badProfiles int
	for i, player := range summary.Details {
		// Set profile
		player.Profile = wgPlayerData[fmt.Sprint(player.ID)]

		// Check if a profile is blank
		if player.Profile == (wgapi.PlayerProfile{}) {
			badProfiles++
			// Check for bad profiles threshold
			if badProfiles > 3 {
				return summary, fmt.Errorf("replays api error: too many bad profiles")
			}
			continue
		}

		// Set team
		if dataprep.IntInSlice(player.ID, summary.Allies) {
			player.Team = 1
		} else {
			player.Team = 2
		}

		summary.Details[i] = player
	}

	return summary, err
}

func getReplayDetails(replayURL string) (summary ReplaySummary, err error) {
	// Make URL
	requestURL, err := url.Parse(config.WotInspectorAPI + replayURL)
	if err != nil {
		return summary, fmt.Errorf("replays api error: %s", err.Error())
	}

	// Send request
	var relayRes ReplayDetailsRes
	err = dataprep.DecodeHTTPResponse("GET", make(map[string]string), requestURL, nil, &relayRes)
	if err != nil {
		return summary, fmt.Errorf("replays api error: %s", err.Error())
	}

	// Set Download and File URLs
	relayRes.Data.Summary.DownloadURL = relayRes.Data.DownloadURL
	relayRes.Data.Summary.FileURL = replayURL

	return relayRes.Data.Summary, nil
}

func realmFromID(pidInt int) string {
	switch {
	case pidInt < 500000000:
		return "RU"
	case pidInt < 1000000000:
		return "EU"
	case pidInt < 2000000000:
		return "NA"
	default:
		return "ASIA"
	}
}

package external

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/cufee/am-stats/config"
	wgapi "github.com/cufee/am-stats/wargamingapi"
)

// CheckUserByUserID - Check user profile by Discord ID
func GreedyClanPlayerCapture(player wgapi.PlayerProfile, realm string) {
	// Get full clan profile
	clan, err := wgapi.ClanDataByID(player.ClanID, realm)
	if err != nil {
		log.Print("Failed to greedy capture clan data ", err.Error())
		return
	}

	// Convert to strings
	var clanPlayers []string
	for _, m := range clan.MembersIds {
		clanPlayers = append(clanPlayers, strconv.Itoa(m))
	}
	if len(clanPlayers) < 5 {
		return
	}

	// Make URL
	requestURL, err := url.Parse(fmt.Sprintf("%s/%s/players/update?idList=%s", config.CacheAPIURL, realm, strings.Join(clanPlayers, ",")))
	if err != nil {
		log.Print("Failed to parse URL for greedy capture ", err.Error())
		return
	}

	// Make headers
	headers := make(map[string]string)
	headers["x-api-key"] = config.AMAPIKey

	// Send request
	err = DecodeHTTPResponse("GET", headers, requestURL, nil, nil)
	if err != nil {
		log.Print("Failed to parse URL for greedy capture ", err.Error())
		return
	}
}

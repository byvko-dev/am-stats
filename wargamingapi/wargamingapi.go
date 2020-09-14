package externalapis

import (
	"fmt"
	"errors"
	"strconv"
	"strings"

	"time"
	"net/http"
	"encoding/json"

	"github.com/cufee/am-stats/config"
)

// API Base URLs
// Players
var wgAPIVehicles string = fmt.Sprintf("/wotb/tanks/stats/?application_id=%s&account_id=", config.WgAPIAppID)
var wgAPIProfileData string = fmt.Sprintf("/wotb/account/info/?application_id=%s&extra=statistics.rating&account_id=", config.WgAPIAppID)
var wgAPIPlayerClan string = fmt.Sprintf("/wotb/clans/accountinfo/?application_id=%s&extra=clan&account_id=", config.WgAPIAppID)
// Clans
var wgAPIClanInfo string = fmt.Sprintf("/wotb/clans/list/?application_id=%s&search=", config.WgAPIAppID)
var wgAPIClanDetails string = fmt.Sprintf("/wotb/clans/info/?application_id=%s&fields=clan_id,name,tag,is_clan_disbanded,members_ids,updated_at,members&extra=members&clan_id=", config.WgAPIAppID)


// HTTP client
var clientHTTP = &http.Client{Timeout: 10 * time.Second}

// getFlatJSON - 
func getJSON(url string, target interface{}) error {
	res, err := clientHTTP.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v. error: %s", res.StatusCode, err)
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}

// getAPIDomain - Get WG API domain using realm
func getAPIDomain(realm string) (string, error) {
	realm = strings.ToUpper(realm)
	if realm == "NA" {
		return "http://api.wotblitz.com", nil

	} else if realm == "EU" {
		return "http://api.wotblitz.eu", nil

	} else if realm == "RU" {
		return "http://api.wotblitz.ru", nil

	} else if realm == "ASIA" || realm == "AS" {
		return "http://api.wotblitz.ru", nil

	} else {
		message := fmt.Sprintf("Realm %s not found", realm)
		return "", errors.New(message)
	}
}

// PlayerVehicleStats - Fetch stats for player Vehicles, returns a slice of all vehicle stats
func PlayerVehicleStats(playerID int, realm string) ([]VehicleStats, error) {
	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return nil, err
	}
	// Get stats
	url := domain + wgAPIVehicles + strconv.Itoa(playerID)
	var rawResponse vehiclesDataToPIDres
	
	err = getJSON(url, &rawResponse)
	if err != nil {
		return nil, err
	}
	finalResponse := rawResponse.Data[strconv.Itoa(playerID)]
	return finalResponse, nil
}

// PlayerProfileData - Fetch general account information and all stats for a player
func PlayerProfileData(playerID int, realm string) (PlayerProfile , error) {
	var finalResponse PlayerProfile
	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return finalResponse, err
	}
	// Get stats
	url := domain + wgAPIProfileData + strconv.Itoa(playerID)
	var rawResponse playerDataToPIDres
	
	err = getJSON(url, &rawResponse)
	if err != nil {
		return finalResponse, err
	}
	finalResponse = rawResponse.Data[strconv.Itoa(playerID)]
	return finalResponse, nil
}

package externalapis

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cufee/am-stats/config"
)

// API Base URLs
// Players
var wgAPIVehicles string = fmt.Sprintf("/wotb/tanks/stats/?application_id=%s&account_id=", config.WgAPIAppID)
var wgAPIProfileData string = fmt.Sprintf("/wotb/account/info/?application_id=%s&extra=statistics.rating&account_id=", config.WgAPIAppID)
var wgAPIPlayerClan string = fmt.Sprintf("/wotb/clans/accountinfo/?application_id=%s&extra=clan&account_id=", config.WgAPIAppID)
var wgAPIPlayerAchievements string = fmt.Sprintf("/wotb/account/achievements/?application_id=%s&account_id=", config.WgAPIAppID)

// Clans
var wgAPIClanSearch string = fmt.Sprintf("/wotb/clans/list/?application_id=%s&search=", config.WgAPIAppID)
var wgAPIClanDetails string = fmt.Sprintf("/wotb/clans/info/?application_id=%s&fields=clan_id,name,tag,members_ids,members_count&clan_id=", config.WgAPIAppID)

// HTTP client
var clientHTTP = &http.Client{Timeout: 750 * time.Millisecond, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

// Mutex lock for rps counter
var waitGroup sync.WaitGroup
var limiterChan chan int = make(chan int, config.OutRPSlimit)

// getJSON -
func getJSON(url string, target interface{}) error {
	// Outgoing rate limiter
	start := time.Now()
	limiterChan <- 1
	defer func() {
		go func() {
			timer := time.Now().Sub(start)

			if timer < (time.Second * 1) {
				toSleep := (time.Second * 1) - timer
				time.Sleep(toSleep)
			}
			<-limiterChan
		}()
	}()

	var resData []byte
	res, err := clientHTTP.Get(url)

	if res == nil {
		// Change timeout to account for cold starts
		timeout := clientHTTP.Timeout
		clientHTTP.Timeout = 2 * time.Second
		defer func() { clientHTTP.Timeout = timeout }()

		// Marshal a request
		proxyReq := struct {
			URL string `json:"url"`
		}{
			URL: url,
		}
		reqData, pErr := json.Marshal(proxyReq)
		if pErr != nil {
			return pErr
		}

		// Make request
		req, pErr := http.NewRequest("GET", config.WGProxyURL, bytes.NewBuffer(reqData))
		if pErr != nil {
			return fmt.Errorf("proxy: no response recieved, error: %v", pErr)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send request
		res, pErr = clientHTTP.Do(req)
		if res == nil {
			return fmt.Errorf("proxy: no response recieved, error: %v", pErr)
		}
		if pErr != nil {
			return pErr
		}

		// Read body
		resData, pErr = ioutil.ReadAll(res.Body)
		if pErr != nil {
			return pErr
		}

		// Check for errors
		var proxyErr struct {
			Message string `json:"error"`
		}
		json.Unmarshal(resData, &proxyErr)
		if proxyErr.Message != "" {
			pErr = fmt.Errorf(proxyErr.Message)
		}

		// Set error to proxy error
		err = pErr
	} else {
		resData, err = ioutil.ReadAll(res.Body)
	}

	// Check error
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.Unmarshal(resData, target)
}

// getAPIDomain - Get WG API domain using realm
func getAPIDomain(realm string) (string, error) {
	realm = strings.ToUpper(realm)
	if realm == "NA" {
		return "https://api.wotblitz.com", nil

	} else if realm == "EU" {
		return "https://api.wotblitz.eu", nil

	} else if realm == "RU" {
		return "https://api.wotblitz.ru", nil

	} else if realm == "ASIA" || realm == "AS" {
		return "https://api.wotblitz.asia", nil

	} else {
		message := fmt.Sprintf("Realm %s not found", realm)
		return "", errors.New(message)
	}
}

// PlayerVehicleStats - Fetch stats for player Vehicles, returns a slice of all vehicle stats
func PlayerVehicleStats(playerID int, tankID int, realm string) (finalResponse []VehicleStats, err error) {
	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return nil, err
	}
	// Get stats
	url := domain + wgAPIVehicles + strconv.Itoa(playerID)

	if tankID != 0 {
		url += fmt.Sprintf("&tank_id=%v", tankID)
	}

	var rawResponse vehiclesDataToPIDres
	err = getJSON(url, &rawResponse)
	if err != nil {
		return nil, err
	}
	if rawResponse.Error.Message != "" {
		return finalResponse, fmt.Errorf("WG error: %s", rawResponse.Error.Message)
	}
	finalResponse = rawResponse.Data[strconv.Itoa(playerID)]
	if len(finalResponse) < 1 {
		return finalResponse, errors.New("no vehicles data available for player")
	}
	return finalResponse, nil
}

// PlayerAchievements - Fetch achievements for player, returns a struct of all achievements
func PlayerAchievements(playerID int, realm string) (AchievementsFrame, error) {
	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return AchievementsFrame{}, err
	}

	// Get stats
	url := domain + wgAPIPlayerAchievements + strconv.Itoa(playerID)
	var rawResponse vehiclesAchievmentsRes
	err = getJSON(url, &rawResponse)
	if err != nil {
		err = fmt.Errorf("error: " + err.Error() + "\nwg responded with: " + rawResponse.Error.Message)
		return AchievementsFrame{}, err
	}
	if rawResponse.Status != "ok" {
		err = fmt.Errorf("wg responded with: " + rawResponse.Error.Message)
		return AchievementsFrame{}, err
	}
	finalResponse := rawResponse.Data[strconv.Itoa(playerID)]
	return finalResponse, nil
}

// PlayerProfileData - Fetch general account information and all stats for a player
func PlayerProfileData(playerID int, realm string) (finalResponse PlayerProfile, err error) {
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
	if rawResponse.Error.Message != "" {
		return finalResponse, fmt.Errorf("WG error: %s", rawResponse.Error.Message)
	}
	if rawResponse.Status != "ok" {
		return finalResponse, fmt.Errorf("WG error: %v", rawResponse.Error.Message)
	}
	if _, ok := rawResponse.Data[strconv.Itoa(playerID)]; ok == false || rawResponse.Data[strconv.Itoa(playerID)].ID != playerID {
		return finalResponse, fmt.Errorf("WG: player not found in response. status %v", rawResponse.Status)
	}
	finalResponse = rawResponse.Data[strconv.Itoa(playerID)]

	// Get clan data
	var clanRes playerDataToPIDres
	url = domain + wgAPIPlayerClan + strconv.Itoa(playerID)
	err = getJSON(url, &clanRes)
	if err != nil {
		return finalResponse, err
	}
	finalResponse.playerClanData = clanRes.Data[strconv.Itoa(playerID)].playerClanData
	return finalResponse, nil
}

// PlayerSliceProfileData - Fetch general account information and all stats for a player
func PlayerSliceProfileData(realm string, playerIDsRaw []int) (finalResponse map[string]PlayerProfile, err error) {
	// Check list length
	if len(playerIDsRaw) > 100 {
		return finalResponse, fmt.Errorf("player_id list is too long")
	}

	// Conver list to strings
	var playerIDs []string
	for _, pid := range playerIDsRaw {
		playerIDs = append(playerIDs, strconv.Itoa(pid))
	}

	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return finalResponse, err
	}

	// Get stats
	url := domain + wgAPIProfileData + strings.Join(playerIDs, ",")
	var rawResponse playerDataToPIDres

	err = getJSON(url, &rawResponse)
	if err != nil {
		return finalResponse, err
	}
	if rawResponse.Error.Message != "" {
		return finalResponse, fmt.Errorf("WG error: %s", rawResponse.Error.Message)
	}
	if rawResponse.Status != "ok" {
		return finalResponse, fmt.Errorf("WG error: %v", rawResponse.Error.Message)
	}

	// Get clan data
	var clanRes playerDataToPIDres
	url = domain + wgAPIPlayerClan + strings.Join(playerIDs, ",")
	err = getJSON(url, &clanRes)
	if err != nil {
		return finalResponse, err
	}

	finalResponse = rawResponse.Data

	// Fill clan data
	for pid, playerData := range finalResponse {
		playerData.playerClanData = clanRes.Data[pid].playerClanData
		finalResponse[pid] = playerData
	}

	return finalResponse, nil
}

// ClanDataByID - Fetch clan profile by clan ID and realm
func ClanDataByID(clanID int, realm string) (data ClanProfile, err error) {
	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return data, err
	}

	// Get clan Profile
	url := domain + wgAPIClanDetails + fmt.Sprint(clanID)
	var rawResponse clanDetailsRes
	err = getJSON(url, &rawResponse)
	if err != nil {
		err = fmt.Errorf("error: " + err.Error() + "\nwg responded with: " + rawResponse.Error.Message)
		return data, err
	}
	if rawResponse.Status != "ok" {
		err = fmt.Errorf("wg responded with: " + rawResponse.Error.Message)
		return data, err
	}

	return rawResponse.Data[fmt.Sprint(clanID)], err
}

// ClanDataByTag - Fetch clan profile by clan tag and realm
func ClanDataByTag(clanTag string, realm string) (data ClanProfile, err error) {
	// Get API domain
	domain, err := getAPIDomain(realm)
	if err != nil {
		return data, err
	}

	// Search clan by tag
	url := domain + wgAPIClanSearch + clanTag
	var rawResponse clanSearchRes
	err = getJSON(url, &rawResponse)
	if err != nil {
		err = fmt.Errorf("error: " + err.Error() + "\nwg responded with: " + rawResponse.Error.Message)
		return data, err
	}
	if rawResponse.Status != "ok" {
		err = fmt.Errorf("wg responded with: " + rawResponse.Error.Message)
		return data, err
	}

	// Look for tag
	var matchID int
	for _, cl := range rawResponse.Data {
		if cl.Tag == strings.ToUpper(clanTag) {
			matchID = cl.ClanID
			break
		}
	}
	if matchID == 0 {
		return data, fmt.Errorf("clan not found on realm")
	}

	// Get profile
	return ClanDataByID(matchID, realm)
}

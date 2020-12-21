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

// Clans
var wgAPIClanInfo string = fmt.Sprintf("/wotb/clans/list/?application_id=%s&search=", config.WgAPIAppID)
var wgAPIClanDetails string = fmt.Sprintf("/wotb/clans/info/?application_id=%s&fields=clan_id,name,tag,is_clan_disbanded,members_ids,updated_at,members&extra=members&clan_id=", config.WgAPIAppID)

// HTTP client
var clientHTTP = &http.Client{Timeout: 250 * time.Millisecond, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

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
		var clientHTTPlocal = &http.Client{Timeout: 2000 * time.Millisecond, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		// Marshal a request
		proxyReq := struct {
			URL string `json:"url"`
		}{
			URL: url,
		}
		reqData, pErr := json.Marshal(proxyReq)
		if pErr != nil {
			return fmt.Errorf("no response recieved from WG API after proxy try, error: %v", pErr)
		}

		// Make request
		req, pErr := http.NewRequest("GET", config.WGProxyURL, bytes.NewBuffer(reqData))
		if pErr != nil {
			return fmt.Errorf("failed to make a proxy request, error: %v", pErr)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send request
		res, pErr = clientHTTPlocal.Do(req)
		if res == nil {
			return fmt.Errorf("no response recieved from WG API after proxy try, error: %v", pErr)
		}
		resData, pErr = ioutil.ReadAll(res.Body)

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
func PlayerVehicleStats(playerID int, realm string) (finalResponse []VehicleStats, err error) {
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
	if rawResponse.Error.Message != "" {
		return finalResponse, fmt.Errorf("WG error: %s", rawResponse.Error.Message)
	}
	finalResponse = rawResponse.Data[strconv.Itoa(playerID)]
	if len(finalResponse) < 1 {
		return finalResponse, errors.New("no vehicles data available for player")
	}
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

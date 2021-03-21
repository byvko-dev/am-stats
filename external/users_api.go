package external

import (
	"fmt"
	"net/url"

	"github.com/cufee/am-stats/config"
)

// UserCheckResponse - reponse from a user check
type UserCheckResponse struct {
	DefaultPID int    `json:"player_id"`
	Locale     string `json:"locale"`

	Premium  bool `json:"premium"`
	Verified bool `json:"verified"`

	CustomBgURL string `json:"bg_url"`

	Banned      bool   `json:"banned"`
	BanReason   string `json:"ban_reason,omitempty"`
	BanNotified bool   `json:"ban_notified,omitempty"`

	Error string `json:"error"`
}

// CheckUserByUserID - Check user profile by Discord ID
func CheckUserByUserID(userIDStr string) (userData UserCheckResponse, err error) {
	// Make URL
	requestURL, err := url.Parse(fmt.Sprintf("%s/users/id/%s", config.UsersAPIURL, userIDStr))
	if err != nil {
		return userData, fmt.Errorf("users api error: %s", err.Error())
	}

	// Make headers
	headers := make(map[string]string)
	headers["x-api-key"] = config.AMAPIKey

	// Send request
	err = DecodeHTTPResponse("GET", headers, requestURL, nil, &userData)
	if err != nil {
		return userData, fmt.Errorf("users api error: %s", err.Error())
	}

	// Check for returned error
	if userData.Error != "" {
		err = fmt.Errorf("users api error: %s", userData.Error)
	}
	return userData, err
}

// CheckUserByPID - Check user profile by player id
func CheckUserByPID(pid int) (userData UserCheckResponse, err error) {
	// Make URL
	requestURL, err := url.Parse(fmt.Sprintf("%s/players/id/%v", config.UsersAPIURL, pid))
	if err != nil {
		return userData, fmt.Errorf("users api error: %s", err.Error())
	}

	// Make headers
	headers := make(map[string]string)
	headers["x-api-key"] = config.AMAPIKey

	// Send request
	err = DecodeHTTPResponse("GET", headers, requestURL, nil, &userData)
	if err != nil {
		return userData, fmt.Errorf("users api error: %s", err.Error())
	}

	// Check error
	if userData.Error != "" {
		err = fmt.Errorf("users api error: %s", userData.Error)
	}
	return userData, err
}

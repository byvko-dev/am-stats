package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTP client
var clientHTTP = &http.Client{Timeout: 2 * time.Second}

// DecodeHTTPResponse - Send HTTP request with a payload and headers, decode to target
func DecodeHTTPResponse(method string, headers map[string]string, requestURL *url.URL, reqData []byte, target interface{}) error {
	// Make request
	req, err := http.NewRequest(strings.ToUpper(method), requestURL.String(), bytes.NewBuffer(reqData))
	if err != nil {
		return err
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Send request
	res, err := clientHTTP.Do(req)

	// Check for response
	if res == nil {
		return fmt.Errorf("no response recieved, error: %v", err)
	}

	// Read body
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if target != nil {
		// Decode body
		return json.Unmarshal(resData, target)
	} else {
		return err
	}
}

// RawHTTPResponse - Send HTTP request with a payload and headers, return raw body
func RawHTTPResponse(method string, headers map[string]string, requestURL *url.URL, reqData []byte) (data []byte, err error) {
	// Make request
	req, err := http.NewRequest(strings.ToUpper(method), requestURL.String(), bytes.NewBuffer(reqData))
	if err != nil {
		return data, err
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Send request
	res, err := clientHTTP.Do(req)

	// Check for response
	if res == nil {
		return data, fmt.Errorf("no response recieved, error: %v", err)
	}

	// Read body
	resData, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	// Decode body
	return resData, err
}

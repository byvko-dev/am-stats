package dataprep

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type responseError struct {
	Error string `json:"error"`
}

// HTTP client
var clientHTTP = &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

// DecodeHTTPResponse - Send HTTP request with a payload and headers, decode to target
func DecodeHTTPResponse(method string, headers map[string]string, requestURL *url.URL, reqData []byte, target interface{}) error {
	// Make request
	req, err := http.NewRequest(strings.ToUpper(method), requestURL.String(), bytes.NewBuffer(reqData))
	if err != nil {
		return err
	}

	// Set headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
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

	// Decode body
	return json.Unmarshal(resData, target)
}

// StringInSlice - Check if a slice contains a string
func StringInSlice(str string, sl []string) bool {
	for _, s := range sl {
		if s == str {
			return true
		}
	}
	return false
}

// IntInSlice - Check if a slice contains an integer
func IntInSlice(n int, sl []int) bool {
	for _, num := range sl {
		if n == num {
			return true
		}
	}
	return false
}

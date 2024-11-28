package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL = "https://s3-eh1.salford.ac.uk/finestra/api/"
	timeout = 60 * time.Second
)

type SessionResponse struct {
	SessionID string `json:"SessionID"`
}

func authenticate(username *string, password *string) (string, error) {
	url := baseURL + "session/logon"
	base64EncodedPassword := base64.StdEncoding.EncodeToString([]byte(*password))
	payload := fmt.Sprintf(`{"Username":"%s","Password":"%s"}`, *username, base64EncodedPassword)

	body, err := callAPI(url, "POST", payload)
	if err != nil {
		return "", err
	}

	var sessionResponse SessionResponse
	if err := json.Unmarshal(body, &sessionResponse); err != nil {
		return "", err
	}

	return sessionResponse.SessionID, nil
}

type LocationDataResponse struct {
	Data string `json:"Data"`
}

func exportLocationData(apiKey string, locationID *string, date string) (string, error) {
	url := fmt.Sprintf(baseURL+"locations/%s/export/range/csv?customStartDate=%s&customEndDate=%s&api_key=%s", *locationID, date, date, apiKey)
	method := "GET"
	body, err := callAPI(url, method, "")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var responseData LocationDataResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println(err)
		return "", err
	}
	return responseData.Data, nil
}

func callAPI(url, method, payload string) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := createRequest(url, method, payload)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func createRequest(url, method, payload string) (*http.Request, error) {
	var body *strings.Reader
	if payload != "" {
		body = strings.NewReader(payload)
	} else {
		body = strings.NewReader("")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	return req, nil
}

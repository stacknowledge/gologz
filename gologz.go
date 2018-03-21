package gologz

import (
	"encoding/json"
	"strings"
)

const (
	euURL = "https://api-eu.logz.io/v1" // European Accounts URL
	usURL = "https://api.logz.io/v1"    // United States Accounts URL
)

// Logzio client structure with needed data to authenticate and which url to access
type Logzio struct {
	Token   string
	BaseURL string
	fetcher func(parameters *request, logzioURL, token string) ([]byte, error)
}

// New returns a new client using a application token and the service region
func New(token, region string) *Logzio {
	client := Logzio{
		Token:   token,
		BaseURL: euURL,
		fetcher: fetch,
	}

	if strings.ToLower(region) == "us" {
		client.BaseURL = usURL
	}

	return &client
}

// Trace a log based on a string input with a start point and a limit
func (logz *Logzio) Trace(input string, start, limit int) (*LogzResponse, error) {
	traces, err := logz.fetcher(&request{
		Endpoint: "/search",
		Query: queryObject{
			From: start,
			Size: limit,
			QueryType: &queryType{
				String: &queryString{
					input,
				},
			},
		},
	}, logz.BaseURL, logz.Token)

	if err != nil {
		return nil, err
	}

	traceTransactionResponse := new(LogzResponse)
	return traceTransactionResponse, json.Unmarshal(traces, &traceTransactionResponse)
}

// TracePairs traces a log with key values input with a start point and a limit
func (logz *Logzio) TracePairs(input map[string]interface{}, start, limit int) (*LogzResponse, error) {
	arrangedPairs := composePairs(input)
	traces, err := logz.fetcher(&request{
		Endpoint: "/search",
		Query: queryObject{
			From: start,
			Size: limit,
			QueryType: &queryType{
				String: &queryString{
					arrangedPairs,
				},
			},
		},
	}, logz.BaseURL, logz.Token)

	if err != nil {
		return nil, err
	}

	traceTransactionResponse := new(LogzResponse)
	return traceTransactionResponse, json.Unmarshal(traces, &traceTransactionResponse)
}

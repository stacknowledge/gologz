package gologz

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	errFailedToParseURL       = errors.New("error: Failed to parse the desired URL")
	errFailedToComposeRequest = errors.New("error: Failed to compose the desired request")
	errFailedToRequest        = errors.New("error: Failed to request. Connectivity problems?")
	errInvalidResponse        = errors.New("error: Invalid response, status code is not 200")
	errFailedToParseBody      = errors.New("error: Failed to parse response body")
)

func composePairs(input map[string]interface{}) string {
	var output string
	if len(input) == 0 {
		return ""
	}

	for key := range input {
		switch input[key].(type) {
		case string:
			output += fmt.Sprintf("%s:\"%s\"", key, input[key])
		default:
			output += fmt.Sprintf("%s:%v", key, input[key])
		}

		output += " AND "
	}

	return output[:len(output)-5] //return without the last AND
}

func fetch(parameters *request, requestURL, token string) ([]byte, error) {
	url, err := url.Parse(requestURL + parameters.Endpoint)

	if err != nil {
		return nil, errFailedToParseURL
	}

	payload, _ := json.Marshal(parameters.Query)
	request, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, errFailedToComposeRequest
	}

	request.Header.Set("X-USER-TOKEN", token)
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, errFailedToRequest
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errInvalidResponse
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errFailedToParseBody
	}

	return body, nil
}

package gologz

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	usecases := []struct {
		name        string
		token       string
		region      string
		expectedURL string
	}{
		{"Should return a valid client with a valid url for EU region", "X-CLIENT-API-TOKEN", "eu", euURL},
		{"Should return a valid client with a valid url for US region", "X-CLIENT-API-TOKEN", "us", usURL},
	}

	for _, usecase := range usecases {
		t.Run(usecase.name, func(t *testing.T) {
			client := New(usecase.token, usecase.region)
			if ok := reflect.TypeOf(client); ok == nil {
				t.Errorf("unexpected type on client instantiation got %v", ok)
			}

			if client.BaseURL != usecase.expectedURL {
				t.Errorf("unexpected base url. expected %s, got %s", usecase.expectedURL, client.BaseURL)
			}
		})
	}
}

func TestTrace(t *testing.T) {
	usecases := []struct {
		name           string
		token          string
		input          string
		expectedResult int
		expectedError  error
	}{
		{"Should return an error when requesting with an unauthorized token",
			"X-UNAUTHORIZED-TOKEN", "*", 0, errInvalidToken},
		{"Should return an LogzResponse with data with an authorized access and a valid query",
			"X-AUTHORIZED-TOKEN", "test_parameter:test_value", 2, nil},
		{"Should return an empty LogzResponse with an authorized access and valid query but no found records",
			"X-AUTHORIZED-TOKEN", "", 0, nil},
	}

	for _, usecase := range usecases {
		t.Run(usecase.name, func(t *testing.T) {
			client := &Logzio{Token: usecase.token, BaseURL: euURL, fetcher: MockFetcher}
			result, err := client.Trace(usecase.input, 0, 2)

			if err != usecase.expectedError {
				t.Errorf("unexpected error, expected %v, got %v", err, usecase.expectedError)
			}

			if result != nil && len(result.Logs) != usecase.expectedResult {
				t.Errorf("unexpected log count, expected %d, got %d", usecase.expectedResult, len(result.Logs))
			}

			if usecase.expectedResult > 0 {
				for index := range result.Logs {
					if result.Logs[index].Log == "" {
						t.Errorf("Expected log value to be parsed, got empty value")
					}
					if result.Logs[index].Service == "" {
						t.Errorf("Expected service value to be parsed, got empty value")
					}
					if result.Logs[index].When == "" {
						t.Errorf("Expected when value to be parsed, got empty value")
					}
				}
			}
		})
	}
}

func TestTracePairs(t *testing.T) {
	usecases := []struct {
		caseName       string
		token          string
		input          map[string]interface{}
		expectedResult int
		expectedError  error
	}{
		{"Should return an error when requesting with an unauthorized token",
			"X-UNAUTHORIZED-TOKEN", map[string]interface{}{}, 0, errInvalidToken},
		{"Should return an LogzResponse with data with an authorized access and a valid query",
			"X-AUTHORIZED-TOKEN", map[string]interface{}{"test_parameter": "test_value", "test_parameter2": 123}, 2, nil},
		{"Should return an empty LogzResponse with an authorized access and valid query but no found records",
			"X-AUTHORIZED-TOKEN", map[string]interface{}{}, 0, nil},
	}

	for _, usecase := range usecases {
		t.Run(usecase.caseName, func(t *testing.T) {
			client := &Logzio{Token: usecase.token, BaseURL: euURL, fetcher: MockFetcher}
			result, err := client.TracePairs(usecase.input, 0, 2)

			if err != usecase.expectedError {
				t.Errorf("unexpected error, expected %v, got %v", err, usecase.expectedError)
			}

			if result != nil && len(result.Logs) != usecase.expectedResult {
				t.Errorf("unexpected log count, expected %d, got %d", usecase.expectedResult, len(result.Logs))
			}

			if usecase.expectedResult > 0 {
				for index := range result.Logs {
					if result.Logs[index].Log == "" {
						t.Errorf("Expected log value to be parsed, got empty value")
					}
					if result.Logs[index].Service == "" {
						t.Errorf("Expected service value to be parsed, got empty value")
					}
					if result.Logs[index].When == "" {
						t.Errorf("Expected when value to be parsed, got empty value")
					}
				}
			}
		})
	}
}

func MockFetcher(parameters *request, logzioURL, token string) ([]byte, error) {
	var response []byte
	var err error

	switch token {
	case "X-UNAUTHORIZED-TOKEN":
		err = errInvalidToken
	case "X-AUTHORIZED-TOKEN":
		err = nil
	}

	switch parameters.Query.QueryType.String.Input {
	case "test_parameter:test_value":
		response = validResponse
	case "test_parameter:\"test_value\" AND test_parameter2:123":
		response = validResponse
	default:
		response = invalidResponse
	}

	return response, err
}

var (
	errInvalidToken = errors.New("Invalid token")

	invalidResponse = []byte(`{}`)
	validResponse   = []byte(`{
		"hits": {
			"total": 17,
			"hits": [
				{
					"_source": {
						"log": "{\"this_is\" : \"a_log\"}",
						"service": "test-service",
						"when": "2018-03-19T16:15:30.291+0000"
					}
				},
				{
					"_source": {
						"log": "{\"this_is\" : \"a_log\"}",
						"service": "test-service",
						"when": "2018-03-19T16:15:30.291+0000"
					}
				}
			]
		}
	}`)
)

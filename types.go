package gologz

import (
	"encoding/json"
)

type request struct {
	Endpoint string
	Query    queryObject
}

type queryObject struct {
	From      int        `json:"from"`
	Size      int        `json:"size"`
	QueryType *queryType `json:"query"`
}

type queryType struct {
	String *queryString `json:"query_string"`
}

type queryString struct {
	Input string `json:"query"`
}

type logBag struct {
	When    string `json:"when"`
	Service string `json:"service"`
	Log     string `json:"log"`
}

type LogzResponse struct {
	Count float64
	Logs  []logBag `json:"hits"`
}

func (lr *LogzResponse) UnmarshalJSON(bytes []byte) error {
	var whatever interface{}
	json.Unmarshal(bytes, &whatever)

	marshalMap := whatever.(map[string]interface{})
	parent := marshalMap["hits"]

	if child, ok := parent.(map[string]interface{}); ok {
		lr.Count = child["total"].(float64)

		for _, value := range child["hits"].([]interface{}) {
			mapIndex := value.(map[string]interface{})
			source := mapIndex["_source"].(map[string]interface{})

			var log, when, service string
			if value, ok := source["log"].(string); ok {
				log = value
			}
			if value, ok := source["when"].(string); ok {
				when = value
			}
			if value, ok := source["service"].(string); ok {
				service = value
			}

			lr.Logs = append(lr.Logs, logBag{
				Log:     log,
				When:    when,
				Service: service,
			})
		}
	}

	return nil
}

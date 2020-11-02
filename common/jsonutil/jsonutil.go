package jsonutil

import "encoding/json"

func ToJSON(s string) map[string]interface{} {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		panic(err)
	}
	return raw
}

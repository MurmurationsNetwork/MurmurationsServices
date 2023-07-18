package jsonutil

import "encoding/json"

func ToJSON(s string) map[string]interface{} {
	var raw map[string]interface{}
	err := json.Unmarshal([]byte(s), &raw)
	if err != nil {
		return map[string]interface{}{}
	}
	return raw
}

func ToString(m map[string]interface{}) string {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(jsonString)
}

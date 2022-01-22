package countries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func FindAlpha2ByName(country interface{}) (countryCode string, err error) {
	f, err := os.Open("countries.json")
	if err != nil {
		return "undefined", err
	}

	file, err := ioutil.ReadAll(f)
	if err != nil {
		return "undefined", err
	}

	var countryNames map[string][]string

	err = json.Unmarshal([]byte(file), &countryNames)

	if err != nil {
		return "undefined", err
	}

	countryStr := fmt.Sprintf("%v", country)
	countryLowerStr := strings.ToLower(countryStr)

	for countryCode, countryName := range countryNames {
		for _, alias := range countryName {
			if countryLowerStr == alias {
				return countryCode, nil
			}
		}
	}

	return "undefined", err
}

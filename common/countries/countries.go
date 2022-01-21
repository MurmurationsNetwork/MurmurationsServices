package countries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func FindAlpha2ByName(country interface{}) (countryCode string, err error) {
	file, _ := ioutil.ReadFile("./common/countries/countries.json")

	var countryNames map[string][]string

	err = json.Unmarshal([]byte(file), &countryNames)

	if err != nil {
		return "undefined", err
	}

	countryStr := fmt.Sprintf("%v", country)

	for countryCode, countryName := range countryNames {
		for _, alias := range countryName {
			if countryStr == alias {
				return countryCode, nil
			}
		}
	}

	return "undefined", err
}

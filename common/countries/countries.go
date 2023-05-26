package countries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FindAlpha2ByName(
	countryURL string,
	country interface{},
) (countryCode string, err error) {
	res, err := http.Get(countryURL)
	if err != nil {
		fmt.Println("Get country map failed")
		return "undefined", err
	}
	defer res.Body.Close()
	countries, err := io.ReadAll(res.Body)

	var countryNames map[string][]string

	err = json.Unmarshal(countries, &countryNames)

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

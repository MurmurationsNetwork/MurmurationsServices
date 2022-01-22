package countries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FindAlpha2ByName(country interface{}) (countryCode string, err error) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(path)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)

	f, err := os.Open("countries.json")
	fmt.Println(err)
	if err != nil {
		return "undefined", err
	}

	file, err := ioutil.ReadAll(f)
	fmt.Println(err)
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

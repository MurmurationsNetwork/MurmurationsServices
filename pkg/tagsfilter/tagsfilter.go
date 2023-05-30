package tagsfilter

import (
	"encoding/json"
	"unicode/utf8"

	"golang.org/x/exp/utf8string"
)

type Data struct {
	Tags []interface{} `json:"tags"`
}

func Filter(
	arraySize int,
	stringLength int,
	profileStr string,
) ([]string, error) {
	var (
		size   int
		result []string
		data   Data
	)

	err := json.Unmarshal([]byte(profileStr), &data)
	if err != nil {
		return nil, err
	}

	for _, value := range data.Tags {
		if size >= arraySize {
			break
		}
		// filter non-string
		switch v := value.(type) {
		case string:
			// filter length (truncate) [#227]
			if utf8.RuneCountInString(v) > stringLength {
				s := utf8string.NewString(v)
				v = s.Slice(0, stringLength)
			}
			result = append(result, v)
		default:
			continue
		}

		size++
	}

	return result, nil
}

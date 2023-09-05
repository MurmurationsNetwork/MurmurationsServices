package schemavalidator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomValidator interface {
	Validate(value interface{}) *ValidationResult
}

type GeolocationValidator struct{}

func (v *GeolocationValidator) Validate(value interface{}) *ValidationResult {
	vr := NewValidationResult()
	if geoValue, ok := value.(map[string]interface{}); ok {
		if lat, exists := geoValue["lat"]; exists {
			if latValue, ok := lat.(json.Number); ok {
				f, err := latValue.Float64()
				if err != nil {
					vr.AppendError(
						"Invalid Latitude Type",
						"Latitude should be a number",
						[]string{"pointer", "/geolocation/lat"},
						http.StatusBadRequest,
					)
				} else if f < -90 || f > 90 {
					vr.AppendError(
						"Invalid Latitude",
						"Latitude should be between -90 and 90",
						[]string{"pointer", "/geolocation/lat"},
						http.StatusBadRequest,
					)
				}
			} else {
				vr.AppendError(
					"Invalid Latitude Type",
					"Latitude should be a json.Number",
					[]string{"pointer", "/geolocation/lat"},
					http.StatusBadRequest,
				)
			}
		}
		if lon, exists := geoValue["lon"]; exists {
			if lonValue, ok := lon.(json.Number); ok {
				f, err := lonValue.Float64()
				if err != nil {
					vr.AppendError(
						"Invalid Longitude Type",
						"Longitude should be a number",
						[]string{"pointer", "/geolocation/lon"},
						http.StatusBadRequest,
					)
				} else if f < -180 || f > 180 {
					vr.AppendError(
						"Invalid Longitude",
						"Longitude should be between -180 and 180",
						[]string{"pointer", "/geolocation/lon"},
						http.StatusBadRequest,
					)
				}
			} else {
				vr.AppendError(
					"Invalid Longitude Type",
					"Longitude should be a json.Number",
					[]string{"pointer", "/geolocation/lon"},
					http.StatusBadRequest,
				)
			}
		}
	} else {
		vr.AppendError(
			"Invalid Geolocation Type",
			"Geolocation should be an object",
			[]string{"pointer", "/geolocation"},
			http.StatusBadRequest,
		)
	}
	return vr
}

type StringValidator struct {
	MaxLength int
}

func (v *StringValidator) Validate(value interface{}) *ValidationResult {
	vr := NewValidationResult()
	if strValue, ok := value.(string); ok {
		if len(strValue) > v.MaxLength {
			vr.AppendError(
				fmt.Sprintf("Invalid Length, max length is %d", v.MaxLength),
				"String length exceeded",
				nil,
				http.StatusBadRequest,
			)
		}
	} else {
		vr.AppendError("Invalid Type", "Should be a string", nil, http.StatusBadRequest)
	}
	return vr
}

type TagsValidator struct{}

func (v *TagsValidator) Validate(value interface{}) *ValidationResult {
	vr := NewValidationResult()
	if tags, ok := value.([]interface{}); ok {
		if len(tags) > 100 {
			vr.AppendError(
				"Too Many Tags",
				"Maximum of 100 tags allowed",
				nil,
				http.StatusBadRequest,
			)
		}
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				if len(tagStr) > 100 {
					vr.AppendError(
						"Tag Too Long",
						"Each tag should be under 100 characters",
						nil,
						http.StatusBadRequest,
					)
				}
			} else {
				vr.AppendError(
					"Invalid Tag Type",
					"Tags should be strings",
					nil,
					http.StatusBadRequest,
				)
			}
		}
	} else {
		vr.AppendError(
			"Invalid Tags Type",
			"Tags should be an array of strings",
			nil,
			http.StatusBadRequest,
		)
	}
	return vr
}

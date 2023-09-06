package schemavalidator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// CustomValidator defines the interface for custom validation.
type CustomValidator interface {
	// Validate performs custom validation.
	Validate(value interface{}) *ValidationResult
}

// GeolocationValidator validates geolocation data.
type GeolocationValidator struct{}

// Validate checks the validity of geolocation data.
func (v *GeolocationValidator) Validate(value interface{}) *ValidationResult {
	vr := NewValidationResult()

	switch geoValue := value.(type) {
	case map[string]interface{}:
		validateLatLon(geoValue, vr)
	case string:
		coords := strings.Split(geoValue, ",")
		if len(coords) != 2 {
			vr.AppendError(
				"Invalid Geolocation Format",
				"Geolocation string should be in 'lat,lon' format",
				[]string{"pointer", "/geolocation"},
				http.StatusBadRequest,
			)
			return vr
		}
		lat, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			vr.AppendError(
				"Invalid Latitude Type",
				"Latitude should be a number",
				[]string{"pointer", "/geolocation/lat"},
				http.StatusBadRequest,
			)
		}
		lon, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			vr.AppendError(
				"Invalid Longitude Type",
				"Longitude should be a number",
				[]string{"pointer", "/geolocation/lon"},
				http.StatusBadRequest,
			)
		}
		validateLatLon(map[string]interface{}{"lat": lat, "lon": lon}, vr)
	default:
		vr.AppendError(
			"Invalid Geolocation Type",
			"Geolocation should be an object or a string",
			[]string{"pointer", "/geolocation"},
			http.StatusBadRequest,
		)
	}
	return vr
}

func validateLatLon(geoValue map[string]interface{}, vr *ValidationResult) {
	validateCoordinate := func(coord interface{}, name string, min, max float64) {
		var f float64
		var err error

		switch v := coord.(type) {
		case json.Number:
			f, err = v.Float64()
		case float64:
			f = v
		default:
			vr.AppendError(
				"Invalid "+name+" Type",
				name+" should be a number",
				[]string{"pointer", "/geolocation/" + name},
				http.StatusBadRequest,
			)
			return
		}

		if err != nil {
			vr.AppendError(
				"Invalid "+name+" Type",
				name+" should be a number",
				[]string{"pointer", "/geolocation/" + name},
				http.StatusBadRequest,
			)
			return
		}

		if f < min || f > max {
			vr.AppendError(
				"Invalid "+name,
				fmt.Sprintf("%s should be between %f and %f", name, min, max),
				[]string{"pointer", "/geolocation/" + name},
				http.StatusBadRequest,
			)
		}
	}

	// Validate latitude if exists.
	if lat, exists := geoValue["lat"]; exists {
		validateCoordinate(lat, "Latitude", -90, 90)
	}

	// Validate longitude if exists.
	if lon, exists := geoValue["lon"]; exists {
		validateCoordinate(lon, "Longitude", -180, 180)
	}

	// Check for extra keys
	for key := range geoValue {
		if key != "lat" && key != "lon" {
			vr.AppendError(
				"Extra Field in Geolocation",
				fmt.Sprintf(
					"Extra field '%s' found in Geolocation object",
					key,
				),
				[]string{"pointer", "/geolocation/" + key},
				http.StatusBadRequest,
			)
		}
	}
}

// StringValidator validates string data with a maximum length constraint.
type StringValidator struct {
	MaxLength int
}

// Validate checks if the given string exceeds the maximum length.
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

// TagsValidator validates an array of tags.
type TagsValidator struct{}

// Validate checks the validity of an array of tags.
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

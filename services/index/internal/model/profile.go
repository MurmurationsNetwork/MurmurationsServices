package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/constant"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/countries"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/tagsfilter"
	"github.com/MurmurationsNetwork/MurmurationsServices/services/index/config"
)

// AllowedFields specifies indexable keys.
//
// Elasticsearch stores the original JSON document, and to prevent the index
// from containing garbage data, we manually filter out unwanted fields.
var AllowedFields = map[string]bool{
	"country":        true,
	"geolocation":    true,
	"last_updated":   true,
	"linked_schemas": true,
	"locality":       true,
	"name":           true,
	"primary_url":    true,
	"profile_url":    true,
	"region":         true,
	"status":         true,
	"tags":           true,
	"expires":        true,
}

// Profile represents the profile data for a node.
type Profile struct {
	// Original profile string.
	str string

	// JSON representation of the profile.
	json map[string]interface{}
}

// NewProfile initializes a new Profile based on a profile string.
func NewProfile(profileStr string) *Profile {
	profileJSON := jsonutil.ToJSON(profileStr)
	return &Profile{
		str:  profileStr,
		json: profileJSON,
	}
}

// GetJSON returns the filtered JSON representation of the profile.
func (p *Profile) GetJSON() map[string]interface{} {
	filteredJSON := make(map[string]interface{})

	for key, value := range p.json {
		if _, ok := AllowedFields[key]; ok {
			filteredJSON[key] = value
		}
	}

	return filteredJSON
}

// Update processes and updates the profile data.
func (p *Profile) Update(
	profileURL string,
	lastUpdated *int64,
) error {
	p.json["profile_url"] = profileURL
	p.json["last_updated"] = lastUpdated

	if err := p.convertGeolocation(); err != nil {
		return err
	}

	if err := p.normalizeCountryCode(); err != nil {
		return err
	}

	if err := p.filterTags(); err != nil {
		return err
	}

	p.setDefaultStatus()

	return nil
}

// convertGeolocation standardizes the profile's geolocation format.
//
// It parses a "latitude,longitude" string or combines separate "latitude" and
// "longitude" fields into a map.
// It won't create a geolocation object if no valid data is found.
func (p *Profile) convertGeolocation() error {
	var lat, lon float64
	var err error

	if geoStr, ok := p.json["geolocation"].(string); ok {
		g := strings.Split(geoStr, ",")
		if len(g) != 2 {
			return fmt.Errorf("invalid geolocation format")
		}

		lat, err = strconv.ParseFloat(g[0], 64)
		if err != nil {
			return fmt.Errorf("invalid latitude: %v", err)
		}

		lon, err = strconv.ParseFloat(g[1], 64)
		if err != nil {
			return fmt.Errorf("invalid longitude: %v", err)
		}

		p.json["geolocation"] = map[string]interface{}{"lat": lat, "lon": lon}
	} else if existingGeo, ok := p.json["geolocation"].(map[string]interface{}); ok {
		p.json["geolocation"] = existingGeo
	} else {
		geoLocation := make(map[string]interface{})
		if existingLat, ok := p.json["latitude"].(float64); ok {
			geoLocation["lat"] = existingLat
		}
		if existingLon, ok := p.json["longitude"].(float64); ok {
			geoLocation["lon"] = existingLon
		}
		if len(geoLocation) > 0 {
			p.json["geolocation"] = geoLocation
		}
	}

	return nil
}

// normalizeCountryCode takes in profile JSON and normalizes the country information present.
func (p *Profile) normalizeCountryCode() error {
	if p.json["country_iso_3166"] != nil {
		p.json["country"] = p.json["country_iso_3166"]
		delete(p.json, "country_iso_3166")
		return nil
	}

	if p.json["country"] != nil || p.json["country_name"] == nil {
		// If 'country' is already defined or 'country_name' does not exist, no need to proceed.
		return nil
	}

	countryCode, err := countries.FindAlpha2ByName(
		config.Values.Library.InternalURL+"/v2/countries",
		p.json["country_name"],
	)
	if err != nil {
		if errors.Is(err, countries.ErrCountryCodeNotFound) {
			logger.Info("Country code not found",
				zap.String("country", p.json["country_name"].(string)),
				zap.String("profile_url", p.json["profile_url"].(string)),
			)
			return nil
		}
		return err
	}

	p.json["country"] = countryCode
	logger.Info("Country code matched",
		zap.String("country", p.json["country_name"].(string)),
		zap.String("code", countryCode),
		zap.String("profile_url", p.json["profile_url"].(string)),
	)

	return nil
}

// filterTags filters the tags based on the provided configuration values and
// attaches them to the profile JSON.
func (p *Profile) filterTags() error {
	arraySize, _ := strconv.Atoi(config.Values.Server.TagsArraySize)
	stringLength, _ := strconv.Atoi(config.Values.Server.TagsStringLength)
	tags, err := tagsfilter.Filter(arraySize, stringLength, p.str)
	if err != nil {
		return err
	}
	if len(tags) != 0 {
		p.json["tags"] = tags
	}
	return nil
}

// setDefaultStatus sets the default status of the profile.
func (p *Profile) setDefaultStatus() {
	p.json["status"] = constant.NodeStatus.Posted
}

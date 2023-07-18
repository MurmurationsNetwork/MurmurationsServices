package model

import (
	"errors"
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

// GetJSON returns the JSON representation of the profile.
func (p *Profile) GetJSON() map[string]interface{} {
	return p.json
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
	p.repackageGeolocation()

	if err := p.normalizeCountryCode(); err != nil {
		return err
	}

	if err := p.filterTags(); err != nil {
		return err
	}

	p.setDefaultStatus()

	return nil
}

// convertGeolocation parses a string-formatted geolocation into latitude and
// longitude values. The geolocation should be in the format "latitude,longitude".
func (p *Profile) convertGeolocation() error {
	if geo, ok := p.json["geolocation"].(string); ok {
		g := strings.Split(geo, ",")
		var err error
		p.json["latitude"], err = strconv.ParseFloat(g[0], 64)
		if err != nil {
			return err
		}
		p.json["longitude"], err = strconv.ParseFloat(g[1], 64)
		if err != nil {
			return err
		}
	}
	return nil
}

// repackageGeolocation reformats the geolocation field as a map with keys "lat" and "lon".
func (p *Profile) repackageGeolocation() {
	if p.json["latitude"] != nil || p.json["longitude"] != nil {
		geoLocation := make(map[string]interface{})

		if p.json["latitude"] != nil {
			geoLocation["lat"] = p.json["latitude"]
		} else {
			geoLocation["lat"] = 0
		}

		if p.json["longitude"] != nil {
			geoLocation["lon"] = p.json["longitude"]
		} else {
			geoLocation["lon"] = 0
		}

		p.json["geolocation"] = geoLocation
	}
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

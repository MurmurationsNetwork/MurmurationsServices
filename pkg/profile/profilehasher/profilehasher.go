package profilehasher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/core"
	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonutil"
)

// defaultFields represents the default fields to be considered for hashing.
var defaultFields = map[string]bool{
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
	"latitude":       true,
	"longitude":      true,
}

// ProfileHash contains data required to hash a profile.
type ProfileHash struct {
	profileURL       string
	libraryURL       string
	profile          map[string]interface{}
	fieldsForHashing map[string]bool
	profileStr       string
}

// New initializes a new ProfileHash instance.
func New(profileURL, libraryURL string) *ProfileHash {
	fieldsForHashing := make(map[string]bool)
	for key := range defaultFields {
		fieldsForHashing[key] = true
	}

	return &ProfileHash{
		profileURL:       profileURL,
		libraryURL:       libraryURL,
		fieldsForHashing: fieldsForHashing,
	}
}

// NewFromString initializes a new ProfileHash instance with a JSON string.
func NewFromString(profileStr, libraryURL string) *ProfileHash {
	fieldsForHashing := make(map[string]bool)
	for key := range defaultFields {
		fieldsForHashing[key] = true
	}

	return &ProfileHash{
		libraryURL:       libraryURL,
		profileStr:       profileStr,
		fieldsForHashing: fieldsForHashing,
	}
}

// Hash computes the hash for the profile.
func (p *ProfileHash) Hash() (string, error) {
	var err error
	if p.profileStr != "" {
		err = json.Unmarshal([]byte(p.profileStr), &p.profile)
	} else {
		var profileData []byte
		profileData, err = p.fetchData(p.profileURL)
		if err == nil {
			err = json.Unmarshal(profileData, &p.profile)
		}
	}

	if err != nil {
		return "", core.ProfileFetchError{Reason: err.Error()}
	}

	err = p.populateFieldsForHashing()
	if err != nil {
		return "", core.ProfileFetchError{Reason: err.Error()}
	}

	filteredProfile := p.filterProfileFields()

	hashedValue, err := jsonutil.Hash(filteredProfile)
	if err != nil {
		return "", fmt.Errorf(
			"error computing hash for filtered profile: %v",
			err,
		)
	}

	return hashedValue, nil
}

func (p *ProfileHash) fetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(
			"error while sending request to %s: %v",
			url,
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request to %s failed with status code %d",
			url,
			resp.StatusCode,
		)
	}

	return io.ReadAll(resp.Body)
}

func (p *ProfileHash) populateFieldsForHashing() error {
	linkedSchemas, ok := p.profile["linked_schemas"].([]interface{})

	if ok {
		for _, schemaIdentifier := range linkedSchemas {
			schemaURL := p.constructSchemaURL(schemaIdentifier.(string))

			schemaData, err := p.fetchData(schemaURL)
			if err != nil {
				return err
			}

			var currentSchema map[string]interface{}
			err = json.Unmarshal(schemaData, &currentSchema)
			if err != nil {
				return err
			}

			properties, ok := currentSchema["properties"].(map[string]interface{})
			if ok {
				for key := range properties {
					p.fieldsForHashing[key] = true
				}
			}
		}
	}

	return nil
}

func (p *ProfileHash) filterProfileFields() map[string]interface{} {
	filteredProfile := make(map[string]interface{})
	for key, value := range p.profile {
		if p.fieldsForHashing[key] {
			filteredProfile[key] = value
		}
	}
	return filteredProfile
}

func (p *ProfileHash) constructSchemaURL(linkedSchema string) string {
	return fmt.Sprintf("%s/v2/schemas/%s", p.libraryURL, linkedSchema)
}

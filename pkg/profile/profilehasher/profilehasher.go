package profilehasher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

// Hash computes the hash for the profile.
func (p *ProfileHash) Hash() (string, error) {
	profileData, err := p.fetchData(p.profileURL)
	if err != nil {
		return "", err
	}

	err = p.parseProfile(profileData)
	if err != nil {
		return "", err
	}

	err = p.populateFieldsForHashing()
	if err != nil {
		return "", err
	}

	filteredProfile := p.filterProfileFields()

	return jsonutil.Hash(filteredProfile)
}

func (p *ProfileHash) fetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (p *ProfileHash) parseProfile(data []byte) error {
	return json.Unmarshal(data, &p.profile)
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

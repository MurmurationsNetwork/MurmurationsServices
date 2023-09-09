package profilevalidator_test

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
)

var StrSchema = `{}`

func TestValidate_CustomValidator(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		valid   bool
		reason  string
	}{
		{
			name:    "Valid Geolocation",
			profile: `{"geolocation": {"lat": 50.123, "lon": -75.123}}`,
			valid:   true,
		},
		{
			name:    "Invalid Latitude",
			profile: `{"geolocation": {"lat": 100.123, "lon": -75.123}}`,
			valid:   false,
			reason:  "Latitude exceeds 90",
		},
		{
			name:    "Invalid Longitude",
			profile: `{"geolocation": {"lat": 50.123, "lon": -200.123}}`,
			valid:   false,
			reason:  "Longitude exceeds 180",
		},
		{
			name:    "Extra Fields in Geolocation",
			profile: `{"geolocation": {"lat": 50.123, "lon": -75.123, "extraField": "invalid"}}`,
			valid:   false,
			reason:  "Extra field 'extraField' found in Geolocation object",
		},
		{
			name:    "Valid Geolocation as String",
			profile: `{"geolocation": "50.123,-75.123"}`,
			valid:   true,
		},
		{
			name:    "Invalid Geolocation as String",
			profile: `{"geolocation": "100.123,-75.123"}`,
			valid:   false,
			reason:  "Latitude exceeds 90",
		},
		{
			name:    "Invalid Geolocation Format",
			profile: `{"geolocation": "50.123"}`,
			valid:   false,
			reason:  "Geolocation string should be in 'lat,lon' format",
		},
		{
			name:    "Valid Name",
			profile: `{"name": "John Doe"}`,
			valid:   true,
		},
		{
			name:    "Invalid Name",
			profile: fmt.Sprintf(`{"name": "%s"}`, randomString(201)),
			valid:   false,
			reason:  "Name length exceeds 200 characters",
		},
		{
			name:    "Valid Locality",
			profile: `{"locality": "New York"}`,
			valid:   true,
		},
		{
			name:    "Invalid Locality",
			profile: fmt.Sprintf(`{"locality": "%s"}`, randomString(101)),
			valid:   false,
			reason:  "Locality length exceeds 100 characters",
		},
		{
			name:    "Valid Country Name",
			profile: `{"country_name": "United States"}`,
			valid:   true,
		},
		{
			name:    "Invalid Country Name",
			profile: fmt.Sprintf(`{"country_name": "%s"}`, randomString(101)),
			valid:   false,
			reason:  "Country name exceeds 100 characters",
		},
		{
			name:    "Valid Country",
			profile: `{"country_iso_3166": "USA"}`,
			valid:   true,
		},
		{
			name:    "Invalid Country",
			profile: `{"country_iso_3166": "United States"}`,
			valid:   false,
			reason:  "Country ISO code exceeds 3 characters",
		},
		{
			name:    "Valid Tags",
			profile: `{"tags": ["tag1", "tag2"]}`,
			valid:   true,
		},
		{
			name:    "Invalid Tags",
			profile: fmt.Sprintf(`{"tags": ["%s", "tag2"]}`, randomString(101)),
			valid:   false,
			reason:  "Tag length exceeds 100 characters",
		},
		{
			name: "Invalid Number of Tags",
			profile: fmt.Sprintf(
				`{"tags": [%s]}`,
				strings.Join(generateTags(101), ","),
			),
			valid:  false,
			reason: "Number of tags exceeds 100",
		},
		{
			name:    "Valid Primary URL",
			profile: `{"primary_url": "https://example.com"}`,
			valid:   true,
		},
		{
			name:    "Invalid Primary URL",
			profile: fmt.Sprintf(`{"primary_url": "%s"}`, randomURL(2001)),
			valid:   false,
			reason:  "Primary URL length exceeds 2000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := profilevalidator.NewBuilder().
				WithStrSchemas([]string{StrSchema}).
				WithStrProfile(tt.profile).
				WithCustomValidation().
				Build()

			require.NoError(t, err)
			result := validator.Validate()
			require.Equal(t, tt.valid, result.Valid)
		})
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		b.WriteByte(charset[randomIndex])
	}
	return b.String()
}

func randomURL(length int) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "example.com",
		// 18 is the length of "https://example.com".
		Path: randomString(length - 18),
	}
	return u.String()
}

func generateTags(count int) []string {
	tags := make([]string, count)
	for i := 0; i < count; i++ {
		tags[i] = fmt.Sprintf("\"tag%d\"", i+1)
	}
	return tags
}

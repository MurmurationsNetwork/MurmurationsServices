package profilevalidator_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
)

func TestStringValidator_Validate(t *testing.T) {
	tests := []struct {
		name            string
		input           interface{}
		maxLength       int
		path            string
		expectedValid   bool
		expectedErrMsgs []string
		expectedDetails []string
		expectedSources [][]string
		expectedStatus  []int
	}{
		{
			name:            "Valid string within max length",
			input:           "valid",
			maxLength:       10,
			path:            "username",
			expectedValid:   true,
			expectedErrMsgs: nil,
			expectedDetails: nil,
			expectedSources: nil,
			expectedStatus:  nil,
		},
		{
			name:            "String exceeds max length",
			input:           "USA",
			maxLength:       2,
			path:            "country_iso_3166",
			expectedValid:   false,
			expectedErrMsgs: []string{"String Length Exceeded"},
			expectedDetails: []string{"Invalid Length, max length is 2"},
			expectedSources: [][]string{{"pointer", "/country_iso_3166"}},
			expectedStatus:  []int{http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := profilevalidator.StringValidator{
				MaxLength: tt.maxLength,
				Path:      tt.path,
			}

			vr := validator.Validate(tt.input)

			// Validate the validity of the result
			require.Equal(t, tt.expectedValid, vr.Valid)

			// Validate the number of errors
			require.Equal(t, len(tt.expectedErrMsgs), len(vr.ErrorMessages))

			// If errors exist, validate their content
			if len(vr.ErrorMessages) > 0 {
				require.Equal(t, tt.expectedErrMsgs, vr.ErrorMessages)
				require.Equal(t, tt.expectedDetails, vr.Details)
				require.Equal(t, tt.expectedSources, vr.Sources)
				require.Equal(t, tt.expectedStatus, vr.ErrorStatus)
			}
		})
	}
}

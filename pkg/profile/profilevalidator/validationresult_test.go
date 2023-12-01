package profilevalidator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/profile/profilevalidator"
)

func TestNewValidationResult(t *testing.T) {
	vr := profilevalidator.NewValidationResult()
	require.True(t, vr.Valid)
	require.Empty(t, vr.ErrorMessages)
	require.Empty(t, vr.Details)
	require.Empty(t, vr.Sources)
	require.Empty(t, vr.ErrorStatus)
}

func TestAppendError(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		detail       string
		source       []string
		status       int
		expVR        *profilevalidator.ValidationResult
	}{
		{
			name:         "Single error",
			errorMessage: "Error 1",
			detail:       "Detail 1",
			source:       []string{"Source1"},
			status:       400,
			expVR: &profilevalidator.ValidationResult{
				Valid:         false,
				ErrorMessages: []string{"Error 1"},
				Details:       []string{"Detail 1"},
				Sources:       [][]string{{"Source1"}},
				ErrorStatus:   []int{400},
			},
		},
		{
			name:         "Error with multiple sources",
			errorMessage: "Error 2",
			detail:       "Detail 2",
			source:       []string{"Source2a", "Source2b"},
			status:       500,
			expVR: &profilevalidator.ValidationResult{
				Valid:         false,
				ErrorMessages: []string{"Error 2"},
				Details:       []string{"Detail 2"},
				Sources:       [][]string{{"Source2a", "Source2b"}},
				ErrorStatus:   []int{500},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := profilevalidator.NewValidationResult()
			vr.AppendError(tt.errorMessage, tt.detail, tt.source, tt.status)
			require.False(t, vr.Valid)
			require.Equal(t, tt.expVR, vr)
		})
	}
}

func TestAppendErrors(t *testing.T) {
	tests := []struct {
		name          string
		errorMessages []string
		details       []string
		sources       [][]string
		statuses      []int
		expVR         *profilevalidator.ValidationResult
	}{
		{
			name:          "Multiple errors",
			errorMessages: []string{"Error 1", "Error 2"},
			details:       []string{"Detail 1", "Detail 2"},
			sources:       [][]string{{"Source1"}, {"Source2"}},
			statuses:      []int{400, 404},
			expVR: &profilevalidator.ValidationResult{
				Valid:         false,
				ErrorMessages: []string{"Error 1", "Error 2"},
				Details:       []string{"Detail 1", "Detail 2"},
				Sources:       [][]string{{"Source1"}, {"Source2"}},
				ErrorStatus:   []int{400, 404},
			},
		},
		{
			name:          "No errors",
			errorMessages: []string{},
			details:       []string{},
			sources:       [][]string{},
			statuses:      []int{},
			expVR: &profilevalidator.ValidationResult{
				Valid:         true,
				ErrorMessages: []string{},
				Details:       []string{},
				Sources:       [][]string{},
				ErrorStatus:   []int{},
			},
		},
		{
			name:          "Mismatched array lengths",
			errorMessages: []string{"Error 1", "Error 2"},
			details:       []string{"Detail 1"},
			sources:       [][]string{{"Source1"}, {"Source2"}},
			statuses:      []int{400},
			expVR: &profilevalidator.ValidationResult{
				Valid:         false,
				ErrorMessages: []string{"Error 1", "Error 2"},
				Details:       []string{"Detail 1", ""},
				Sources:       [][]string{{"Source1"}, {"Source2"}},
				ErrorStatus:   []int{400, 0},
			},
		},
		{
			name:          "Empty error messages",
			errorMessages: []string{},
			details:       []string{"Detail 1", "Detail 2"},
			sources:       [][]string{{"Source1"}, {"Source2"}},
			statuses:      []int{400, 404},
			expVR: &profilevalidator.ValidationResult{
				Valid:         false,
				ErrorMessages: []string{"", ""},
				Details:       []string{"Detail 1", "Detail 2"},
				Sources:       [][]string{{"Source1"}, {"Source2"}},
				ErrorStatus:   []int{400, 404},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := profilevalidator.NewValidationResult()
			vr.AppendErrors(
				tt.errorMessages,
				tt.details,
				tt.sources,
				tt.statuses,
			)
			require.Equal(t, tt.expVR, vr)
		})
	}
}

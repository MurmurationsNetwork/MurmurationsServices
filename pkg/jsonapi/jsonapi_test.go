package jsonapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
)

func TestNewError(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		titles         []string
		details        []string
		sources        [][]string
		status         []int
		expectedLen    int
		expectedErrors []jsonapi.Error
	}{
		{
			name:    "Standard case",
			titles:  []string{"Title1", "Title2"},
			details: []string{"Detail1", "Detail2"},
			sources: [][]string{{"Key1", "Value1"}, {"Key2", "Value2"}},
			status:  []int{200, 404},
			expectedErrors: []jsonapi.Error{
				{
					Status: 200,
					Title:  "Title1",
					Detail: "Detail1",
					Source: map[string]string{"Key1": "Value1"},
				},
				{
					Status: 404,
					Title:  "Title2",
					Detail: "Detail2",
					Source: map[string]string{"Key2": "Value2"},
				},
			},
		},
		{
			name:    "Single error with incomplete source",
			titles:  []string{"ErrorTitle"},
			details: []string{"ErrorDetail"},
			sources: [][]string{{"SourceKey"}},
			status:  []int{500},
			expectedErrors: []jsonapi.Error{
				{
					Status: 500,
					Title:  "ErrorTitle",
					Detail: "ErrorDetail",
					Source: map[string]string{"SourceKey": ""},
				},
			},
		},
		{
			name:    "No details and sources",
			titles:  []string{"ErrorTitle1", "ErrorTitle2"},
			details: []string{},
			sources: [][]string{},
			status:  []int{400, 403},
			expectedErrors: []jsonapi.Error{
				{Status: 400, Title: "ErrorTitle1"},
				{Status: 403, Title: "ErrorTitle2"},
			},
		},
		{
			name:    "More statuses than titles",
			titles:  []string{"ErrorTitle"},
			details: []string{"ErrorDetail"},
			sources: [][]string{{"Key", "Value"}},
			status:  []int{401, 500},
			expectedErrors: []jsonapi.Error{
				{
					Status: 401,
					Title:  "ErrorTitle",
					Detail: "ErrorDetail",
					Source: map[string]string{"Key": "Value"},
				},
			},
		},
		{
			name:           "Empty titles with details and sources",
			titles:         []string{},
			details:        []string{"ErrorDetail1", "ErrorDetail2"},
			sources:        [][]string{{"Key1", "Value1"}, {"Key2", "Value2"}},
			status:         []int{402, 404},
			expectedErrors: nil,
		},
		{
			name:           "Nil slices",
			titles:         nil,
			details:        nil,
			sources:        nil,
			status:         nil,
			expectedErrors: nil, // No errors expected as input slices are nil
		},
		{
			name:    "Titles with empty strings",
			titles:  []string{"", ""},
			details: []string{"Detail1", "Detail2"},
			sources: [][]string{{"Key1", "Value1"}, {"Key2", ""}},
			status:  []int{405, 406},
			expectedErrors: []jsonapi.Error{
				{
					Status: 405,
					Title:  "",
					Detail: "Detail1",
					Source: map[string]string{"Key1": "Value1"},
				},
				{
					Status: 406,
					Title:  "",
					Detail: "Detail2",
					Source: map[string]string{"Key2": ""},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errors := jsonapi.NewError(
				tc.titles,
				tc.details,
				tc.sources,
				tc.status,
			)
			require.Equal(
				t,
				tc.expectedErrors,
				errors,
				"Errors do not match the expected output",
			)
		})
	}
}

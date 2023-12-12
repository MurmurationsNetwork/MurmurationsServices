package jsonapi_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/jsonapi"
)

// Mock request to create a gin context.
func mockRequest(method, path string) *gin.Context {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}

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

func TestNewLinks(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int64
		totalPage   int64
		requestPath string
		expected    *jsonapi.Link
	}{
		{
			name:        "TestFirstPage",
			currentPage: 1,
			totalPage:   5,
			requestPath: "/test?page=1",
			expected: &jsonapi.Link{
				First: "",
				Prev:  "",
				Self:  "http://example.com/test?page=1",
				Next:  "http://example.com/test?page=2",
				Last:  "http://example.com/test?page=5",
			},
		},
		{
			name:        "TestMiddlePage",
			currentPage: 3,
			totalPage:   5,
			requestPath: "/test?page=3",
			expected: &jsonapi.Link{
				First: "http://example.com/test?page=1",
				Prev:  "http://example.com/test?page=2",
				Self:  "http://example.com/test?page=3",
				Next:  "http://example.com/test?page=4",
				Last:  "http://example.com/test?page=5",
			},
		},
		{
			name:        "TestLastPage",
			currentPage: 5,
			totalPage:   5,
			requestPath: "/test?page=5",
			expected: &jsonapi.Link{
				First: "http://example.com/test?page=1",
				Prev:  "http://example.com/test?page=4",
				Self:  "http://example.com/test?page=5",
				Next:  "",
				Last:  "",
			},
		},
		{
			name:        "TestSinglePage",
			currentPage: 1,
			totalPage:   1,
			requestPath: "/test?page=1",
			expected: &jsonapi.Link{
				First: "",
				Prev:  "",
				Self:  "http://example.com/test?page=1",
				Next:  "",
				Last:  "",
			},
		},
		{
			name:        "TestURLWithOtherParams",
			currentPage: 2,
			totalPage:   4,
			requestPath: "/test?param=value&page=2",
			expected: &jsonapi.Link{
				First: "http://example.com/test?page=1&param=value",
				Prev:  "http://example.com/test?page=1&param=value",
				Self:  "http://example.com/test?page=2&param=value",
				Next:  "http://example.com/test?page=3&param=value",
				Last:  "http://example.com/test?page=4&param=value",
			},
		},
		{
			name:        "TestCurrentPageGreaterThanTotal",
			currentPage: 6,
			totalPage:   5,
			requestPath: "/test?page=6",
			expected: &jsonapi.Link{
				First: "http://example.com/test?page=1",
				Prev:  "http://example.com/test?page=4",
				Self:  "http://example.com/test?page=5",
				Next:  "",
				Last:  "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := mockRequest("GET", tc.requestPath)

			link := jsonapi.NewLinks(c, tc.currentPage, tc.totalPage)

			require.Equal(t, tc.expected, link)
		})
	}
}

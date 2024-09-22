package jsonapi

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MurmurationsNetwork/MurmurationsServices/pkg/logger"
)

type JSONAPI struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
	Links  *Link       `json:"links,omitempty"`
	Meta   *Meta       `json:"meta,omitempty"`
}

type Error struct {
	Status int               `json:"status,omitempty"`
	Source map[string]string `json:"source,omitempty"`
	Title  string            `json:"title,omitempty"`
	Detail string            `json:"detail,omitempty"`
}

type Link struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Self  string `json:"self,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

type Meta struct {
	Message         string        `json:"message,omitempty"`
	NodeID          string        `json:"node_id,omitempty"`
	ProfileURL      string        `json:"profile_url,omitempty"`
	NumberOfResults int64         `json:"number_of_results,omitempty"`
	TotalPages      int64         `json:"total_pages,omitempty"`
	Sort            []interface{} `json:"sort,omitempty"`
	BatchID         string        `json:"batch_id,omitempty"`
}

// JSON API Response Combination

func Response(
	data interface{},
	errors []Error,
	link *Link,
	meta *Meta,
) JSONAPI {
	return JSONAPI{
		Data:   data,
		Errors: errors,
		Links:  link,
		Meta:   meta,
	}
}

// JSON API Internal Data

// NewErrors constructs a slice of Error objects from the provided slices.
// Each Error object is populated with the corresponding elements from the
// input slices.
//
// Parameters:
//   - titles: a slice of error titles (required).
//   - details: a slice of error details (optional).
//   - sources: a slice of slices, where each inner slice contains key-value
//     pairs for the Source map in an Error object (optional).
//   - statuses: a slice of status codes (optional).
//
// The function handles cases where the optional slices have fewer elements
// than the titles slice. If an optional slice is shorter, the missing values
// are omitted for those Error objects.
//
// Example usage:
//
//	titles := []string{"Error 1", "Error 2"}
//	details := []string{"Detail 1", "Detail 2"}
//	sources := [][]string{
//	  {"file", "main.go", "line", "42"},
//	  {"file", "utils.go"},
//	}
//	statuses := []int{400, 404}
//
//	errors := NewErrors(titles, details, sources, statuses)
//
//	// The errors slice will contain:
//	errors[0] = Error{
//	  Title:  "Error 1",
//	  Detail: "Detail 1",
//	  Status: 400,
//	  Source: map[string]string{
//	    "file": "main.go",
//	    "line": "42",
//	  },
//	}
//	errors[1] = Error{
//	  Title:  "Error 2",
//	  Detail: "Detail 2",
//	  Status: 404,
//	  Source: map[string]string{
//	    "file": "utils.go",
//	  },
//	}
func NewError(
	titles []string,
	details []string,
	sources [][]string,
	statuses []int,
) []Error {
	errors := make([]Error, len(titles))

	for i, title := range titles {
		// Initialize the Error object with the title
		err := Error{
			Title: title,
		}

		// Set Status if available
		if i < len(statuses) {
			err.Status = statuses[i]
		}

		// Set Detail if available
		if i < len(details) {
			err.Detail = details[i]
		}

		// Set Source if available
		if i < len(sources) && len(sources[i]) > 0 {
			err.Source = make(map[string]string)
			sourcePairs := sources[i]

			// Iterate over the sourcePairs slice in steps of 2 to get key-value pairs
			for j := 0; j+1 < len(sourcePairs); j += 2 {
				key := sourcePairs[j]
				value := sourcePairs[j+1]
				err.Source[key] = value
			}

			// If there's an odd number of elements, assign an empty string to the last key
			if len(sourcePairs)%2 != 0 {
				key := sourcePairs[len(sourcePairs)-1]
				err.Source[key] = ""
			}
		}

		// Assign the Error object to the errors slice
		errors[i] = err
	}

	return errors
}

// NewLinks creates pagination links (first, previous, current, next, last)
// based on the current page and total pages.
// It adjusts the links according to the page position in the context of the
// total pagination.
func NewLinks(c *gin.Context, currentPage int64, totalPage int64) *Link {
	scheme := getURLScheme(c)
	base := getBaseURL(c, scheme)
	u, err := removePageParam(c.Request.RequestURI)
	if err != nil {
		logger.Error("Error removing page parameter", err)
		// Generate a special error link.
		errorLink := base + "?error=link-generation-failed"
		return &Link{
			First: errorLink,
			Prev:  errorLink,
			Self:  errorLink,
			Next:  errorLink,
			Last:  errorLink,
		}
	}

	// Ensure currentPage is within the valid range。
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPage {
		currentPage = totalPage
	}

	first, prev, self, next, last := createPaginationLinks(
		base,
		u,
		currentPage,
		totalPage,
	)

	return &Link{
		First: first,
		Prev:  prev,
		Self:  self,
		Next:  next,
		Last:  last,
	}
}

func getURLScheme(c *gin.Context) string {
	// First, check the X-Forwarded-Proto header.
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		return proto
	}

	// Fallback to checking if TLS is not nil.
	if c.Request.TLS != nil {
		return "https"
	}

	// Default to http if none of the above conditions are met.
	return "http"
}

func getBaseURL(c *gin.Context, scheme string) string {
	return scheme + "://" + c.Request.Host
}

func removePageParam(requestURI string) (*url.URL, error) {
	// Parse the original request URL.
	u, err := url.Parse(requestURI)
	if err != nil {
		return nil, err
	}

	// Get query values.
	queryValues := u.Query()

	// Remove the 'page' parameter.
	queryValues.Del("page")

	// Rebuild the query string without the 'page' parameter.
	u.RawQuery = queryValues.Encode()

	// Return the modified url.
	return u, nil
}

func createPaginationLinks(
	base string,
	u *url.URL,
	currentPage, totalPage int64,
) (first, prev, self, next, last string) {
	if currentPage > 1 {
		first = buildPageURL(base, u, 1)
		prev = buildPageURL(base, u, currentPage-1)
	}
	self = buildPageURL(base, u, currentPage)
	if currentPage < totalPage {
		next = buildPageURL(base, u, currentPage+1)
		last = buildPageURL(base, u, totalPage)
	}
	return
}

func buildPageURL(base string, u *url.URL, page int64) string {
	queryValues := u.Query()
	queryValues.Set("page", strconv.FormatInt(page, 10))
	return base + u.Path + "?" + queryValues.Encode()
}

func NewMeta(message string, nodeID string, profileURL string) *Meta {
	return &Meta{
		Message:    message,
		NodeID:     nodeID,
		ProfileURL: profileURL,
	}
}

func NewSearchMeta(
	message string,
	numberOfResults int64,
	totalPages int64,
) *Meta {
	return &Meta{
		Message:         message,
		NumberOfResults: numberOfResults,
		TotalPages:      totalPages,
	}
}

func NewBlockSearchMeta(sort []interface{}) *Meta {
	return &Meta{
		Sort: sort,
	}
}

func NewBatchMeta(message string, batchID string) *Meta {
	return &Meta{
		Message: message,
		BatchID: batchID,
	}
}

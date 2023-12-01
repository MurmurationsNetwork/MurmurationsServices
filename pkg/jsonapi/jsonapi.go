package jsonapi

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

func NewError(
	titles []string,
	details []string,
	sources [][]string,
	status []int,
) []Error {
	var errors []Error

	for i := 0; i < len(titles); i++ {
		newError := Error{
			Status: status[i],
			Title:  titles[i],
		}

		if i < len(details) {
			newError.Detail = details[i]
		}

		if i < len(sources) && len(sources[i]) > 0 {
			newError.Source = make(map[string]string)
			newError.Source[sources[i][0]] = ""
			if len(sources[i]) > 1 {
				newError.Source[sources[i][0]] = sources[i][1]
			}
		}

		errors = append(errors, newError)
	}

	return errors
}

func NewLinks(c *gin.Context, currentPage int64, totalPage int64) *Link {
	// check the request is http or https
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	url := scheme + "://" + c.Request.Host

	query := c.Request.RequestURI
	if currentPage != 0 {
		if strings.Contains(
			query,
			"?page="+strconv.Itoa(int(currentPage))+"&",
		) {
			// if page is the first query and query also has the other parameters
			query = strings.Replace(
				query,
				"page="+strconv.Itoa(int(currentPage))+"&",
				"",
				-1,
			)
		} else if strings.Contains(query, "?page="+strconv.Itoa(int(currentPage))) {
			// if query only has page parameter
			query = strings.Replace(query, "?page="+strconv.Itoa(int(currentPage)), "", -1)
		} else {
			query = strings.Replace(query, "&page="+strconv.Itoa(int(currentPage)), "", -1)
		}
	}

	// if query don't have a question mark, means it didn't have query, we need to add "?" to add the page query.
	// otherwise, it has query, we need to add "&" for adding page later.
	if !strings.Contains(query, "?") {
		query += "?"
	} else {
		query += "&"
	}

	// if page is not set, set to default 1
	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPage {
		currentPage = totalPage
	}

	// define the page
	first := url + query + "page=1"
	prev := url + query + "page=" + strconv.Itoa(int(currentPage-1))
	self := url + query + "page=" + strconv.Itoa(int(currentPage))
	next := url + query + "page=" + strconv.Itoa(int(currentPage+1))
	last := url + query + "page=" + strconv.Itoa(int(totalPage))

	if currentPage == 1 {
		first = ""
		prev = ""
	}
	if currentPage == totalPage {
		next = ""
		last = ""
	}
	if currentPage == 2 {
		first = ""
	}
	if currentPage+1 == totalPage {
		last = ""
	}

	return &Link{
		First: first,
		Prev:  prev,
		Self:  self,
		Next:  next,
		Last:  last,
	}
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

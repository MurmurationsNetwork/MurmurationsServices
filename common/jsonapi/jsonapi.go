package jsonapi

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type JsonApi struct {
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
	NodeId          string        `json:"node_id,omitempty"`
	ProfileUrl      string        `json:"profile_url,omitempty"`
	NumberOfResults int64         `json:"number_of_results,omitempty"`
	TotalPages      int64         `json:"total_pages,omitempty"`
	Sort            []interface{} `json:"sort,omitempty"`
}

// JSON API Response Combination

func Response(data interface{}, errors []Error, link *Link, meta *Meta) JsonApi {
	return JsonApi{
		Data:   data,
		Errors: errors,
		Links:  link,
		Meta:   meta,
	}
}

// JSON API Internal Data

func NewError(titles []string, details []string, sources [][]string, status []int) []Error {
	var errors []Error
	for i := 0; i < len(titles); i++ {
		if len(details) != 0 && len(sources) != 0 {
			errors = append(errors, Error{
				Status: status[i],
				Title:  titles[i],
				Detail: details[i],
				Source: map[string]string{
					sources[i][0]: sources[i][1],
				},
			})
		} else if len(details) != 0 {
			errors = append(errors, Error{
				Status: status[i],
				Title:  titles[i],
				Detail: details[i],
			})
		} else {
			errors = append(errors, Error{
				Status: status[i],
				Title:  titles[i],
			})
		}
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
		if strings.Contains(query, "?page="+strconv.Itoa(int(currentPage))+"&") {
			// if page is the first query and query also has the other parameters
			query = strings.Replace(query, "page="+strconv.Itoa(int(currentPage))+"&", "", -1)
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

func NewMeta(message string, nodeId string, profileUrl string) *Meta {
	return &Meta{
		Message:    message,
		NodeId:     nodeId,
		ProfileUrl: profileUrl,
	}
}

func NewSearchMeta(message string, numberOfResults int64, totalPages int64) *Meta {
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

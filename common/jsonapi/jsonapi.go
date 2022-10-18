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
	First string `json:"first"`
	Prev  string `json:"prev"`
	Self  string `json:"self"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}

type Meta struct {
	Message         string `json:"message,omitempty"`
	NodeId          string `json:"node_id,omitempty"`
	ProfileUrl      string `json:"profile_url,omitempty"`
	NumberOfResults int64  `json:"number_of_results"`
	TotalPages      int64  `json:"total_pages"`
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

func NewError(titles []string, details []string, sources []string, status []int) []Error {
	var errors []Error
	for i := 0; i < len(titles); i++ {
		if len(details) != 0 && len(sources) != 0 {
			errors = append(errors, Error{
				Status: status[i],
				Title:  titles[i],
				Detail: details[i],
				Source: map[string]string{
					"pointer": sources[i],
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

func NewLinks(c *gin.Context, queryPage int64, currentPage int64, totalPage int64) *Link {
	// check the request is http or https
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	url := scheme + "://" + c.Request.Host

	query := c.Request.RequestURI
	if queryPage != 0 {
		query = strings.Replace(query, "&page="+strconv.Itoa(int(queryPage)), "", -1)
	}

	// if query don't have a question mark, means it didn't have query, we need to add "?" to add the page query.
	// otherwise, it has query, we need to add "&" for adding page later.
	if !strings.Contains(query, "?") {
		query += "?"
	} else {
		query += "&"
	}

	prev := currentPage - 1
	if prev < 1 {
		prev = 1
	}

	next := currentPage + 1
	if next > totalPage {
		next = totalPage
	}

	return &Link{
		First: url + query + "page=1",
		Prev:  url + query + "page=" + strconv.Itoa(int(prev)),
		Self:  url + query + "page=" + strconv.Itoa(int(currentPage)),
		Next:  url + query + "page=" + strconv.Itoa(int(next)),
		Last:  url + query + "page=" + strconv.Itoa(int(totalPage)),
	}
}

func NewMeta(message string, nodeId string, profileUrl string) *Meta {
	return &Meta{
		Message:    message,
		NodeId:     nodeId,
		ProfileUrl: profileUrl,
	}
}

func NewSearchMeta(totalPages int64, numberOfResults int64) *Meta {
	return &Meta{
		NumberOfResults: totalPages,
		TotalPages:      numberOfResults,
	}
}

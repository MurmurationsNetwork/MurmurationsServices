package jsonapi

type JsonApi struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
	Meta   *Meta       `json:"meta,omitempty"`
}

type Error struct {
	Status int               `json:"status,omitempty"`
	Source map[string]string `json:"source,omitempty"`
	Title  string            `json:"title,omitempty"`
	Detail string            `json:"detail,omitempty"`
}

type Meta struct {
	Message string `json:"message,omitempty"`
}

// JSON API Response Combination

func Response(data interface{}, errors []Error, meta *Meta) JsonApi {
	return JsonApi{
		Data:   data,
		Errors: errors,
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

func NewMeta(message string) *Meta {
	return &Meta{
		Message: message,
	}
}

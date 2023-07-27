package main

type Response struct {
	Body       []byte            `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
	empty      bool
}

func NewEmptyResponse() Response {
	return Response{
		Body:       make([]byte, 0),
		Headers:    make(map[string]string),
		StatusCode: -1,
		empty:      true,
	}
}

func NewResponseFromString(body string, status int) Response {
	bodyBytes := []byte(body)

	return Response{
		Body:       bodyBytes,
		Headers:    make(map[string]string),
		StatusCode: status,
		empty:      true,
	}
}

func NewResponse(body []byte, status int) Response {
	return Response{
		Body:       body,
		Headers:    make(map[string]string),
		StatusCode: status,
		empty:      true,
	}
}

func (res *Response) IsEmpty() bool {
	return res.empty
}

func (res *Response) WithContentType(contentType string) Response {
	result := *res
	result.Headers["Content-Type"] = contentType
	return result
}

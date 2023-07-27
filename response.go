package main

type Response struct {
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
	empty      bool
}

func NewEmptyResponse() Response {
	return Response{
		Body:       "",
		Headers:    make(map[string]string),
		StatusCode: -1,
		empty:      true,
	}
}

func NewResponse(body string, status int) Response {
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

package api

import "net/http"

type Response struct {
	StatusCode int
	Body       interface{}
	Header     http.Header
}

func NewResponse(status int, body interface{}) *Response {
	return &Response{
		StatusCode: status,
		Body:       body,
	}
}

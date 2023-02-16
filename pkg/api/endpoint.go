package api

type Endpoint struct {
	Path   string
	Method string
	Auth   bool
	Handle func(r *Request) *Response
}

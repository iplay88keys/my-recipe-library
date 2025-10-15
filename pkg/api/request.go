package api

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Req    *http.Request
	UserID int64
}

func (r *Request) Decode(out interface{}) error {
	defer r.Req.Body.Close()
	err := json.NewDecoder(r.Req.Body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}

package action

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	statusCode int
	result     interface{}
}

func (r *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	return json.NewEncoder(w).Encode(r.result)
}

func Success(result interface{}, status int) *Response {
	return &Response{
		statusCode: status,
		result:     result,
	}
}

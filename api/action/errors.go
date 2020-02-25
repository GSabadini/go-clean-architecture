package action

import (
	"encoding/json"
	"net/http"
)

var (
	ErrInternalServer = Error{statusCode: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrNotFound       = Error{statusCode: http.StatusNotFound, Message: "Not Found"}
)

//Error armazena a estrutura de erro da API
type Error struct {
	statusCode int
	Message    string `json:"message,omitempty"`
}

//Send
func (e Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	return json.NewEncoder(w).Encode(e)
}

//ErrorMessage
func ErrorMessage(err error, status int) *Error {
	return &Error{
		statusCode: status,
		Message:    err.Error(),
	}
}

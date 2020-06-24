package response

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrParameterInvalid = errors.New("parameter invalid")
)

//Error armazena a estrutura de response com error da API
type Error struct {
	statusCode int
	Message    string `json:"message,omitempty"`
}

//Send envia um response de error
func (e Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	return json.NewEncoder(w).Encode(e)
}

//NewError constrói uma estrutura de response com error
func NewError(err error, status int) *Error {
	return &Error{
		statusCode: status,
		Message:    err.Error(),
	}
}

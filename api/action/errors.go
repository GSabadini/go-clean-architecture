package action

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrParameterInvalid = errors.New("Parameter invalid")
)

//Error armazena a estrutura de error da API
type Error struct {
	statusCode int
	Message    string `json:"message,omitempty"`
}

//Send envia uma resposta de error
func (e Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	return json.NewEncoder(w).Encode(e)
}

//ErrorMessage constrói uma estrutura de error
func ErrorMessage(err error, status int) *Error {
	return &Error{
		statusCode: status,
		Message:    err.Error(),
	}
}

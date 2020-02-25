package action

import (
	"encoding/json"
	"net/http"
)

//Response armazena a estrutura de retorno da API
type Response struct {
	statusCode int
	result     interface{}
}

//Send envia uma resposta
func (r *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	return json.NewEncoder(w).Encode(r.result)
}

//Success constrói um response
func Success(result interface{}, status int) *Response {
	return &Response{
		statusCode: status,
		result:     result,
	}
}

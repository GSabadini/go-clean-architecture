package action

import (
	"encoding/json"
	"net/http"
)

//TODO rever nomenclatura
//Response armazena a estrutura de retorno da API
type Response struct {
	statusCode int
	result     interface{}
}

//Send envia uma resposta de sucesso
func (r *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	return json.NewEncoder(w).Encode(r.result)
}

//Success constrói uma estrutura de response
func Success(result interface{}, status int) *Response {
	return &Response{
		statusCode: status,
		result:     result,
	}
}

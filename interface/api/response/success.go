package response

import (
	"encoding/json"
	"net/http"
)

//Success armazena a estrutura de response com sucesso da API
type Success struct {
	statusCode int
	result     interface{}
}

//Send envia um response de sucesso
func (r *Success) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.statusCode)
	return json.NewEncoder(w).Encode(r.result)
}

//NewSuccess constrói uma estrutura de response com sucesso
func NewSuccess(result interface{}, status int) *Success {
	return &Success{
		statusCode: status,
		result:     result,
	}
}

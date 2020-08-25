package input

//Account armazena a estrutura de dados de entrada da API
type Account struct {
	Name    string `json:"name" validate:"required"`
	CPF     string `json:"cpf" validate:"required"`
	Balance int64  `json:"balance" validate:"gt=0,required"`
}

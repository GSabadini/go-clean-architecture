package usecase

//Transfer armazena a estrutura de dados de entrada da API
type TransferInput struct {
	AccountOriginID      string `json:"account_origin_id" validate:"required,uuid4"`
	AccountDestinationID string `json:"account_destination_id" validate:"required,uuid4"`
	Amount               int64  `json:"amount" validate:"gt=0,required"`
}

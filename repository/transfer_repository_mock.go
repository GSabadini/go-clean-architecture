package repository

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"

	"github.com/pkg/errors"
)

//TransferRepositoryMockSuccess implementa a interface de TransferRepository com resultados de sucesso
type TransferRepositoryMockSuccess struct{}

//Store cria uma transferência pré-definida
func (t TransferRepositoryMockSuccess) Store(_ domain.Transfer) (domain.Transfer, error) {
	return domain.Transfer{
		ID:                   "5e570851adcef50116aa7a5a",
		AccountOriginID:      "5e570851adcef50116aa7a5d",
		AccountDestinationID: "5e570851adcef50116aa7a5c",
		Amount:               20,
		CreatedAt:            time.Time{},
	}, nil
}

//FindAll retorna uma lista de transferências pré-definidas
func (t TransferRepositoryMockSuccess) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{
		{
			ID:                   "5e570851adcef50116aa7a5a",
			AccountOriginID:      "5e570851adcef50116aa7a5d",
			AccountDestinationID: "5e570851adcef50116aa7a5c",
			Amount:               100,
			CreatedAt:            time.Time{},
		},
		{
			ID:                   "5e570851adcef50116aa7a5b",
			AccountOriginID:      "5e570851adcef50116aa7a5d",
			AccountDestinationID: "5e570851adcef50116aa7a5c",
			Amount:               500,
			CreatedAt:            time.Time{},
		},
	}, nil
}

//TransferRepositoryMockError implementa a interface de TransferRepository com resultados de erro
type TransferRepositoryMockError struct{}

//Store retorna um error ao criar uma transferência
func (t TransferRepositoryMockError) Store(_ domain.Transfer) (domain.Transfer, error) {
	return domain.Transfer{}, errors.New("Error")
}

//FindAll retorna um error ao listar as transferências
func (t TransferRepositoryMockError) FindAll() ([]domain.Transfer, error) {
	return []domain.Transfer{}, errors.New("Error")
}

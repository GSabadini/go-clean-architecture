package usecase

import (
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferOutput armazena a estrutura de dados de retorno do caso de uso
type TransferOutput struct {
	ID                   string    `json:"id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

//Transfer armazena as dependências para os casos de uso de Transfer
type Transfer struct {
	transferRepo domain.TransferRepository
	accountRepo  domain.AccountRepository
	//sync.Mutex
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransfer(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
) Transfer {
	return Transfer{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
	}
}

//Store cria uma nova Transfer
func (t Transfer) Store(
	accountOriginID,
	accountDestinationID domain.AccountID,
	amount float64,
) (TransferOutput, error) {
	if err := t.process(accountOriginID, accountDestinationID, amount); err != nil {
		return TransferOutput{}, err
	}

	var transfer = domain.NewTransfer(
		domain.TransferID(domain.NewUUID()),
		accountOriginID,
		accountDestinationID,
		amount,
		time.Now(),
	)

	transfer, err := t.transferRepo.Store(transfer)
	if err != nil {
		return TransferOutput{}, err
	}

	return TransferOutput{
		ID:                   string(transfer.ID),
		AccountOriginID:      string(transfer.AccountOriginID),
		AccountDestinationID: string(transfer.AccountDestinationID),
		Amount:               transfer.Amount,
		CreatedAt:            transfer.CreatedAt,
	}, nil
}

/* TODO melhorar processsamento de transação */
func (t Transfer) process(
	accountOriginID,
	accountDestinationID domain.AccountID,
	amount float64,
) error {
	//uc.Lock()
	//defer uc.Unlock()

	origin, err := t.accountRepo.FindByID(accountOriginID)
	if err != nil {
		return err
	}

	destination, err := t.accountRepo.FindByID(accountDestinationID)
	if err != nil {
		return err
	}

	if err := origin.Withdraw(amount); err != nil {
		return err
	}

	destination.Deposit(amount)

	if err = t.accountRepo.UpdateBalance(origin.ID, origin.Balance); err != nil {
		return err
	}

	if err = t.accountRepo.UpdateBalance(destination.ID, destination.Balance); err != nil {
		return err
	}

	return nil
}

//FindAll retorna uma lista de transferências
func (t Transfer) FindAll() ([]TransferOutput, error) {
	var output = make([]TransferOutput, 0)

	transfers, err := t.transferRepo.FindAll()
	if err != nil {
		return output, err
	}

	for _, transfer := range transfers {
		output = append(output, TransferOutput{
			ID:                   string(transfer.ID),
			AccountOriginID:      string(transfer.AccountOriginID),
			AccountDestinationID: string(transfer.AccountDestinationID),
			Amount:               transfer.Amount,
			CreatedAt:            transfer.CreatedAt,
		})
	}

	return output, nil
}

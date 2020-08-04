package usecase

import (
	"context"
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
	ctxTimeout   time.Duration
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransfer(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
	t time.Duration,
) Transfer {
	return Transfer{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		ctxTimeout:   t,
	}
}

//Store cria uma nova Transfer
func (t Transfer) Store(
	ctx context.Context,
	accountOriginID,
	accountDestinationID domain.AccountID,
	amount domain.Money,
) (TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	if err := t.process(ctx, accountOriginID, accountDestinationID, amount); err != nil {
		return TransferOutput{}, err
	}

	var transfer = domain.NewTransfer(
		domain.TransferID(domain.NewUUID()),
		accountOriginID,
		accountDestinationID,
		amount,
		time.Now(),
	)

	transfer, err := t.transferRepo.Store(ctx, transfer)
	if err != nil {
		return TransferOutput{}, err
	}

	return TransferOutput{
		ID:                   transfer.ID.String(),
		AccountOriginID:      transfer.AccountOriginID.String(),
		AccountDestinationID: transfer.AccountDestinationID.String(),
		Amount:               transfer.Amount.Float64(),
		CreatedAt:            transfer.CreatedAt,
	}, nil
}

/* TODO melhorar processsamento de transação */
func (t Transfer) process(
	ctx context.Context,
	accountOriginID,
	accountDestinationID domain.AccountID,
	amount domain.Money,
) error {
	origin, err := t.accountRepo.FindByID(ctx, accountOriginID)
	if err != nil {
		return err
	}

	destination, err := t.accountRepo.FindByID(ctx, accountDestinationID)
	if err != nil {
		return err
	}

	if err := origin.Withdraw(amount); err != nil {
		return err
	}

	destination.Deposit(amount)

	if err = t.accountRepo.UpdateBalance(ctx, origin.ID, origin.Balance); err != nil {
		return err
	}

	if err = t.accountRepo.UpdateBalance(ctx, destination.ID, destination.Balance); err != nil {
		return err
	}

	return nil
}

//FindAll retorna uma lista de transferências
func (t Transfer) FindAll(ctx context.Context) ([]TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	var output = make([]TransferOutput, 0)

	transfers, err := t.transferRepo.FindAll(ctx)
	if err != nil {
		return output, err
	}

	for _, transfer := range transfers {
		output = append(output, TransferOutput{
			ID:                   transfer.ID.String(),
			AccountOriginID:      transfer.AccountOriginID.String(),
			AccountDestinationID: transfer.AccountDestinationID.String(),
			Amount:               transfer.Amount.Float64(),
			CreatedAt:            transfer.CreatedAt,
		})
	}

	return output, nil
}

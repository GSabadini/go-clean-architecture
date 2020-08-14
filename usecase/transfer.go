package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//Transfer armazena as dependências para os casos de uso de Transfer
type Transfer struct {
	transferRepo domain.TransferRepository
	accountRepo  domain.AccountRepository
	presenter    TransferPresenter
	ctxTimeout   time.Duration
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransfer(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
	presenter TransferPresenter,
	t time.Duration,
) Transfer {
	return Transfer{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		presenter:    presenter,
		ctxTimeout:   t,
	}
}

//Store cria uma nova Transfer
func (t Transfer) Store(
	ctx context.Context,
	accountOriginID domain.AccountID,
	accountDestinationID domain.AccountID,
	amount domain.Money,
) (TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	if err := t.process(ctx, accountOriginID, accountDestinationID, amount); err != nil {
		return t.presenter.Output(domain.Transfer{}), err
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
		return t.presenter.Output(domain.Transfer{}), err
	}

	return t.presenter.Output(transfer), nil
}

func (t Transfer) process(
	ctx context.Context,
	accountOriginID domain.AccountID,
	accountDestinationID domain.AccountID,
	amount domain.Money,
) error {
	origin, err := t.accountRepo.FindByID(ctx, accountOriginID)
	if err != nil {
		return err
	}

	if err := origin.Withdraw(amount); err != nil {
		return err
	}

	destination, err := t.accountRepo.FindByID(ctx, accountDestinationID)
	if err != nil {
		return err
	}

	destination.Deposit(amount)

	if err = t.accountRepo.UpdateBalance(ctx, origin.ID(), origin.Balance()); err != nil {
		return err
	}

	if err = t.accountRepo.UpdateBalance(ctx, destination.ID(), destination.Balance()); err != nil {
		return err
	}

	return nil
}

//FindAll retorna uma lista de transferências
func (t Transfer) FindAll(ctx context.Context) ([]TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	transfers, err := t.transferRepo.FindAll(ctx)
	if err != nil {
		return t.presenter.OutputList([]domain.Transfer{}), err
	}

	return t.presenter.OutputList(transfers), nil
}

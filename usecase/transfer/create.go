package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//TransferCreateInteractor armazena as dependências para o caso de uso Create de Transfer
type TransferCreateInteractor struct {
	transferRepo domain.TransferRepository
	accountRepo  domain.AccountRepository
	presenter    TransferPresenter
	ctxTimeout   time.Duration
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransferCreateInteractor(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
	presenter TransferPresenter,
	t time.Duration,
) TransferCreateInteractor {
	return TransferCreateInteractor{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		presenter:    presenter,
		ctxTimeout:   t,
	}
}

//Store cria uma nova Transfer
func (t TransferCreateInteractor) Execute(
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

func (t TransferCreateInteractor) process(
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

package usecase

import (
	"context"
	"github.com/gsabadini/go-bank-transfer/usecase/input"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//CreateTransfer é uma abstração de caso de uso de Transfer
type CreateTransfer interface {
	Execute(context.Context, input.Transfer) (output.TransferOutput, error)
}

//CreateTransferInteractor armazena as dependências para o caso de uso Create de Transfer
type CreateTransferInteractor struct {
	transferRepo domain.TransferRepository
	accountRepo  domain.AccountRepository
	presenter    output.TransferPresenter
	ctxTimeout   time.Duration
}

//NewCreateTransferInteractor constrói um CreateTransferInteractor com suas dependências
func NewCreateTransferInteractor(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
	presenter output.TransferPresenter,
	t time.Duration,
) CreateTransferInteractor {
	return CreateTransferInteractor{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		presenter:    presenter,
		ctxTimeout:   t,
	}
}

//Execute cria uma nova Transfer
func (t CreateTransferInteractor) Execute(ctx context.Context, input input.Transfer) (output.TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	if err := t.process(ctx, input); err != nil {
		return t.presenter.Output(domain.Transfer{}), err
	}

	var transfer = domain.NewTransfer(
		domain.TransferID(domain.NewUUID()),
		domain.AccountID(input.AccountOriginID),
		domain.AccountID(input.AccountDestinationID),
		domain.Money(input.Amount),
		time.Now(),
	)

	transfer, err := t.transferRepo.Store(ctx, transfer)
	if err != nil {
		return t.presenter.Output(domain.Transfer{}), err
	}

	return t.presenter.Output(transfer), nil
}

func (t CreateTransferInteractor) process(ctx context.Context, input input.Transfer) error {
	origin, err := t.accountRepo.FindByID(ctx, domain.AccountID(input.AccountOriginID))
	if err != nil {
		return err
	}

	if err := origin.Withdraw(domain.Money(input.Amount)); err != nil {
		return err
	}

	destination, err := t.accountRepo.FindByID(ctx, domain.AccountID(input.AccountDestinationID))
	if err != nil {
		return err
	}

	destination.Deposit(domain.Money(input.Amount))

	if err = t.accountRepo.UpdateBalance(ctx, origin.ID(), origin.Balance()); err != nil {
		return err
	}

	if err = t.accountRepo.UpdateBalance(ctx, destination.ID(), destination.Balance()); err != nil {
		return err
	}

	return nil
}

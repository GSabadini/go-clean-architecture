package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type (
	// CreateTransferUseCase input port
	CreateTransferUseCase interface {
		Execute(context.Context, CreateTransferInput) (CreateTransferOutput, error)
	}

	// CreateTransferInput input data
	CreateTransferInput struct {
		AccountOriginID      string `json:"account_origin_id" validate:"required,uuid4"`
		AccountDestinationID string `json:"account_destination_id" validate:"required,uuid4"`
		Amount               int64  `json:"amount" validate:"gt=0,required"`
	}

	// CreateTransferPresenter output port
	CreateTransferPresenter interface {
		Output(domain.Transfer) CreateTransferOutput
	}

	// CreateTransferOutput output data
	CreateTransferOutput struct {
		ID                   string  `json:"id"`
		AccountOriginID      string  `json:"account_origin_id"`
		AccountDestinationID string  `json:"account_destination_id"`
		Amount               float64 `json:"amount"`
		CreatedAt            string  `json:"created_at"`
	}

	createTransferInteractor struct {
		transferRepo domain.TransferRepository
		accountRepo  domain.AccountRepository
		presenter    CreateTransferPresenter
		ctxTimeout   time.Duration
	}
)

// NewCreateTransferInteractor creates new createTransferInteractor with its dependencies
func NewCreateTransferInteractor(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
	presenter CreateTransferPresenter,
	t time.Duration,
) CreateTransferUseCase {
	return createTransferInteractor{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		presenter:    presenter,
		ctxTimeout:   t,
	}
}

// Execute orchestrates the use case
func (t createTransferInteractor) Execute(ctx context.Context, input CreateTransferInput) (CreateTransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	var (
		transfer domain.Transfer
		err      error
	)

	err = t.transferRepo.WithTransaction(ctx, func(ctxTx context.Context) error {
		if err = t.process(ctxTx, input); err != nil {
			return err
		}

		transfer = domain.NewTransfer(
			domain.TransferID(domain.NewUUID()),
			domain.AccountID(input.AccountOriginID),
			domain.AccountID(input.AccountDestinationID),
			domain.Money(input.Amount),
			time.Now(),
		)

		transfer, err = t.transferRepo.Create(ctxTx, transfer)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return t.presenter.Output(domain.Transfer{}), err
	}

	return t.presenter.Output(transfer), nil
}

func (t createTransferInteractor) process(ctx context.Context, input CreateTransferInput) error {
	origin, err := t.accountRepo.FindByID(ctx, domain.AccountID(input.AccountOriginID))
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			return domain.ErrAccountOriginNotFound
		default:
			return err
		}
	}

	if err := origin.Withdraw(domain.Money(input.Amount)); err != nil {
		return err
	}

	destination, err := t.accountRepo.FindByID(ctx, domain.AccountID(input.AccountDestinationID))
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			return domain.ErrAccountDestinationNotFound
		default:
			return err
		}
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

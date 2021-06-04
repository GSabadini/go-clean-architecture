package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

type (
	// FindAllTransferUseCase input port
	FindAllTransferUseCase interface {
		Execute(context.Context) ([]FindAllTransferOutput, error)
	}

	// FindAllTransferPresenter output port
	FindAllTransferPresenter interface {
		Output([]domain.Transfer) []FindAllTransferOutput
	}

	// FindAllTransferOutput output data
	FindAllTransferOutput struct {
		ID                   string  `json:"id"`
		AccountOriginID      string  `json:"account_origin_id"`
		AccountDestinationID string  `json:"account_destination_id"`
		Amount               float64 `json:"amount"`
		CreatedAt            string  `json:"created_at"`
	}

	findAllTransferInteractor struct {
		repo       domain.TransferRepository
		presenter  FindAllTransferPresenter
		ctxTimeout time.Duration
	}
)

// NewFindAllTransferInteractor creates new findAllTransferInteractor with its dependencies
func NewFindAllTransferInteractor(
	repo domain.TransferRepository,
	presenter FindAllTransferPresenter,
	t time.Duration,
) FindAllTransferUseCase {
	return findAllTransferInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (t findAllTransferInteractor) Execute(ctx context.Context) ([]FindAllTransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	transfers, err := t.repo.FindAll(ctx)
	if err != nil {
		return t.presenter.Output([]domain.Transfer{}), err
	}

	return t.presenter.Output(transfers), nil
}

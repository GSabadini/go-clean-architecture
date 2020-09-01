package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
)

type FindAllTransfer interface {
	Execute(context.Context) ([]output.Transfer, error)
}

type FindAllTransferInteractor struct {
	repo       domain.TransferRepository
	presenter  output.TransferPresenter
	ctxTimeout time.Duration
}

func NewFindAllTransferInteractor(
	repo domain.TransferRepository,
	presenter output.TransferPresenter,
	t time.Duration,
) FindAllTransferInteractor {
	return FindAllTransferInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

func (t FindAllTransferInteractor) Execute(ctx context.Context) ([]output.Transfer, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	transfers, err := t.repo.FindAll(ctx)
	if err != nil {
		return t.presenter.OutputList([]domain.Transfer{}), err
	}

	return t.presenter.OutputList(transfers), nil
}

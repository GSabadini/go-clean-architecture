package usecase

import (
	"context"
	"github.com/gsabadini/go-bank-transfer/usecase/output"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//FindAllTransfer é uma abstração de caso de uso de Transfer
type FindAllTransfer interface {
	Execute(context.Context) ([]output.TransferOutput, error)
}

//FindAllTransferInteractor armazena as dependências para o casos de uso FindAll de Transfer
type FindAllTransferInteractor struct {
	repo       domain.TransferRepository
	presenter  output.TransferPresenter
	ctxTimeout time.Duration
}

//NewFindAllTransferInteractor constrói um Transfer com suas dependências
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

//Execute retorna uma lista de transferências
func (t FindAllTransferInteractor) Execute(ctx context.Context) ([]output.TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	transfers, err := t.repo.FindAll(ctx)
	if err != nil {
		return t.presenter.OutputList([]domain.Transfer{}), err
	}

	return t.presenter.OutputList(transfers), nil
}

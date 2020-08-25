package usecase

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
)

//FindAllInteractor armazena as dependências para o casos de uso FindAll de Transfer
type TransferFindAllInteractor struct {
	transferRepo domain.TransferRepository
	presenter    TransferPresenter
	ctxTimeout   time.Duration
}

//NewTransfer constrói um Transfer com suas dependências
func NewTransferFindAllInteractor(
	transferRepo domain.TransferRepository,
	presenter TransferPresenter,
	t time.Duration,
) TransferFindAllInteractor {
	return TransferFindAllInteractor{
		transferRepo: transferRepo,
		presenter:    presenter,
		ctxTimeout:   t,
	}
}

//Execute retorna uma lista de transferências
func (t TransferFindAllInteractor) Execute(ctx context.Context) ([]TransferOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	transfers, err := t.transferRepo.FindAll(ctx)
	if err != nil {
		return t.presenter.OutputList([]domain.Transfer{}), err
	}

	return t.presenter.OutputList(transfers), nil
}

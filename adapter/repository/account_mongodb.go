package repository

import (
	"context"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountBSON struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	CPF       string    `bson:"cpf"`
	Balance   int64     `bson:"balance"`
	CreatedAt time.Time `bson:"created_at"`
}

type AccountNoSQL struct {
	collectionName string
	db             NoSQL
}

func NewAccountNoSQL(db NoSQL) AccountNoSQL {
	return AccountNoSQL{
		db:             db,
		collectionName: "accounts",
	}
}

func (a AccountNoSQL) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	var accountBSON = accountBSON{
		ID:        account.ID().String(),
		Name:      account.Name(),
		CPF:       account.CPF(),
		Balance:   account.Balance().Int64(),
		CreatedAt: account.CreatedAt(),
	}

	if err := a.db.Store(ctx, a.collectionName, accountBSON); err != nil {
		return domain.Account{}, errors.Wrap(err, "error creating account")
	}

	return account, nil
}

func (a AccountNoSQL) UpdateBalance(ctx context.Context, ID domain.AccountID, balance domain.Money) error {
	var (
		query  = bson.M{"id": ID}
		update = bson.M{"$set": bson.M{"balance": balance}}
	)

	if err := a.db.Update(ctx, a.collectionName, query, update); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return errors.Wrap(domain.ErrAccountNotFound, "error updating account balance")
		default:
			return errors.Wrap(err, "error updating account balance")
		}
	}

	return nil
}

func (a AccountNoSQL) FindAll(ctx context.Context) ([]domain.Account, error) {
	var accountsBSON = make([]accountBSON, 0)

	if err := a.db.FindAll(ctx, a.collectionName, bson.M{}, &accountsBSON); err != nil {
		switch err {
		case mongo.ErrNilDocument:
			return []domain.Account{}, errors.Wrap(domain.ErrAccountNotFound, "error listing accounts")
		default:
			return []domain.Account{}, errors.Wrap(err, "error listing accounts")
		}
	}

	var accounts = make([]domain.Account, 0)

	for _, accountBSON := range accountsBSON {
		var account = domain.NewAccount(
			domain.AccountID(accountBSON.ID),
			accountBSON.Name,
			accountBSON.CPF,
			domain.Money(accountBSON.Balance),
			accountBSON.CreatedAt,
		)

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (a AccountNoSQL) FindByID(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		accountBSON = &accountBSON{}
		query       = bson.M{"id": ID}
	)

	if err := a.db.FindOne(ctx, a.collectionName, query, nil, accountBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.Account{}, domain.ErrAccountNotFound
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account")
		}
	}

	return domain.NewAccount(
		domain.AccountID(accountBSON.ID),
		accountBSON.Name,
		accountBSON.CPF,
		domain.Money(accountBSON.Balance),
		accountBSON.CreatedAt,
	), nil
}

func (a AccountNoSQL) FindBalance(ctx context.Context, ID domain.AccountID) (domain.Account, error) {
	var (
		accountBSON = &accountBSON{}
		query       = bson.M{"id": ID}
		projection  = bson.M{"balance": 1, "_id": 0}
	)

	if err := a.db.FindOne(ctx, a.collectionName, query, projection, accountBSON); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return domain.Account{}, domain.ErrAccountNotFound
		default:
			return domain.Account{}, errors.Wrap(err, "error fetching account balance")
		}
	}

	return domain.NewAccountBalance(domain.Money(accountBSON.Balance)), nil
}

package usecase

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/repository"
	"reflect"
	"testing"
	"time"
)

func TestStoreTransfer(t *testing.T) {
	type args struct {
		transferRepository repository.TransferRepository
		accountRepository  repository.AccountRepository
		transfer           domain.Transfer
	}
	tests := []struct {
		name     string
		args     args
		expected domain.Transfer
		wantErr  bool
	}{
		{
			name: "Success store transfer",
			args: args{
				transferRepository: repository.TransferRepositoryMockSuccess{},
				accountRepository:  repository.AccountRepositoryMock{},
				transfer: domain.Transfer{
					AccountOriginID:      "5e570851adcef50116aa7a5d",
					AccountDestinationID: "5e570851adcef50116aa7a5c",
					Amount:               20,
				},
			},
			expected: domain.Transfer{
				ID:                   "5e570851adcef50116aa7a5a",
				AccountOriginID:      "5e570851adcef50116aa7a5d",
				AccountDestinationID: "5e570851adcef50116aa7a5c",
				Amount:               100,
				CreatedAt:            time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StoreTransfer(tt.args.transferRepository, tt.args.accountRepository, tt.args.transfer)

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) || got != tt.expected {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}

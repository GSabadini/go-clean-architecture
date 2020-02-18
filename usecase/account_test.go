package usecase

import (
	"testing"

	"github.com/gsabadini/go-stone/domain"
	"github.com/gsabadini/go-stone/infrastructure/database"
	"github.com/gsabadini/go-stone/repository"
)

func TestCreate(t *testing.T) {
	type args struct {
		repository repository.Account
		account    domain.Account
	}

	tests := []struct {
		name     string
		args     args
		expected interface{}
	}{
		{
			name: "Create account successful",
			args: args{
				repository: repository.NewAccount(database.MongoHandlerSuccessMock{}),
				account: domain.Account{
					Id:        "1",
					Name:      "Test",
					Cpf:       "123123123",
					Ballance:  0,
					CreatedAt: "01/01/01",
				},
			},
			expected: nil,
		},
		//{
		//	name: "Create account error",
		//	args: args{
		//		repository: repository.NewAccount(database.MongoHandlerErrorMock{}),
		//		account: domain.Account{
		//			Id:        "1",
		//			Name:      "Test",
		//			Cpf:       "123123123",
		//			Ballance:  0,
		//			CreatedAt: "01/01/01",
		//		},
		//	},
		//	expected: "Fail",
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := Create(tt.args.repository, tt.args.account); err != nil {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, err, tt.expected)
			}
		})
	}
}

package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"reflect"
	"testing"
	"time"
)

func Test_createAccountPresenter_Output(t *testing.T) {
	type args struct {
		account domain.Account
	}
	tests := []struct {
		name string
		args args
		want usecase.CreateAccountOutput
	}{
		{
			name: "Create account output",
			args: args{
				account: domain.NewAccount(
					"3c096a40-ccba-4b58-93ed-57379ab04680",
					"Testing",
					"07091054965",
					1000,
					time.Time{},
				),
			},
			want: usecase.CreateAccountOutput{
				ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
				Name:      "Testing",
				CPF:       "07091054965",
				Balance:   10,
				CreatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateAccountPresenter()
			if got := pre.Output(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

package presenter

import (
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

func Test_findAllAccountPresenter_Output(t *testing.T) {
	type args struct {
		accounts []domain.Account
	}
	tests := []struct {
		name string
		args args
		want []usecase.FindAllAccountOutput
	}{
		{
			name: "Find all account output",
			args: args{
				accounts: []domain.Account{
					domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04680",
						"Testing",
						"07091054965",
						1000,
						time.Time{},
					),
					domain.NewAccount(
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						"Testing",
						"07091054965",
						99,
						time.Time{},
					),
				},
			},
			want: []usecase.FindAllAccountOutput{
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04680",
					Name:      "Testing",
					CPF:       "07091054965",
					Balance:   10,
					CreatedAt: "0001-01-01T00:00:00Z",
				},
				{
					ID:        "3c096a40-ccba-4b58-93ed-57379ab04682",
					Name:      "Testing",
					CPF:       "07091054965",
					Balance:   0.99,
					CreatedAt: "0001-01-01T00:00:00Z",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewFindAllAccountPresenter()
			if got := pre.Output(tt.args.accounts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"reflect"
	"testing"
)

func TestNewFindAccountBalancePresenter(t *testing.T) {
	tests := []struct {
		name string
		want usecase.FindAccountBalancePresenter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFindAccountBalancePresenter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFindAccountBalancePresenter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findAccountBalancePresenter_Output(t *testing.T) {
	type args struct {
		balance domain.Money
	}
	tests := []struct {
		name string
		args args
		want usecase.FindAccountBalanceOutput
	}{
		{
			name: "Find account balance ouitput",
			args: args{
				balance: 1099,
			},
			want: usecase.FindAccountBalanceOutput{
				Balance: 10.99,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewFindAccountBalancePresenter()
			if got := pre.Output(tt.args.balance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

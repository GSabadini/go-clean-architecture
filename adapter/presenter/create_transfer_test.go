package presenter

import (
	"reflect"
	"testing"
	"time"

	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
)

func Test_createTransferPresenter_Output(t *testing.T) {
	type args struct {
		transfer domain.Transfer
	}
	tests := []struct {
		name string
		args args
		want usecase.CreateTransferOutput
	}{
		{
			name: "Create transfer output",
			args: args{
				transfer: domain.NewTransfer(
					"3c096a40-ccba-4b58-93ed-57379ab04680",
					"3c096a40-ccba-4b58-93ed-57379ab04681",
					"3c096a40-ccba-4b58-93ed-57379ab04682",
					1000,
					time.Time{},
				),
			},
			want: usecase.CreateTransferOutput{
				ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
				AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
				AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
				Amount:               10,
				CreatedAt:            "0001-01-01T00:00:00Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateTransferPresenter()
			if got := pre.Output(tt.args.transfer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

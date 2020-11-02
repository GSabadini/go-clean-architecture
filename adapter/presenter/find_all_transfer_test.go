package presenter

import (
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/usecase"
	"reflect"
	"testing"
	"time"
)

func Test_findAllTransferPresenter_Output(t *testing.T) {
	type args struct {
		transfers []domain.Transfer
	}
	tests := []struct {
		name string
		args args
		want []usecase.FindAllTransferOutput
	}{
		{
			name: "Find all transfer output",
			args: args{
				transfers: []domain.Transfer{
					domain.NewTransfer(
						"3c096a40-ccba-4b58-93ed-57379ab04680",
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						1000,
						time.Time{},
					),
					domain.NewTransfer(
						"3c096a40-ccba-4b58-93ed-57379ab04680",
						"3c096a40-ccba-4b58-93ed-57379ab04681",
						"3c096a40-ccba-4b58-93ed-57379ab04682",
						99,
						time.Time{},
					),
				},
			},
			want: []usecase.FindAllTransferOutput{
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               10,
					CreatedAt:            "0001-01-01T00:00:00Z",
				},
				{
					ID:                   "3c096a40-ccba-4b58-93ed-57379ab04680",
					AccountOriginID:      "3c096a40-ccba-4b58-93ed-57379ab04681",
					AccountDestinationID: "3c096a40-ccba-4b58-93ed-57379ab04682",
					Amount:               0.99,
					CreatedAt:            "0001-01-01T00:00:00Z",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewFindAllTransferPresenter()
			if got := pre.Output(tt.args.transfers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

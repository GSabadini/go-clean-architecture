package usecase

//AccountUseCase é uma abstração para os casos de uso de Account
type AccountUseCase interface {
	Store(string, string, float64) (AccountOutput, error)
	FindAll() ([]AccountOutput, error)
	FindBalance(string) (AccountBalanceOutput, error)
}

//TransferUseCase é uma abstração para os casos de uso de Transfer
type TransferUseCase interface {
	Store(string, string, float64) (TransferOutput, error)
	FindAll() ([]TransferOutput, error)
}

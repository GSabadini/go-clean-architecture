package usecase

//AccountUseCase é uma abstração para os casos de uso de Account
type AccountUseCase interface {
	Store(string, string, float64) (accountOutput, error)
	FindAll() ([]accountOutput, error)
	FindBalance(string) (accountBalanceOutput, error)
}

//TransferUseCase é uma abstração para os casos de uso de Transfer
type TransferUseCase interface {
	Store(string, string, float64) (transferOutput, error)
	FindAll() ([]transferOutput, error)
}

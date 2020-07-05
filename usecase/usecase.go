package usecase

type AccountUseCase interface {
	Store(string, string, float64) (accountOutput, error)
	FindAll() ([]accountOutput, error)
	FindBalance(string) (accountBalanceOutput, error)
}

type TransferUseCase interface {
	Store(string, string, float64) (transferOutput, error)
	FindAll() ([]transferOutput, error)
}

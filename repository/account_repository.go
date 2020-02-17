package repository

type AccountRepository struct {
	dbHandler string
}

func NewAccountRepository(dbHandler string) AccountRepository {
	return AccountRepository{dbHandler: dbHandler}
}

func (a AccountRepository) Store() {

}


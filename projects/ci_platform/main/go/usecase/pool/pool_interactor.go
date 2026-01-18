package usePool

type interactor struct {
	pool Pool
}

func New(pool Pool) Interactor {
	return &interactor{pool}
}

func (interactor *interactor) Pool() Pool {
	return interactor.pool
}

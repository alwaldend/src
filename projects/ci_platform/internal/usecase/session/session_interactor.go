package useSession

type interactor struct {
	repository Repository
}

func New(repository Repository) Interactor {
	return &interactor{repository: repository}
}

func (interactor *interactor) Get(context interface{}) (Session, error) {
	return interactor.repository.Get(context)
}

func (interactor *interactor) Reset() error {
	return interactor.repository.Reset()
}

func (interactor *interactor) Close() error {
	return interactor.repository.Close()
}

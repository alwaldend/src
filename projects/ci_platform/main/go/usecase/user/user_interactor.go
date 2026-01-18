package useUser

import (
	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg"
)

type userInteractor struct {
	repository Repository
}

func New(repository Repository) Interactor {
	return &userInteractor{
		repository: repository,
	}
}

func (interactor *userInteractor) Validate(user *entity.User) error {
	return pkg.ValidateStruct(user)
}

func (interactor *userInteractor) Create(user *entity.User) error {
	if err := pkg.ValidateStruct(user); err != nil {
		return err
	}
	err := pkg.ValidateStruct(&entity.Password{Value: user.Password})
	if err != nil {
		return err
	}
	user.Password = interactor.PasswordHashCreate(user.Password)
	return interactor.repository.Create(user)
}

func (interactor *userInteractor) Update(user *entity.User) error {
	err := interactor.Validate(user)
	if err != nil {
		return err
	}
	return interactor.repository.Save(user)
}

func (interactor *userInteractor) GetAll() ([]entity.User, error) {
	users, err := interactor.repository.GetAll()
	return users, err
}

func (interactor *userInteractor) Get(id uuid.UUID) (*entity.User, error) {
	return interactor.repository.Get(id)
}

func (interactor *userInteractor) FindByUsername(username string) (
	*entity.User, error,
) {
	return interactor.repository.FindOne(&entity.User{Username: username})
}

func (interactor *userInteractor) PasswordHashCheck(
	user *entity.User, password string,
) (bool, error) {
	return pkg.EncodedHashCompare(password, user.Password)
}

func (interactor *userInteractor) PasswordHashCreate(
	password string,
) string {
	return pkg.EncodedHashGenerate(password)
}

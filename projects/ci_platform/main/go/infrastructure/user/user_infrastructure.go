package infraUser

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	usecaseUser "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/user"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) usecaseUser.Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]entity.User, error) { return r.Find() }

func (r *repository) FindOne(condition *entity.User) (*entity.User, error) {
	user := &entity.User{}
	return user, r.db.First(&user, condition).Error
}

func (r *repository) Get(id uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	return user, r.db.First(&user, id).Error
}

func (r *repository) Find(
	conditions ...interface{},
) ([]entity.User, error) {
	users := make([]entity.User, 0)
	err := r.db.Find(&users, conditions...).Error
	return users, err
}

func (r *repository) Save(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *repository) Delete(user *entity.User) error {
	return r.db.Delete(user).Error
}

func (r *repository) Migrate() error {
	return r.db.AutoMigrate(
		&entity.Capability{}, &entity.Permission{},
		&entity.Role{}, &entity.User{},
	)
}

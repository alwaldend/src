package useUser

import (
	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
)

// interactor for entity.User
type Interactor interface {
	// check whether a given user is valid
	Validate(user *entity.User) error
	// create a user
	// password will be hashed
	Create(user *entity.User) error
	// update a user
	Update(user *entity.User) error
	// return all users in the database
	GetAll() ([]entity.User, error)
	// get one by id
	Get(id uuid.UUID) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)

	// check whether a given unhashed password matches the hashed password of
	// the user
	PasswordHashCheck(user *entity.User, password string) (
		matches bool, err error,
	)
	// creates a hashed password from the given string
	PasswordHashCreate(password string) (hashedPassword string)
}

// repository for entity.User
type Repository interface {
	// return all users in the database
	GetAll() ([]entity.User, error)
	// return a user based on the condition
	// repository.FindOne(&model.User{Username: "username0"})
	FindOne(condition *entity.User) (*entity.User, error)
	// get by id
	Get(id uuid.UUID) (*entity.User, error)
	// return an array of users based on the conditions
	// https://gorm.io/docs/query.html
	Find(conditions ...interface{}) ([]entity.User, error)

	// save a given user (if it doesn't exist, create)
	Create(user *entity.User) error
	// save a given user (if it doesn't exist, create)
	Save(user *entity.User) error
	// delete a given user
	Delete(user *entity.User) error

	// automatically migrate userRegistry tables
	// https://gorm.io/docs/migration.html
	// NOTE: AutoMigrate will create tables, missing foreign keys, constraints,
	// columns and indexes
	// It will change existing column’s type if its size, precision, nullable
	// changed
	// It WON’T delete unused columns to protect your data
	Migrate() error
}

package useJob

import (
	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
)

type Interactor interface {
	Repository() Repository
}

// repository for entity.Job
type Repository interface {
	// return all jobs in the database
	GetAll() ([]entity.Job, error)
	// return a job based on the condition
	// repository.FindOne(&entity.Job{ID: id})
	FindOne(condition *entity.Job) (*entity.Job, error)
	// get by id
	Get(id uuid.UUID) (*entity.Job, error)
	// get by redis id
	GetByRedisID(id string) (*entity.Job, error)
	// return an array of jobs based on conditions
	// https://gorm.io/docs/query.html
	Find(conditions ...interface{}) ([]entity.Job, error)

	// save a given job (if it doesn't exist, create)
	Create(job *entity.Job) error
	// create a given job
	Save(job *entity.Job) error
	// delete a job
	Delete(job *entity.Job) error

	// automatically migrate job tables
	// https://gorm.io/docs/migration.html
	// NOTE: AutoMigrate will create tables, missing foreign keys, constraints,
	// columns and indexes
	// It will change existing column’s type if its size, precision, nullable
	// changed
	// It WON’T delete unused columns to protect your data
	Migrate() error
}

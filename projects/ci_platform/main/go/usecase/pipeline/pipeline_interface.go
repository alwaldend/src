package usePipeline

import (
	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
)

const (
	InteractorKey = "usecasePipeline.Interactor"
)

type Interactor interface {
	Repository() Repository
	RunJob(job *entity.Job) error
	SchedulePipelines(branch *entity.Branch) ([]entity.Pipeline, error)
	JobErrorHandler(job *entity.Job, err error) error
}

// repository for entity.Pipeline
type Repository interface {
	// return all pipelines in the database
	GetAll() ([]entity.Pipeline, error)
	// return a pipeline based on the condition
	// repository.FindOne(&entity.Pipeline{ID: id})
	FindOne(condition *entity.Pipeline) (*entity.Pipeline, error)
	// get by id
	Get(id uuid.UUID) (*entity.Pipeline, error)
	// return an array of pipelines based on conditions
	// https://gorm.io/docs/query.html
	Find(conditions ...interface{}) ([]entity.Pipeline, error)

	// save a given Pipeline (if it doesn't exist, create)
	Create(pipeline *entity.Pipeline) error
	// create a given Pipeline
	Save(pipeline *entity.Pipeline) error
	// delete a pipeline
	Delete(pipeline *entity.Pipeline) error

	// automatically migrate Pipeline tables
	// https://gorm.io/docs/migration.html
	// NOTE: AutoMigrate will create tables, missing foreign keys, constraints,
	// columns and indexes
	// It will change existing column’s type if its size, precision, nullable
	// changed
	// It WON’T delete unused columns to protect your data
	Migrate() error
}

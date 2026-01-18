package useProject

import (
	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
)

// interactor for entity.project
type Interactor interface {
	// check whether a given project is valid
	Validate(project *entity.Project) error
	// create a project
	Create(project *entity.Project) error
	// update a project
	Update(project *entity.Project) error
	// return all projects in the database
	GetAll() ([]entity.Project, error)
	// get one by id
	Get(id uuid.UUID) (*entity.Project, error)
	GetByName(name string) (*entity.Project, error)
	GetByPath(path string) (*entity.Project, error)
}

type Repository interface {
	// return all projects in the database
	GetAll() ([]entity.Project, error)
	// return a project based on the condition
	// repository.FindOne(&model.project{projectname: "projectname0"})
	FindOne(condition *entity.Project) (*entity.Project, error)
	// get by id
	Get(id uuid.UUID) (*entity.Project, error)

	Create(project *entity.Project) error
	// save a given project (if it doesn't exist, create)
	Save(project *entity.Project) error

	// update a project
	Update(project *entity.Project) error
	// delete a given project
	Delete(project *entity.Project) error

	Migrate() error
}

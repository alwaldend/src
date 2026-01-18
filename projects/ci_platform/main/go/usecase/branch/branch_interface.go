package useBranch

import (
	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
)

// interactor for entity.branch
type Interactor interface {
	Validate(branch *entity.Branch) error
	Create(branch *entity.Branch) error
	Update(branch *entity.Branch) error
	GetProjectBranch(projectID uuid.UUID, name string) (*entity.Branch, error)
	GetProjectBranches(projectID uuid.UUID) ([]entity.Branch, error)
	GetAll() ([]entity.Branch, error)
	Get(id uuid.UUID) (*entity.Branch, error)
	// get a branch from the remote source
	// clone the branch from source with depth 1, then parse the file
	GetBranchFromRemote(project *entity.Project, branch ...string) (
		*entity.Branch, error,
	)
	// update the branch
	UpdateBranchFromRemote(project *entity.Project, branchName string) (
		branch *entity.Branch, err error,
	)
}

type Repository interface {
	// clone the branch from source with depth 1, then read the file
	GetRemotePipelineConfig(source string, branch string, filePath string) (
		config *entity.PipelineConfig, err error,
	)

	// return all branches in the database
	GetAll() ([]entity.Branch, error)
	// return a branch based on the condition
	// repository.FindOne(&entity.Branch{Name: "branchname0"})
	FindOne(condition *entity.Branch) (*entity.Branch, error)
	// get by id
	Get(id uuid.UUID) (*entity.Branch, error)
	// return an array of branches based on the conditions
	// https://gorm.io/docs/query.html
	Find(conditions ...interface{}) ([]entity.Branch, error)

	// save a given branch (if it doesn't exist, create)
	Save(branch *entity.Branch) error

	// update a branch
	Update(branch *entity.Branch) error
	// delete a given branch
	Delete(branch *entity.Branch) error

	Migrate() error
}

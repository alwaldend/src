package useBranch

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg"
)

type interactor struct {
	repository Repository
}

func New(repository Repository) Interactor {
	return &interactor{
		repository: repository,
	}
}

func (interactor *interactor) Validate(branch *entity.Branch) error {
	return pkg.ValidateStruct(branch)
}

func (interactor *interactor) Create(branch *entity.Branch) error {
	if err := pkg.ValidateStruct(branch); err != nil {
		return err
	}
	return interactor.repository.Save(branch)
}

func (interactor *interactor) Update(branch *entity.Branch) error {
	if err := pkg.ValidateStruct(branch); err != nil {
		return err
	}
	return interactor.repository.Update(branch)
}

func (interactor *interactor) GetAll() ([]entity.Branch, error) {
	return interactor.repository.GetAll()
}

func (interactor *interactor) Get(id uuid.UUID) (*entity.Branch, error) {
	return interactor.repository.Get(id)
}

func (interactor *interactor) GetProjectBranch(projectID uuid.UUID, name string) (
	*entity.Branch, error,
) {
	return interactor.repository.FindOne(
		&entity.Branch{Name: name, ProjectID: projectID},
	)
}

func (interactor *interactor) GetProjectBranches(projectID uuid.UUID) (
	[]entity.Branch, error,
) {
	return interactor.repository.Find(&entity.Branch{ProjectID: projectID})
}

func (interactor *interactor) GetBranchFromRemote(
	project *entity.Project, branch ...string,
) (*entity.Branch, error) {
	var branchName string
	if len(branch) != 0 {
		branchName = branch[0]
	} else {
		branchName = project.DefaultBranch
	}
	if project.NamespaceID == nil {
		project.Path = strings.ToLower(project.Name)
	} else {
		return nil, errors.New("namespaces are not implemented yet")
		// project.Path = project.Namespace.Path + "/" + project.Name
	}
	config, err := interactor.repository.GetRemotePipelineConfig(
		project.Source, branchName, ".ci_platform/config.hcl",
	)
	if err != nil {
		return nil, err
	}
	return &entity.Branch{Name: branchName, Config: config}, nil
}

func (interactor *interactor) UpdateBranchFromRemote(
	project *entity.Project, branchName string,
) (branch *entity.Branch, err error) {
	branch, err = interactor.GetProjectBranch(project.ID, branchName)
	if err != nil {
		return
	}
	branchNew, err := interactor.GetBranchFromRemote(project, branchName)
	if err != nil {
		return
	}
	branch.Config = branchNew.Config
	interactor.Update(branch)
	return
}

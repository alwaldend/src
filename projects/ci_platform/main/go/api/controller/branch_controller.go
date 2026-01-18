package controller

import (
	"github.com/gofiber/fiber/v2"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	useBranch "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/branch"
	useProject "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/project"
)

type Branch interface {
	Get(ctx *fiber.Ctx) error
	GetProjectBranches(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type branchController struct {
	branch  useBranch.Interactor
	project useProject.Interactor
}

func NewBranch(
	branch useBranch.Interactor, project useProject.Interactor,
) Branch {
	return &branchController{branch: branch, project: project}
}

type (
	ResponseBranch   presenter.Response[entity.Branch]
	ResponseBranches presenter.Response[[]entity.Branch]
)

// get a project branch
// @Summary get a project branch
// @Description get a project branch
// @ID get-branch
// @Tags branch
// @Security ApiKeyAuth
//
// @Param   project_id  path     string  true  "project id"
// @Param   branch_name  path     string  true  "branch name"
//
// @Success 200 {object} ResponseBranch
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id}/branches/{branch_name} [GET]
func (c *branchController) Get(ctx *fiber.Ctx) error {
	id, err := getProjectID(ctx)
	if err != nil {
		return err
	}
	name, err := getBranchName(ctx)
	if err != nil {
		return err
	}
	branch, err := c.branch.GetProjectBranch(id, name)
	if err != nil {
		return presenter.CouldNotFindProjectBranch(err)
	}
	return presenter.Success(ctx, branch)
}

// get project branches
// @Summary get all project branches
// @Description get all branches
// @ID get-branches
// @Tags branch
// @Security ApiKeyAuth
//
// @Param   project_id  path     string  true  "project id"
//
// @Success 200 {object} ResponseBranches
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id}/branches [GET]
func (c *branchController) GetProjectBranches(ctx *fiber.Ctx) error {
	id, err := getProjectID(ctx)
	if err != nil {
		return err
	}
	branches, err := c.branch.GetProjectBranches(id)
	if err != nil {
		return presenter.CouldNotFindProjectBranch(err)
	}
	return presenter.Success(ctx, branches)
}

type RequestBranchCreate struct {
	Name string `json:"name" form:"name"  validate:"required,min=3,max=64"`
}

// create a branch
// @Summary create a branch
// @Description create a branch
// @ID create-branch
// @Tags branch
// @Security ApiKeyAuth
//
// @Param project_id   path string                 true  "project id"
// @Param Body body RequestBranchCreate true  "creation request body"
//
// @Success 200 {object} ResponseBranch
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id}/branches [POST]
func (c *branchController) Create(ctx *fiber.Ctx) error {
	id, err := getProjectID(ctx)
	if err != nil {
		return err
	}
	request := &RequestBranchCreate{}
	if err = parse(ctx, request); err != nil {
		return err
	}
	branch := &entity.Branch{ProjectID: id, Name: request.Name}
	if err := c.branch.Create(branch); err != nil {
		return presenter.CouldNotCreateProjectBranch(err)
	}
	return presenter.Success(ctx, branch)
}

// update a branch
// @Summary update a branch
// @Description update a branch
// @ID update-branch
// @Tags branch
// @Security ApiKeyAuth
//
// @Param   project_id  path     string  true  "project id"
// @Param branch_name path string                 true  "branch name"
//
// @Success 200 {object} ResponseBranch
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id}/branches/{branch_name} [PUT]
func (c *branchController) Update(ctx *fiber.Ctx) error {
	id, err := getProjectID(ctx)
	if err != nil {
		return err
	}
	branchName, err := getBranchName(ctx)
	if err != nil {
		return err
	}
	project, err := c.project.Get(id)
	if err != nil {
		return presenter.CouldNotFindProject(err)
	}
	branch, err := c.branch.UpdateBranchFromRemote(project, branchName)
	if err != nil {
		return presenter.CouldNotUpdateProject(err)
	}
	return presenter.Success(ctx, branch)
}

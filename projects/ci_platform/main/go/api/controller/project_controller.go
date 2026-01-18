package controller

import (
	"github.com/gofiber/fiber/v2"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	usecaseProject "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/project"
)

type Project interface {
	GetByID(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type (
	ResponseProject  = presenter.Response[entity.Project]
	ResponseProjects = presenter.Response[[]entity.Project]
)

type project struct {
	inter usecaseProject.Interactor
}

func NewProject(inter usecaseProject.Interactor) Project {
	return &project{inter}
}

// get a project
// @Summary get a project
// @Description get a project
// @ID get-project
// @Tags project
// @Security ApiKeyAuth
//
// @Param   project_id  path     int  true  "project id"
//
// @Success 200 {object} ResponseProject
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id} [GET]
func (c *project) GetByID(ctx *fiber.Ctx) error {
	id, err := getProjectID(ctx)
	if err != nil {
		return err
	}
	project, err := c.inter.Get(id)
	if err != nil {
		return presenter.CouldNotFindProject(err)
	}
	return presenter.Success(ctx, project)
}

// get projects
// @Summary get projects
// @Description get projects
// @ID get-projects
// @Tags project
// @Security ApiKeyAuth
//
// @Param        path    query     string  false  "search by path"
//
// @Success 200 {object} ResponseProjects
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects [GET]
func (c *project) Get(ctx *fiber.Ctx) error {
	path := ctx.Query("path")
	if path == "" {
		projects, err := c.inter.GetAll()
		if err != nil {
			return presenter.CouldNotFindProject(err)
		}
		return presenter.Success(ctx, projects)
	}
	project, err := c.inter.GetByPath(path)
	if err != nil {
		return presenter.CouldNotFindProject(err)
	}
	return presenter.Success(ctx, []entity.Project{*project})
}

type ProjectRequestCreate struct {
	Name          string `json:"name" validate:"required"`
	Source        string `json:"source" validate:"required"`
	DefaultBranch string `json:"default_branch"`
}

// create a project
// @Summary create a project
// @Description create a project
// @ID create-project
// @Tags project
// @Security ApiKeyAuth
//
// @Param Body body ProjectRequestCreate true "Project body"
//
// @Success 200 {object} ResponseProject
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects [POST]
func (c *project) Create(ctx *fiber.Ctx) error {
	request := &ProjectRequestCreate{}
	if err := parse(ctx, request); err != nil {
		return err
	}
	project := &entity.Project{
		Name:          request.Name,
		Source:        request.Source,
		DefaultBranch: request.DefaultBranch,
	}
	if err := c.inter.Create(project); err != nil {
		return presenter.CouldNotCreateProject(err)
	}
	return presenter.Success(ctx, project)
}

// update a project
// @Summary update a project
// @Description update a project
// @ID update-project
// @Tags project
// @Security ApiKeyAuth
//
// @Param   project_id  path     string  true  "project id"
// @Param Body body entity.Project true "Project body"
//
// @Success 200 {object} ResponseProject
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id} [PUT]
func (c *project) Update(ctx *fiber.Ctx) error {
	id, err := getProjectID(ctx)
	if err != nil {
		return err
	}
	project := &entity.Project{}
	if err = parse(ctx, project); err != nil {
		return err
	}
	project.ID = id
	if err = c.inter.Update(project); err != nil {
		return presenter.CouldNotUpdateProject(err)
	}
	return presenter.Success(ctx, project)
}

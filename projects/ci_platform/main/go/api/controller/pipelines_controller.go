package controller

import (
	"github.com/gofiber/fiber/v2"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	useBranch "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/branch"
	usePipeline "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/pipeline"
	useProject "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/project"
)

type Pipeline interface {
	Create(ctx *fiber.Ctx) error
}

type pipeline struct {
	branch  useBranch.Interactor
	project useProject.Interactor
	inter   usePipeline.Interactor
}

func NewPipeline(
	branch useBranch.Interactor, project useProject.Interactor,
	interactor usePipeline.Interactor,
) Pipeline {
	return &pipeline{branch, project, interactor}
}

type (
	ResponsePipeline = presenter.Response[entity.PipelineConfig]
)

// create a pipeline
// @Summary create a pipeline
// @Description create a pipeline
// @ID create-pipeline
// @Tags pipeline
// @Security ApiKeyAuth
//
// @Param   project_id  path     string  true  "project id"
// @Param   branch_name  path     string  true  "branch name"
// @Param   pipeline_name  path     string  true  "pipeline name"
//
// @Success 200 {object} map[string]any
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/projects/{project_id}/branches/{branch_name}/pipelines/{pipeline_name} [POST]
func (c *pipeline) Create(ctx *fiber.Ctx) error {
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
	pipelines, err := c.inter.SchedulePipelines(branch)
	if err != nil {
		return presenter.CouldNotSchedulePipelines(err)
	}
	return presenter.Success(ctx, pipelines)
}

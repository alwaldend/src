package useJob

import "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"

type interactor struct {
	repository Repository
}

func New(repository Repository) Interactor {
	return &interactor{repository}
}

func (interactor *interactor) Repository() Repository {
	return interactor.repository
}

func (interactor *interactor) Create(
	pipeline *entity.PipelineConfigPipeline,
) *entity.Pipeline {
	return nil
}

package usePipeline

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	usecaseJob "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/job"
)

type pipeline struct {
	repository Repository
	job        usecaseJob.Interactor
}

func New(repository Repository, job usecaseJob.Interactor) Interactor {
	return &pipeline{repository, job}
}

func (p *pipeline) Create(
	pipeline *entity.PipelineConfigPipeline,
) *entity.Pipeline {
	return nil
}

func (p *pipeline) RunJob(job *entity.Job) error {
	request, err := http.NewRequest(
		job.Config.RequestType, job.Config.URL, nil,
	)
	if err != nil {
		return err
	}
	query := request.URL.Query()
	for key, value := range job.Config.Query {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var responseMap map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		return err
	}
	return nil
}

func (p *pipeline) SchedulePipelines(branch *entity.Branch) (
	[]entity.Pipeline, error,
) {
	return nil, nil
}

func (p *pipeline) JobErrorHandler(
	job *entity.Job, err error,
) error {
	if err == nil {
		return nil
	}
	job.Log += fmt.Sprintf("an error has occured:\n%+v", err)
	job.Result = false
	job.Status = entity.JobError
	return p.job.Repository().Save(job)
}

func (p *pipeline) Repository() Repository {
	return p.repository
}

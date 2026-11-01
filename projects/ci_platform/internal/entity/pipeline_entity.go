package entity

import "github.com/google/uuid"

type JobStatus string

var (
	JobCreated JobStatus = "CREATED"
	JobRunning JobStatus = "RUNNING"
	JobError   JobStatus = "ERROR"
)

type Pipeline struct {
	CommonFields
	ProjectID uuid.UUID `gorm:"type:uuid;not null;" json:"project_id"`
	BranchID  uuid.UUID `gorm:"type:uuid;not null;" json:"branch_id"`
}

type Job struct {
	CommonFields
	PipelineID   uuid.UUID         `gorm:"type:uuid;not null;" json:"pipeline_id"`
	RedisJobID   string            `gorm:"not null;" json:"redis_job_id"`
	RedisJobName string            `gorm:"not null;" json:"redis_job_name"`
	Result       bool              `json:"result"`
	Status       JobStatus         `gorm:"not null;" json:"status"`
	Log          string            `json:"log"`
	Config       PipelineConfigJob `json:"-"`
}

type PipelineConfig struct {
	Jobs      []PipelineConfigJob      `hcl:"job,block" json:"jobs"`
	Pipelines []PipelineConfigPipeline `hcl:"pipeline,block" json:"pipelines"`
}

type PipelineConfigJob struct {
	Name        string            `hcl:"name,label" json:"name"`
	RequestType string            `hcl:"request_type,label" json:"request_type"`
	URL         string            `hcl:"url" json:"url"`
	Query       map[string]string `hcl:"query" json:"query"`
}

type PipelineConfigPipelineJob struct {
	Name     string    `hcl:"name,label" json:"name"`
	Job      string    `hcl:"job,label" json:"job"`
	Requires *[]string `hcl:"requires" json:"requires"`
}

type PipelineConfigPipeline struct {
	Name string                      `hcl:"name,label" json:"name"`
	Jobs []PipelineConfigPipelineJob `hcl:"job,block" json:"jobs"`
}

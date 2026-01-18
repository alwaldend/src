package entity

import (
	"github.com/google/uuid"
)

type Branch struct {
	CommonFields
	ProjectID uuid.UUID       `gorm:"type:uuid;not null;index:branch_name,unique;" json:"project_id"`
	Name      string          `gorm:"not null;index:branch_name,unique;size:256" json:"name" validate:"required,min=1,max=256"`
	Config    *PipelineConfig `gorm:"not null;serializer:json" json:"config"`
}

// project
type Project struct {
	CommonFields
	// foreign key for the parent namespace
	NamespaceID *uuid.UUID `gorm:"type:uuid;index:project_name,unique" json:"namespace_id"`
	// display name
	Name          string   `gorm:"index:project_name,unique;size:256" json:"name" validate:"required,min=1,max=256"`
	Path          string   `gorm:"uniqueIndex" json:"path"`
	Description   string   `json:"description"`
	Source        string   `json:"source" validate:"required"`
	DefaultBranch string   `gorm:"not null;size:256;default:master" json:"default_branch" validate:"required,min=1,max=256"`
	Branches      []Branch `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

// contains groups and projects
type Namespace struct {
	CommonFields
	// foreign key for the parent group
	ParentID *uuid.UUID `gorm:"type:uuid;index:namespace_name,unique" json:"parent_id"`
	// display name
	Name       string      `gorm:"index:namespace_name,unique;size:256" json:"name" validate:"required,min=1,max=256"`
	Path       string      `json:"path"`
	Namespaces []Namespace `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Projects   []Project   `gorm:"foreignKey:NamespaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

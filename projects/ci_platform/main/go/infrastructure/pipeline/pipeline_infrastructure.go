package infraPipeline

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	usecasePipeline "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/pipeline"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) usecasePipeline.Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]entity.Pipeline, error) {
	return r.Find()
}

func (r *repository) FindOne(condition *entity.Pipeline,
) (*entity.Pipeline, error) {
	pipeline := &entity.Pipeline{}
	return pipeline, r.db.First(&pipeline, condition).Error
}

func (r *repository) Get(id uuid.UUID) (*entity.Pipeline, error) {
	pipeline := &entity.Pipeline{}
	return pipeline, r.db.First(&pipeline, id).Error
}

func (r *repository) Find(
	conditions ...interface{},
) ([]entity.Pipeline, error) {
	pipelines := make([]entity.Pipeline, 0)
	return pipelines, r.db.Find(&pipelines, conditions...).Error
}

func (r *repository) Save(Pipeline *entity.Pipeline) error {
	return r.db.Save(Pipeline).Error
}

func (r *repository) Create(Pipeline *entity.Pipeline) error {
	return r.db.Create(Pipeline).Error
}

func (r *repository) Delete(Pipeline *entity.Pipeline) error {
	return r.db.Delete(Pipeline).Error
}

func (r *repository) Migrate() error {
	return r.db.AutoMigrate(&entity.Pipeline{})
}

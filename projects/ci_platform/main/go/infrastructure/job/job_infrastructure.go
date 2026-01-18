package infra

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	useJob "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/job"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) useJob.Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]entity.Job, error) { return r.Find() }

func (r *repository) FindOne(condition *entity.Job) (*entity.Job, error) {
	j := &entity.Job{}
	return j, r.db.First(&j, condition).Error
}

func (r *repository) Get(id uuid.UUID) (*entity.Job, error) {
	j := &entity.Job{}
	return j, r.db.First(&j, id).Error
}

func (r *repository) GetByRedisID(id string) (*entity.Job, error) {
	j := &entity.Job{}
	return j, r.db.First(&j, &entity.Job{RedisJobID: id}).Error
}

func (r *repository) Find(conditions ...interface{}) ([]entity.Job, error) {
	jobs := make([]entity.Job, 0)
	return jobs, r.db.Find(&jobs, conditions...).Error
}

func (r *repository) Save(job *entity.Job) error {
	return r.db.Save(job).Error
}

func (r *repository) Create(job *entity.Job) error {
	return r.db.Create(job).Error
}

func (r *repository) Delete(job *entity.Job) error {
	return r.db.Delete(job).Error
}

func (r *repository) Migrate() error {
	return r.db.AutoMigrate(&entity.Job{})
}

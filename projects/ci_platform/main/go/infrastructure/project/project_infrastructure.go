package infraProject

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	usecaseProject "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/project"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) usecaseProject.Repository { return &repository{db: db} }

func (r *repository) GetAll() ([]entity.Project, error) { return r.Find() }

func (r *repository) FindOne(condition *entity.Project) (
	*entity.Project, error,
) {
	p := &entity.Project{}
	return p, r.db.First(&p, condition).Error
}

func (r *repository) Get(id uuid.UUID) (*entity.Project, error) {
	p := &entity.Project{}
	return p, r.db.First(&p, id).Error
}

func (r *repository) Find(conditions ...interface{}) (
	[]entity.Project, error,
) {
	projects := make([]entity.Project, 0)
	return projects, r.db.Find(&projects, conditions...).Error
}

func (r *repository) Update(project *entity.Project) error {
	return r.db.Updates(project).Error
}

func (r *repository) Create(project *entity.Project) error {
	return r.db.Create(project).Error
}

func (r *repository) Save(project *entity.Project) error {
	return r.db.Save(project).Error
}

func (r *repository) Delete(project *entity.Project) error {
	return r.db.Delete(project).Error
}

func (r *repository) Migrate() error {
	return r.db.AutoMigrate(
		&entity.Namespace{},
		&entity.Project{},
		&entity.Branch{},
	)
}

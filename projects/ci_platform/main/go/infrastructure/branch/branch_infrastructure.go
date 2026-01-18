package infra

import (
	"io"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"gorm.io/gorm"

	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	useBranch "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/branch"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) useBranch.Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]entity.Branch, error) {
	return r.Find()
}

func (r *repository) FindOne(condition *entity.Branch) (
	*entity.Branch, error,
) {
	branch := &entity.Branch{}
	return branch, r.db.Limit(1).Find(&branch, condition).Error
}

func (r *repository) Get(id uuid.UUID) (*entity.Branch, error) {
	branch := &entity.Branch{}
	return branch, r.db.First(&branch, id).Error
}

func (r *repository) Find(conditions ...interface{}) (
	[]entity.Branch, error,
) {
	branches := make([]entity.Branch, 0)
	return branches, r.db.Find(&branches, conditions...).Error
}

func (r *repository) Update(branch *entity.Branch) error {
	return r.db.Updates(branch).Error
}

func (r *repository) Save(branch *entity.Branch) error {
	return r.db.Save(branch).Error
}

func (r *repository) Create(branch *entity.Branch) error {
	return r.db.Create(branch).Error
}

func (r *repository) Delete(branch *entity.Branch) error {
	return r.db.Delete(branch).Error
}

func (r *repository) Migrate() error {
	return r.db.AutoMigrate(&entity.Branch{})
}

func (r *repository) GetRemotePipelineConfig(
	source string, branch string, filePath string,
) (config *entity.PipelineConfig, err error) {
	filesystem := memfs.New()
	_, err = git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{
		URL:           source,
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return nil, err
	}
	fileObject, err := filesystem.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileBytes, err := io.ReadAll(fileObject)
	if err != nil {
		return nil, nil
	}
	config = &entity.PipelineConfig{}
	err = hclsimple.Decode(filePath, fileBytes, nil, config)
	return config, err
}

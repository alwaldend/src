package infra

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	useAccess "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/access"
	"gorm.io/gorm"
)

// https://github.com/casbin/casbin
// https://github.com/casdoor/casdoor
// https://casbin.org/docs/en/adapters
// https://github.com/casbin/gorm-adapter
type access struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
	confPath string
	adapter  *gormadapter.Adapter
}

// https://github.com/casbin/gorm-adapter#customize-table-columns-example
type casbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:512;uniqueIndex:unique_index"`
	V0    string `gorm:"size:512;uniqueIndex:unique_index"`
	V1    string `gorm:"size:512;uniqueIndex:unique_index"`
	V2    string `gorm:"size:512;uniqueIndex:unique_index"`
	V3    string `gorm:"size:512;uniqueIndex:unique_index"`
	V4    string `gorm:"size:512;uniqueIndex:unique_index"`
	V5    string `gorm:"size:512;uniqueIndex:unique_index"`
}

func NewService(db *gorm.DB, confPath string) useAccess.Interactor {
	return &access{db: db, confPath: confPath}
}

func (a *access) SetupEnforcer() {
	gormadapter.TurnOffAutoMigrate(a.db)
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(
		a.db, &casbinRule{},
	)
	a.adapter = adapter
	if err != nil {
		log.Panicln("could not create the adapter: ", err)
	}
	enforcer, err := casbin.NewEnforcer(a.confPath, adapter)
	if err != nil {
		log.Panicln("could not create the enforcer: ", err)
	}
	a.enforcer = enforcer
	if err = a.Load(); err != nil {
		log.Panicln("could not load the policy: ", err)
	}
}

func (a *access) Migrate() error { return a.db.AutoMigrate(&casbinRule{}) }
func (a *access) Load() error    { return a.enforcer.LoadPolicy() }
func (a *access) Enforce(subject, object, action interface{}) (bool, error) {
	return a.enforcer.Enforce(subject, object, action)
}

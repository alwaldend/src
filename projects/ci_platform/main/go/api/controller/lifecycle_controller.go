package controller

import (
	"github.com/gofiber/fiber/v2"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	useLifecycle "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/lifecycle"
)

type Lifecycle interface {
	Migrate(ctx *fiber.Ctx) error
}

type lifecycleController struct {
	interactor useLifecycle.Interactor
}

func NewLifecycle(interactor useLifecycle.Interactor) Lifecycle {
	return &lifecycleController{interactor: interactor}
}

// migrate databases
// @Summary migrate databases
// @Description migrate databases
// @ID database-migrate
// @Tags database
// @Security ApiKeyAuth
//
// @Success 200 {object} presenter.SuccessModel
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/database/migrate [POST]
func (c *lifecycleController) Migrate(ctx *fiber.Ctx) error {
	if err := c.interactor.Migrate(); err != nil {
		return presenter.CouldNotMigrateDatabase(err)
	}
	return presenter.SuccessDefault(ctx)
}

package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// serve embedded frontend files
type Frontend interface {
	Serve(ctx *fiber.Ctx) error
	ServeSwagger(ctx *fiber.Ctx) error
}

type frontendController struct {
	filesystem fiber.Handler
	swagger    fiber.Handler
}

func NewFrontend(frontend http.FileSystem, swagger http.FileSystem) Frontend {
	return &frontendController{
		filesystem: filesystem.New(filesystem.Config{
			Root:         frontend,
			PathPrefix:   "dist",
			NotFoundFile: "dist/index.html",
		}),
		swagger: filesystem.New(filesystem.Config{
			Root: swagger,
		}),
	}
}

// get the current OpenAPI schema
// @Summary get the current OpenAPI schema
// @Description get the current OpenAPI schema
// @ID get-openapi
// @Tags openapi
// @Security ApiKeyAuth
//
// @Success 200 {object} map[string]any
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/swagger/swagger.json [PUT]
func (c *frontendController) ServeSwagger(ctx *fiber.Ctx) error {
	return c.swagger(ctx)
}

func (c *frontendController) Serve(ctx *fiber.Ctx) error {
	return c.filesystem(ctx)
}

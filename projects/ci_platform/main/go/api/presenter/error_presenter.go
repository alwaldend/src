package presenter

import (
	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
	"github.com/gofiber/fiber/v2"
)

func Error(context *fiber.Ctx, err pkgError.Error) error {
	return context.Status(err.StatusCode()).JSON(
		&ResponseError{Error: err.Error()},
	)
}

package presenter

import (
	"github.com/gofiber/fiber/v2"
	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
)

func Error(context *fiber.Ctx, err pkgError.Error) error {
	return context.Status(err.StatusCode()).JSON(
		&ResponseError{Error: err.Error()},
	)
}

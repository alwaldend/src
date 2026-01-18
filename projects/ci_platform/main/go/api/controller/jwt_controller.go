package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
)

type JWT interface {
	// JSON Web Token (JWT) auth middleware.
	// For valid token, it sets the user in Ctx.Locals and calls next handler.
	// For invalid token, it returns "401 - Unauthorized" error.
	// For missing token, it returns "400 - Bad Request" error.
	// https://pkg.go.dev/github.com/gofiber/jwt/v3@v3.2.14#section-readme
	CheckJWT(ctx *fiber.Ctx) error
}

type jwtController struct {
	handler fiber.Handler
}

func NewJWT(secret []byte) JWT {
	config := jwtware.Config{
		SigningKey:  secret,
		ContextKey:  "user",
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ErrorHandler: func(context *fiber.Ctx, err error) error {
			return pkgError.NewWithCode(
				fmt.Errorf("JWT error: %w", err), http.StatusUnauthorized,
			)
		},
	}
	return &jwtController{handler: jwtware.New(config)}
}

func (c *jwtController) CheckJWT(ctx *fiber.Ctx) error {
	return c.handler(ctx)
}

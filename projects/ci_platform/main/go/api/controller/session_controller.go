package controller

import (
	"github.com/gofiber/fiber/v2"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	useJWT "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/jwt"
	useSession "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/session"
)

var UnauthorizedRoutes = map[string]bool{
	"/api/v1/users":       true,
	"/api/v1/users/login": true,
}

type Session interface {
	CheckSession(ctx *fiber.Ctx) error
}

type session struct {
	interactor    useSession.Interactor
	jwtController JWT
	jwtInteractor useJWT.Interactor
}

func NewSession(
	inter useSession.Interactor, jwtCtrl JWT, jwt useJWT.Interactor,
) Session {
	return &session{
		interactor:    inter,
		jwtController: jwtCtrl,
		jwtInteractor: jwt,
	}
}

func (c *session) CheckSession(ctx *fiber.Ctx) error {
	// allow registration and login without a token
	if _, defined := UnauthorizedRoutes[ctx.Path()]; defined {
		return nil
	}
	session, err := c.interactor.Get(ctx)
	if err != nil {
		return presenter.CouldNotGetSession(err)
	}
	// if current session contains user_id, the user is logged in
	if session.Get(useSession.UserID) != nil {
		return nil
	}
	// check if there is a JWT in the header
	if err := c.jwtController.CheckJWT(ctx); err != nil {
		return err
	}
	claims, err := c.jwtInteractor.Claims(ctx.Locals("user"))
	if err != nil {
		return err
	}
	session.Set(useSession.UserID, claims.ID)
	return nil
}

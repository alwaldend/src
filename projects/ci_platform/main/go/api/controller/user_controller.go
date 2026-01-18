package controller

import (
	"github.com/gofiber/fiber/v2"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/entity"
	useSession "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/session"
	useUser "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/user"
)

type User interface {
	GetAll(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type user struct {
	inter   useUser.Interactor
	session useSession.Interactor
}

// NewUser: constructor, dependency injection from user service and firebase service
func NewUser(
	inter useUser.Interactor, session useSession.Interactor,
) User {
	return &user{inter, session}
}

type (
	ResponseUser  = presenter.Response[entity.User]
	ResponseUsers = presenter.Response[entity.User]
)

// Get all users
// @Summary get all users
// @Description get all users
// @ID get-users
// @Tags user
// @Security ApiKeyAuth
//
// @Success 200 {object} ResponseUsers
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/users [GET]
func (c *user) GetAll(ctx *fiber.Ctx) error {
	users, err := c.inter.GetAll()
	if err != nil {
		return presenter.CouldNotFindUser(err)
	}
	return presenter.Success(ctx, users)
}

type UserRequestCreate struct {
	Username string `json:"username" form:"username"  validate:"required,min=3,max=64"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password"  validate:"required,min=8,max=64"`
}

// Create user
// @Summary create a user
// @Description create a user
// @ID create-user
// @Tags user
//
// @Param Body body UserRequestCreate true "User body"
//
// @Success 200 {object} ResponseUser
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/users [POST]
func (c *user) Create(ctx *fiber.Ctx) error {
	request := &UserRequestCreate{}
	if err := parse(ctx, request); err != nil {
		return err
	}
	user := &entity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
	if err := c.inter.Create(user); err != nil {
		return presenter.CouldNotCreateUser(err)
	}
	return presenter.Success(ctx, user)
}

// get a user
// @Summary get a user
// @Description get a user
// @ID get-user
// @Tags user
// @Security ApiKeyAuth
//
// @Param   user_id  path     string  true  "account id"
//
// @Success 200 {object} ResponseUser
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/users/user/{user_id} [GET]
func (c *user) Get(ctx *fiber.Ctx) error {
	id, err := getUserID(ctx)
	if err != nil {
		return err
	}
	user, err := c.inter.Get(id)
	if err != nil {
		return presenter.CouldNotFindUser(err)
	}
	return presenter.Success(ctx, user)
}

type UserRequestLogin struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// login
// @Summary login
// @Description login
// @ID login-user
// @Tags user
//
// @Success 200 {object} ResponseUser
// @Failure 500 {object} presenter.ResponseError
// @Router /api/v1/users/login [GET]
func (c *user) Login(ctx *fiber.Ctx) error {
	request := &UserRequestLogin{}
	if err := parse(ctx, request); err != nil {
		return err
	}
	user, err := c.inter.FindByUsername(request.Username)
	if err != nil {
		return presenter.InvalidLoginOrPassword()
	}
	matches, err := c.inter.PasswordHashCheck(
		user,
		request.Password,
	)
	if err != nil {
		return presenter.CouldNotParsePasswordHash(err)
	}
	if !matches {
		return presenter.InvalidLoginOrPassword()
	}
	session, err := c.session.Get(ctx)
	if err != nil {
		return presenter.CouldNotGetSession(err)
	}
	ctx.Cookie(&fiber.Cookie{
		Name:  useSession.SessionCookie,
		Value: session.ID(),
	})
	session.Set(useSession.UserID, user.ID.String())
	if err := session.Save(); err != nil {
		return presenter.CouldNotSaveSession(err)
	}
	return presenter.Success(ctx, user)
}

package presenter

import (
	"errors"
	"fmt"
	"net/http"

	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
)

func InvalidLoginOrPassword() error {
	return pkgError.NewWithCode(
		errors.New("invalid login or password"), http.StatusForbidden,
	)
}

func CouldNotParsePasswordHash(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not parse password hash: %w", err),
		http.StatusInternalServerError,
	)
}

func InvalidUserID(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("invalid user id: %w", err), http.StatusBadRequest,
	)
}

func CouldNotFindUser(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not find user(s): %w", err),
		http.StatusNotFound,
	)
}

func CouldNotUpdateUser(err error) error {
	return pkgError.New(fmt.Errorf("could not update user: %w", err))
}

func CouldNotCreateUser(err error) error {
	return pkgError.New(fmt.Errorf("could not create user: %w", err))
}

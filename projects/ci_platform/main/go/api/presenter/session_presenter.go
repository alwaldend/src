package presenter

import (
	"fmt"
	"net/http"

	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
)

func CouldNotGetSession(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not get session: %w", err),
		http.StatusInternalServerError,
	)
}

func CouldNotSaveSession(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not save session: %w", err),
		http.StatusInternalServerError,
	)
}

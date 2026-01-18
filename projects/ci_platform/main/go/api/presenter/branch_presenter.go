package presenter

import (
	"fmt"
	"net/http"

	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
)

func InvalidBranchName(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("invalid branch name: %w", err), http.StatusBadRequest,
	)
}

func CouldNotFindProjectBranch(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not find project branch(es): %w", err),
		http.StatusNotFound,
	)
}

func CouldNotCreateProjectBranch(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not find create branch: %w", err),
		http.StatusNotFound,
	)
}

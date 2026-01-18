package presenter

import (
	"fmt"
	"net/http"

	pkgError "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/pkg/error"
)

func CouldNotMigrateDatabase(err error) error {
	return pkgError.NewWithCode(
		fmt.Errorf("could not migrate database: %w", err),
		http.StatusInternalServerError,
	)
}

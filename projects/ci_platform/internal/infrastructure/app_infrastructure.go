package infra

import (
	infraConfig "git.alwaldend.com/alwaldend/src/projects/ci_platform/internal/infrastructure/config"
	infraLifecycle "git.alwaldend.com/alwaldend/src/projects/ci_platform/internal/infrastructure/lifecycle"
	useLifecycle "git.alwaldend.com/alwaldend/src/projects/ci_platform/internal/usecase/lifecycle"
)

func NewApp(configs ...*infraConfig.App) useLifecycle.Interactor {
	return useLifecycle.New(infraLifecycle.New(configs...).Setup())
}

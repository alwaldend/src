package infra

import (
	infraConfig "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure/config"
	infraLifecycle "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure/lifecycle"
	useLifecycle "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/lifecycle"
)

func NewApp(configs ...*infraConfig.App) useLifecycle.Interactor {
	return useLifecycle.New(infraLifecycle.New(configs...).Setup())
}

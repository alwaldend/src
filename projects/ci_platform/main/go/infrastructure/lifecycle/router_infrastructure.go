package infraLifecycle

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/controller"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/middleware"
	"git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/api/presenter"
	_ "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/docs"
	infraConfig "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure/config"
)

func newRouter(config *fiber.Config) *fiber.App {
	config.ErrorHandler = controller.Error
	return fiber.New(*config)
}

func setupRouter(
	router *fiber.App,
	app *controller.App,
	config *infraConfig.App,
) {
	// use global middlewares
	// https://docs.gofiber.io/api/middleware
	router.Use(
		// limiter.New(),
		middleware.NewCORS(),
		compress.New(),
		etag.New(),
		// csrf.New(csrf.Config{CookieSameSite: "Strict"}),
		logger.New(),
		recover.New(recover.Config{EnableStackTrace: true}),
	)

	// groups
	// /api
	apiRoot := router.Group("api",
		middleware.NewGroup(app.Session.CheckSession),
	)
	// /api/v1
	api := apiRoot.Group("v1")
	// /api/v1/users
	users := api.Group("users")
	// /api/v1/projects
	projects := api.Group("projects")
	// /api/v1/database
	database := api.Group("database")
	// /api/v1/projects/:project_id
	project := projects.Group(":" + controller.PathProjectID)
	// /api/v1/projects/:project_id/branches
	branches := project.Group("branches")
	// /api/v1/projects/:project_id/branches/:branch_name
	branch := branches.Group(":" + controller.PathBranchName)
	// /api/v1/projects/:project_id/branches/:branch_name/pipelines
	pipelines := branch.Group("pipelines")
	// /api/v1/projects/:project_id/branches/:branch_name/pipelines/:pipeline_name
	pipeline := pipelines.Group(":" + controller.PathPipelineName)

	// user routes
	users.Get("", app.User.GetAll)
	users.Post("", app.User.Create)
	users.Post("login", app.User.Login)
	users.Get("user/:"+controller.PathUserID, app.User.Get)

	// project routes
	projects.Get("", app.Project.Get)
	projects.Post("", app.Project.Create)

	project.Get("", app.Project.GetByID)
	project.Put("", app.Project.Update)

	branches.Post("", app.Branch.Create)
	branches.Get("", app.Branch.GetProjectBranches)

	branch.Get("", app.Branch.Get)
	branch.Put("", app.Branch.Update)

	pipeline.Post("", app.Pipeline.Create)

	// swagger routes
	api.Use("swagger", app.Frontend.ServeSwagger)

	// database routes
	database.Post("migrate", app.Lifecycle.Migrate)

	// api catch all in order to not send html in response to
	// invalid api requests
	apiRoot.Use(presenter.RouteDoesNotExist)

	// frontend controller
	//
	// returns embedded files from ../infrastructure/ui/web2/dist
	// (the folder is created and populated by 'yarn build')
	//
	// it is a Single Page Application, routing is done by Javascript, that's
	// why it has to return index.html on missing files
	router.Use(app.Frontend.Serve)
}

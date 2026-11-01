package useLifecycle

// general lifecycle actions
type Interactor interface {
	// setup everything
	Setup() Interactor
	// run the server
	Run()
	// migrate everything
	Migrate() error
}

// specific server actions
type App interface {
	// listen for requests, blocks the goroutine
	Listen()
	// gracefully shutdown
	Shutdown(exitCode int)
	// gracefully shutdown on panic
	// if there is no panic, do nothing
	ShutdownOnPanic()
	// gracefully shutdown when an interruption signal is sent
	// blocks the goroutine untill the signal is recieved
	// listens for os.Interrupt and syscall.SIGTERM
	ShutdownOnInterruptionSignal()
	// automigrate all repositories
	Migrate() error

	// setup everything
	Setup() App
	// populate the config file
	SetupConfig() App
	// connect to the database
	SetupDatabase() App
	// setup swagger
	SetupSwagger() App
	// setup router
	SetupRouter() App
	// setup infrastructure
	SetupInfrastructure() App
	// setup interactors
	SetupInteractor() App
	// setup controllers
	SetupController() App
}

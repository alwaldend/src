package useLifecycle

type lifecycle struct {
	app App
}

func New(app App) Interactor {
	return &lifecycle{app: app}
}

func (l *lifecycle) Run() {
	go l.app.Listen()
	l.app.ShutdownOnInterruptionSignal()
}

func (l *lifecycle) Setup() Interactor {
	l.app.Setup()
	return l
}

func (l *lifecycle) Migrate() error {
	return l.app.Migrate()
}

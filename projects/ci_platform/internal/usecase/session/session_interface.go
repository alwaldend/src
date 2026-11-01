package useSession

const (
	UserID        = "user_id"
	SessionCookie = "session_id"
)

type Interactor interface {
	Get(context interface{}) (Session, error)
	Reset() error
	Close() error
}

type Repository interface {
	Get(context interface{}) (Session, error)
	Reset() error
	Close() error
}

type Session interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Delete(key string)
	Destroy() error
	Regenerate() error
	Save() error
	Fresh() bool
	ID() string
	Keys() []string
}

// https://docs.gofiber.io/api/middleware/session
package infraSession

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	useSession "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/session"
)

type repository struct {
	session *session.Store
}

func New(config *redis.Config) useSession.Repository {
	return &repository{
		session: session.New(session.Config{
			Storage: redis.New(*config),
		}),
	}
}

func (r *repository) Get(context interface{}) (
	useSession.Session, error,
) {
	contextAsserted, ok := context.(*fiber.Ctx)
	if !ok {
		return nil, fmt.Errorf("invalid context type: %v+", context)
	}
	return r.session.Get(contextAsserted)
}

func (r *repository) Reset() error { return r.session.Reset() }

func (r *repository) Close() error { return r.session.Storage.Close() }

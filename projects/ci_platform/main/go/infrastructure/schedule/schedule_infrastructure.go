package infraSchedule

import (
	"context"

	"github.com/go-redis/redis/v8"
	useSchedule "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/schedule"
)

type jobTypes map[string]useSchedule.Job

type schedule struct {
	rdb      *redis.Client
	jobTypes jobTypes
}

func New(rdb *redis.Client) useSchedule.Interactor {
	return &schedule{rdb, make(jobTypes, 0)}
}

func (s *schedule) Schedule(key string, values ...any) error {
	return s.rdb.LPush(context.TODO(), key, values...).Err()
}

func (s *schedule) Consume(key string) (string, error) {
	return s.rdb.LPop(context.TODO(), key).Result()
}

func (s *schedule) Job(name string, j useSchedule.Job) useSchedule.Interactor {
	s.jobTypes[name] = j
	return s
}

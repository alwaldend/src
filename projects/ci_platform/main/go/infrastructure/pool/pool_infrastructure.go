package infraPool

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	infraConfig "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/infrastructure/config"
	useJob "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/job"
	usePool "git.alwaldend.com/alwaldend/src/projects/ci_platform/main/go/usecase/pool"
)

type pool struct {
	workerPool *work.WorkerPool
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	jobRepo    useJob.Repository
}

type workJob func(*work.Job) error

func New(
	config *infraConfig.App, jobRepository useJob.Repository,
) usePool.Pool {
	redisPool := &redis.Pool{
		MaxActive: 15,
		MaxIdle:   15,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
			)
		},
	}
	namespace := "{ci_platform}"
	workerPool := work.NewWorkerPool(Context{}, 10, namespace, redisPool)
	enqueuer := work.NewEnqueuer(namespace, redisPool)
	return &pool{workerPool, redisPool, enqueuer, jobRepository}
}

func (p *pool) Start() { p.workerPool.Start() }

func (p *pool) Close() error {
	p.workerPool.Stop()
	return p.redisPool.Close()
}

func (p *pool) Job(
	name string, handler usePool.Handler, errorHandler usePool.ErrorHandler,
) {
	p.workerPool.JobWithOptions(
		name, work.JobOptions{MaxFails: 1},
		p.normalizeHandler(handler, errorHandler),
	)
}

func (p *pool) normalizeHandler(
	handler usePool.Handler, errorHandler usePool.ErrorHandler,
) workJob {
	return func(jobRedis *work.Job) error {
		job, err := p.jobRepo.GetByRedisID(jobRedis.ID)
		if err != nil {
			return fmt.Errorf(
				"could not get the job from the database: %w", err,
			)
		}
		if err = handler(job); err != nil {
			return errorHandler(job, err)
		}
		return nil
	}
}

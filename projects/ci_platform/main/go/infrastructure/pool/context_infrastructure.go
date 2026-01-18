package infraPool

import (
	"fmt"

	"github.com/gocraft/work"
)

type Context struct{}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

package fp

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
)

type WaitGroupE struct {
	wg    sync.WaitGroup
	err   error
	errMx sync.Mutex
}

func (self *WaitGroupE) Go(fs ...func() error) {
	for _, f := range fs {
		self.wg.Go(func() {
			var err error
			defer func() {
				if r := recover(); r != nil {
					switch rTyped := r.(type) {
					case error:
						err = errors.Join(err, fmt.Errorf("panic: %s\n%s: %w", rTyped, string(debug.Stack()), rTyped))
					default:
						err = errors.Join(err, fmt.Errorf("panic: %s\n%s", r, string(debug.Stack())))
					}
				}
				self.errMx.Lock()
				defer self.errMx.Unlock()
				self.err = errors.Join(self.err, err)
			}()
			err = f()
		})
	}
}

func (self *WaitGroupE) Wait() error {
	self.wg.Wait()
	return self.err
}

func (self *WaitGroupE) WaitCtx(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- self.Wait()
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timed out")
	}
}

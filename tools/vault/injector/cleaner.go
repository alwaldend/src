package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Cleaner struct {
	consume chan string
	mx      *sync.RWMutex
	logger  *log.Logger
	closed  bool
}

func NewCleaner(stderr io.Writer) *Cleaner {
	return &Cleaner{
		consume: make(chan string),
		mx:      &sync.RWMutex{},
		logger:  log.New(stderr, "com.alwaldend.src.tools.vault.injector.cleaner ", log.Flags()),
	}
}

func (self *Cleaner) Add(paths ...string) error {
	self.mx.RLock()
	defer self.mx.RUnlock()
	if self.closed {
		return fmt.Errorf("could not add: closed")
	}
	for _, p := range paths {
		self.consume <- p
	}
	return nil
}

func (self *Cleaner) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		self.logger.Printf("Shutting down cleaner\n")
		self.mx.Lock()
		defer self.mx.Unlock()
		close(self.consume)
		self.closed = true
	})
	wg.Go(func() {
		toClean := []string{}
		for {
			path, ok := <-self.consume
			if !ok {
				break
			}
			toClean = append(toClean, path)
		}
		for _, c := range toClean {
			wg.Go(func() {
				self.logger.Printf("Cleaning path %s\n", c)
				if err := os.RemoveAll(c); err != nil {
					self.logger.Printf("could not clean path %s: could not remove: %s", c, err)
				}
			})
		}
	})
	wg.Wait()
}

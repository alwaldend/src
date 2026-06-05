package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al_plugin"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
)

type plugin struct {
	cleanConsume chan string
	cleanFiles   []string
	cleanStart   chan struct{}
	ctx          context.Context
	v            *al.VaultStore
	wg           *sync.WaitGroup
	return_code  int
	logger       *log.Logger
}

func (self plugin) PluginStart(ctx context.Context, req *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (self *plugin) cleanRun() {
	<-self.cleanStart
	for _, c := range self.cleanFiles {
		self.logger.Printf("cleaning path %s", c)
		if err := os.RemoveAll(c); err != nil {
			self.logger.Printf("could not clean path %s: could not remove: %s", c, err)
		}
	}
}

func (self *plugin) cleanConsumeRun() {
	for {
		select {
		case c := <-self.cleanConsume:
			self.cleanFiles = append(self.cleanFiles, c)
		case <-self.ctx.Done():
			close(self.cleanConsume)
			close(self.cleanStart)
			return
		}
	}
}

func (self *plugin) run() int {
	self.logger.Printf("starting plugin")
	self.wg.Go(self.cleanRun)
	self.wg.Go(self.cleanConsumeRun)
	self.wg.Go(func() {
		if _, err := al_plugin.Serve(self.ctx, os.Stdin, os.Stdout, self).Get(); err != nil {
			self.logger.Printf("could not serve the plugin: %s", err)
			self.return_code = 1
		}
	})
	self.wg.Wait()
	return self.return_code
}

func main() {
	p := &plugin{
		cleanConsume: make(chan string),
		cleanStart:   make(chan struct{}, 1),
		wg:           &sync.WaitGroup{},
		logger:       log.New(os.Stderr, "com.alwaldend.src.tools.vault.injector ", log.Flags()),
	}
	p.ctx, _ = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	os.Exit(p.run())
}

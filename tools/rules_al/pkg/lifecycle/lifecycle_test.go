package lifecycle

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"
)

type testResource struct {
	start func(context.Context) error
	stop  func(context.Context) error
}

func (r testResource) Start(c context.Context) error { return r.start(c) }
func (r testResource) Stop(c context.Context) error  { return r.stop(c) }
func TestStartFailureRollsBackNestedPartialResources(t *testing.T) {
	var stopped []int
	failure := errors.New("start failure")
	var inner, outer Manager
	for i := range 3 {
		inner.Add(testResource{start: func(context.Context) error {
			if i == 1 {
				return failure
			}
			return nil
		}, stop: func(c context.Context) error {
			if c.Err() != nil {
				t.Error("rollback used canceled startup context")
			}
			stopped = append(stopped, i)
			return nil
		}})
	}
	outer.Add(&inner)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := outer.Start(ctx); !errors.Is(err, failure) {
		t.Fatalf("start error %v", err)
	}
	if !reflect.DeepEqual(stopped, []int{2, 1, 0}) {
		t.Fatalf("cleanup order %v", stopped)
	}
	if err := outer.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(stopped) != 3 {
		t.Fatal("cleanup repeated")
	}
	if err := inner.Add(StoppableFunc0(func() error { return nil })); err == nil {
		t.Fatal("accepted registration after shutdown")
	}
}

func TestStopContinuesAfterPanicAndError(t *testing.T) {
	var m Manager
	called := false
	expected := errors.New("cleanup failed")
	m.AddState(StateStarted, StoppableFunc0(func() error { called = true; return expected }), StoppableFunc0(func() error { panic("cleanup panic") }))
	if err := m.Stop(context.Background()); !errors.Is(err, expected) {
		t.Fatalf("error %v", err)
	}
	if !called {
		t.Fatal("cleanup skipped after panic")
	}
}

func TestProcessHelper(t *testing.T) {
	mode := os.Getenv("AL_PROCESS_TEST")
	if mode == "" {
		return
	}
	signals := make(chan os.Signal, 1)
	if mode == "default" {
		signal.Reset(syscall.SIGTERM)
	} else if mode == "ignore" {
		signal.Ignore(syscall.SIGTERM)
	} else {
		signal.Notify(signals, syscall.SIGTERM)
	}
	os.Stdout.Write([]byte("ready"))
	if mode == "ignore" || mode == "default" {
		for {
			time.Sleep(time.Hour)
		}
	}
	<-signals
	if mode == "fail" {
		os.Exit(23)
	}
	os.Exit(0)
}

func TestStopProcessWaitsAndReaps(t *testing.T) {
	for _, mode := range []string{"clean", "fail", "ignore"} {
		t.Run(mode, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestProcessHelper$")
			cmd.Env = append(os.Environ(), "AL_PROCESS_TEST="+mode)
			out, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatal(err)
			}
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			ready := make([]byte, 5)
			if _, err := io.ReadFull(out, ready); err != nil {
				t.Fatal(err)
			}
			wait := make(chan error, 1)
			go func() { wait <- cmd.Wait() }()
			// Race-instrumented helpers delay normal process exit by one second.
			grace := 3 * time.Second
			if mode == "ignore" {
				grace = 100 * time.Millisecond
			}
			ctx, cancel := context.WithTimeout(context.Background(), grace)
			defer cancel()
			err = StopProcess(ctx, cmd, wait)
			if (err != nil) != (mode != "clean") {
				t.Fatalf("stop error %v", err)
			}
			if cmd.ProcessState == nil {
				t.Fatal("child not reaped")
			}
		})
	}
}

func TestStartFailureCancelsSibling(t *testing.T) {
	var m Manager
	failure := errors.New("failure")
	entered := make(chan struct{})
	m.Add(testResource{start: func(ctx context.Context) error { close(entered); <-ctx.Done(); return ctx.Err() }, stop: func(context.Context) error { return nil }})
	m.Add(testResource{start: func(context.Context) error { <-entered; return failure }, stop: func(context.Context) error { return nil }})
	done := make(chan error, 1)
	go func() { done <- m.Start(context.Background()) }()
	select {
	case err := <-done:
		if !errors.Is(err, failure) {
			t.Fatalf("start error %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("failed start did not cancel sibling")
	}
}

func TestLifecyclePanicDoesNotExposeValue(t *testing.T) {
	err := rescueLifecycleFunc(func(context.Context) error { panic("synthetic-secret-value") })(context.Background())
	if err == nil || strings.Contains(err.Error(), "synthetic-secret-value") {
		t.Fatalf("panic was not redacted: %v", err)
	}
}

func TestStopBeforeStartReleasesConstructedResources(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	var inner, outer Manager
	inner.Add(StoppableFunc0(reader.Close), StoppableFunc0(writer.Close))
	outer.Add(&inner)
	if err := outer.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("secret")); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("writer remains open: %v", err)
	}
	if _, err := reader.Read(make([]byte, 1)); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("reader remains open: %v", err)
	}
	if err := outer.Stop(context.Background()); err != nil {
		t.Fatalf("cleanup repeated: %v", err)
	}
}

func TestStopResourceProcess(t *testing.T) {
	for _, mode := range []string{"default", "clean", "fail", "ignore"} {
		t.Run(mode, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestProcessHelper$")
			cmd.Env = append(os.Environ(), "AL_PROCESS_TEST="+mode)
			out, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatal(err)
			}
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			if _, err := io.ReadFull(out, make([]byte, 5)); err != nil {
				t.Fatal(err)
			}
			wait := make(chan error, 1)
			go func() { wait <- cmd.Wait() }()
			// Race-instrumented helpers delay normal process exit by one second.
			grace := 3 * time.Second
			if mode == "ignore" {
				grace = 100 * time.Millisecond
			}
			ctx, cancel := context.WithTimeout(context.Background(), grace)
			defer cancel()
			err = StopResourceProcess(ctx, cmd, wait)
			if (err != nil) != (mode == "fail" || mode == "ignore") {
				t.Fatalf("shutdown error %v", err)
			}
			if cmd.ProcessState == nil {
				t.Fatal("child not reaped")
			}
		})
	}
}

func TestStopResourceProcessPreservesPriorExit(t *testing.T) {
	expected := errors.New("prior process failure")
	wait := make(chan error, 1)
	wait <- expected
	if err := StopResourceProcess(context.Background(), &exec.Cmd{Process: &os.Process{}}, wait); !errors.Is(err, expected) {
		t.Fatalf("prior error lost: %v", err)
	}
}

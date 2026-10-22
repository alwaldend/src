package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/cli"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/pkg/ankiascode"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"github.com/spf13/cobra"
)

// run composes optional server adapters and same-executable workers.
func run(ctx context.Context, arguments []string) error {
	runServer := nativeServerRunner
	if runServer == nil {
		runServer = wasmServerRunner
	}
	var binary, defaultAddress string
	if runServer != nil {
		var err error
		binary, err = os.Executable()
		if err != nil {
			return fmt.Errorf("resolve application executable: %w", err)
		}
		defaultAddress = ankigrpc.SubprocessAddress(binary, "server", "serve")
	}
	engines := map[string]func([]string) error{"cgo": nativeServerRunner, "wasm": wasmServerRunner}
	var command *cobra.Command
	hostedRunner := func(arguments []string) error {
		selectedEngine, err := command.PersistentFlags().GetString("server-engine")
		if err != nil {
			return fmt.Errorf("read server engine selector: %w", err)
		}
		runner := runServer
		if selectedEngine != "" {
			runner = engines[selectedEngine]
			if runner == nil {
				return fmt.Errorf("run Anki server: engine %q is unavailable", selectedEngine)
			}
		}
		if runner == nil {
			return fmt.Errorf("run Anki server: no local engines are compiled")
		}
		// The server loop does not observe main's Go cancellation context, so
		// restore SIGINT termination instead of only cancelling that context.
		signal.Reset(os.Interrupt)
		// Native writes to a disconnected client can raise SIGPIPE; let the
		// server handle the write error instead of terminating the worker.
		signal.Ignore(syscall.SIGPIPE)
		if err := runner(arguments); err != nil {
			return fmt.Errorf("host configured Anki server: %w", err)
		}
		return nil
	}
	workerFactory := ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](ankigrpc.NewSubprocessClient)
	factory := ankigrpc.NewFactory(defaultAddress)
	factory.RegisterClientFactory("subprocess", workerFactory)
	for engine, runner := range engines {
		if runner != nil {
			factory.RegisterClientFactory("subprocess+"+engine, engineWorkerFactory(engine, binary))
		}
	}
	if wasmClientFactory != nil {
		factory.RegisterClientFactory("inprocess+wasm", wasmClientFactory)
		if nativeServerRunner == nil {
			factory.RegisterClientFactory("inprocess", wasmClientFactory)
		}
	}
	application := ankiascode.New(ankiascode.Options{
		NewClient:    factory,
		ServerRunner: hostedRunner,
	})
	command = cli.New(application)
	command.SetArgs(arguments)
	if err := command.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

// engineWorkerFactory selects a compiled engine for same-executable workers.
// Independently supplied executables retain their exact argument vectors.
func engineWorkerFactory(engine, binary string) ankigrpc.ClientFactory {
	return ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](func(address string) (*ankiclient.ClientSTD, error) {
		options, err := ankigrpc.ParseSubprocessAddress(address)
		if err != nil {
			return nil, fmt.Errorf("parse %s worker address: %w", engine, err)
		}
		if sameExecutable(binary, options.Binary) {
			for _, argument := range options.Arguments {
				if argument == "--server-engine" || strings.HasPrefix(argument, "--server-engine=") {
					return nil, fmt.Errorf("qualified subprocess scheme owns engine selection; omit --server-engine from worker arguments")
				}
			}
			options.Arguments = append([]string{"--server-engine=" + engine}, options.Arguments...)
		}
		return ankiclient.NewSTD(options), nil
	})
}

// sameExecutable recognizes direct paths, PATH lookup and filesystem aliases.
// Unresolvable workers retain their arguments and fail during normal startup.
func sameExecutable(binary, worker string) bool {
	if binary == worker {
		return true
	}
	path, err := exec.LookPath(worker)
	if err != nil {
		return false
	}
	application, err := os.Stat(binary)
	if err != nil {
		return false
	}
	executable, err := os.Stat(path)
	return err == nil && os.SameFile(application, executable)
}

// main owns executable cancellation, command invocation, and failure reporting.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// newAPICommand assembles the complete service tree from compiled upstream descriptors.
func newAPICommand(application ankiascodeapi.App) *cobra.Command {
	command := &cobra.Command{Use: "api", Short: "Invoke upstream Anki services using protobuf JSON or TOML", Args: cobra.NoArgs}
	services, err := rpc.Services()
	if err != nil {
		command.RunE = func(*cobra.Command, []string) error { return fmt.Errorf("assemble Anki API commands: %w", err) }
		return command
	}
	for _, service := range services {
		group := &cobra.Command{Use: service.Name, Short: string(service.FullName), Args: cobra.NoArgs}
		for _, method := range service.Methods {
			group.AddCommand(newAPIMethodCommand(application, method))
		}
		command.AddCommand(group)
	}
	return command
}

// methodHelp describes the selected schemas without relying on stripped source comments.
func methodHelp(method rpc.Method) string {
	var help strings.Builder
	fmt.Fprintf(&help, "%s\n\nRequest: %s\nResponse: %s\n", method.Path, method.Descriptor.Input().FullName(), method.Descriptor.Output().FullName())
	fields := method.Descriptor.Input().Fields()
	if fields.Len() != 0 {
		help.WriteString("\nRequest fields:\n")
	}
	for index := 0; index < fields.Len(); index++ {
		field := fields.Get(index)
		kind := field.Kind().String()
		if field.Kind() == protoreflect.MessageKind {
			kind = string(field.Message().FullName())
		} else if field.Kind() == protoreflect.EnumKind {
			kind = string(field.Enum().FullName())
		}
		fmt.Fprintf(&help, "  %s: %s %s\n", field.Name(), field.Cardinality(), kind)
	}
	return help.String()
}

// readAPIRequest bounds input acquisition without closing caller-owned streams.
// A canceled arbitrary reader can finish when its owner closes it; process exit releases CLI stdin.
func readAPIRequest(ctx context.Context, command *cobra.Command, input string, request protoreflect.MessageDescriptor) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("start API input: %w", errors.Join(err, context.Cause(ctx)))
	}
	reader := command.InOrStdin()
	skipTerminal := request.Fields().Len() == 0 && !command.Flags().Changed("input")
	result := make(chan []byte, 1)
	failure := make(chan error, 1)
	go func() {
		data, err := readAPIInput(reader, input, skipTerminal)
		if err != nil {
			failure <- err
			return
		}
		result <- data
	}()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for API request input: %w", errors.Join(ctx.Err(), context.Cause(ctx)))
	case err := <-failure:
		return nil, fmt.Errorf("read API input: %w", err)
	case data := <-result:
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("finish API input: %w", errors.Join(err, context.Cause(ctx)))
		}
		return data, nil
	}
}

// readAPIInput reads explicit files or piped stdin, avoiding a terminal read for empty requests.
func readAPIInput(reader io.Reader, input string, skipTerminal bool) ([]byte, error) {
	if input != "" && input != "-" {
		data, err := os.ReadFile(input)
		if err != nil {
			return nil, fmt.Errorf("read API request file: %w", err)
		}
		return data, nil
	}
	if skipTerminal {
		if file, ok := reader.(*os.File); ok {
			info, err := file.Stat()
			if err != nil {
				return nil, fmt.Errorf("inspect API request input: %w", err)
			}
			if info.Mode()&os.ModeCharDevice != 0 {
				return nil, nil
			}
		}
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read API request stdin: %w", err)
	}
	return data, nil
}

// newAPIMethodCommand shares payload handling and connection execution across every RPC.
func newAPIMethodCommand(application ankiascodeapi.App, method rpc.Method) *cobra.Command {
	var source, address, input, inputFormat, outputFormat string
	var timeout time.Duration
	writeBack := true
	command := &cobra.Command{
		Use: method.Name, Short: method.Path, Long: methodHelp(method), Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if (source == "") == (address == "") {
				return fmt.Errorf("API calls require exactly one of --collection and --address")
			}
			ctx, cancel := context.WithTimeout(command.Context(), timeout)
			defer cancel()
			data, err := readAPIRequest(ctx, command, input, method.Descriptor.Input())
			if err != nil {
				return fmt.Errorf("read API request: %w", err)
			}
			response, err := application.APICall(ctx, ankiascodeapi.APIRequest{
				Collection: ankiascodeapi.CollectionOptions{Path: source, Address: address}, WriteBack: writeBack, Method: method.Path,
				Input: data, InputFormat: inputFormat, OutputFormat: outputFormat,
			})
			if err != nil {
				return fmt.Errorf("execute Anki API command: %w", err)
			}
			if _, err := command.OutOrStdout().Write(response); err != nil {
				return fmt.Errorf("write API response: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&source, "collection", "", "Offline .anki2 database or .colpkg archive")
	addressFlag.bind(command, &address)
	command.Flags().StringVar(&input, "input", "-", "Request file, or - for stdin")
	command.Flags().StringVar(&inputFormat, "input-format", "json", "Request encoding: json or toml")
	command.Flags().StringVar(&outputFormat, "output-format", "json", "Response encoding: json or toml")
	command.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Input, startup and request timeout")
	command.Flags().BoolVar(&writeBack, "write-back", true, "Write offline archives back on close")
	command.Flags().BoolFunc("no-write-back", "Disable offline archive write-back", func(value string) error {
		disabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("parse no-write-back: %w", err)
		}
		writeBack = !disabled
		return nil
	})
	return command
}

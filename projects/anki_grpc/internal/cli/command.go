package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// options holds shared connection flags and the upstream message resolver.
type options struct {
	// worker selects a managed stdio executable, resolved from runfiles by default.
	worker string
	// target selects an existing TCP or Unix gRPC server instead of a process.
	target string
	// collection opens a local collection for one managed invocation.
	collection string
	// mediaFolder overrides the one-shot collection's media directory.
	mediaFolder string
	// mediaDB overrides the one-shot collection's media database path.
	mediaDB string
	// timeout bounds startup and the selected native operation.
	timeout time.Duration
	// types resolves messages and Any payloads from the embedded upstream schema.
	types *dynamicpb.Types
}

// NewCommand constructs the complete generated command tree without starting Anki.
// The executable entry point supplies context and its standard streams through Cobra.
func NewCommand() (*cobra.Command, error) {
	set := new(descriptorpb.FileDescriptorSet)
	if err := proto.Unmarshal([]byte(schemaData), set); err != nil {
		return nil, fmt.Errorf("parse embedded Anki schema: %w", err)
	}
	files, err := protodesc.NewFiles(set)
	if err != nil {
		return nil, fmt.Errorf("resolve embedded Anki schema: %w", err)
	}
	o := &options{types: dynamicpb.NewTypes(files)}
	root := &cobra.Command{Use: "anki_rpc", Short: "Call Anki's protobuf services", SilenceErrors: true, SilenceUsage: true}
	root.CompletionOptions.DisableDefaultCmd = true
	flags := root.PersistentFlags()
	flags.StringVar(&o.worker, "worker", "", "Managed stdio server executable (default: Bazel runfile or anki_grpc in PATH)")
	flags.StringVar(&o.target, "target", "", "Existing plaintext gRPC target (host:port or unix:///path)")
	flags.StringVar(&o.collection, "collection", "", "Open this local collection for a managed one-shot call")
	flags.StringVar(&o.mediaFolder, "media-folder", "", "Local media directory (default: COLLECTION.media)")
	flags.StringVar(&o.mediaDB, "media-db", "", "Local media database (default: COLLECTION.media.db)")
	flags.DurationVar(&o.timeout, "timeout", 30*time.Second, "Startup and RPC deadline")
	root.MarkFlagsMutuallyExclusive("target", "worker")
	root.MarkFlagsMutuallyExclusive("target", "collection")
	root.MarkFlagsMutuallyExclusive("target", "media-folder")
	root.MarkFlagsMutuallyExclusive("target", "media-db")
	if err := addGenerated(root, o, files); err != nil {
		return nil, fmt.Errorf("register Anki commands: %w", err)
	}
	return root, nil
}

// addService attaches a generated service and its explicitly generated method names.
func (o *options) addService(root *cobra.Command, files *protoregistry.Files, name string, methods []string) error {
	descriptor, err := files.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil {
		return fmt.Errorf("find service %q: %w", name, err)
	}
	service, ok := descriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		return fmt.Errorf("descriptor %q is not a service", name)
	}
	command := &cobra.Command{Use: name, Short: "Anki service " + name}
	for _, name := range methods {
		method := service.Methods().ByName(protoreflect.Name(name))
		if method == nil {
			return fmt.Errorf("service %q has no method %q", service.FullName(), name)
		}
		command.AddCommand(o.methodCommand(method))
	}
	root.AddCommand(command)
	return nil
}

// methodCommand builds one JSON command around the upstream request/response contract.
func (o *options) methodCommand(method protoreflect.MethodDescriptor) *cobra.Command {
	var input string
	command := &cobra.Command{
		Use: string(method.Name()), Short: fmt.Sprintf("%s → %s", method.Input().FullName(), method.Output().FullName()),
		Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if o.timeout <= 0 {
				return fmt.Errorf("timeout must be positive")
			}
			if o.collection == "" && (o.mediaFolder != "" || o.mediaDB != "") {
				return fmt.Errorf("media paths require --collection")
			}
			data, err := readRequest(command.InOrStdin(), input)
			if err != nil {
				return fmt.Errorf("read request: %w", err)
			}
			request := dynamicpb.NewMessage(method.Input())
			if err := (protojson.UnmarshalOptions{Resolver: o.types}).Unmarshal(data, request); err != nil {
				return fmt.Errorf("parse %s: %w", method.Input().FullName(), err)
			}
			response := dynamicpb.NewMessage(method.Output())
			ctx, cancel := context.WithTimeout(command.Context(), o.timeout)
			defer cancel()
			if err := o.call(ctx, method, request, response); err != nil {
				return fmt.Errorf("call %s: %w", method.FullName(), err)
			}
			output, err := (protojson.MarshalOptions{UseProtoNames: true, Indent: "  ", Resolver: o.types}).Marshal(response)
			if err != nil {
				return fmt.Errorf("encode response: %w", err)
			}
			if _, err := fmt.Fprintln(command.OutOrStdout(), string(output)); err != nil {
				return fmt.Errorf("write response: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&input, "request", "{}", "ProtoJSON request, @file, or - for stdin")
	return command
}

// readRequest obtains JSON without interpreting protobuf fields itself.
func readRequest(stdin io.Reader, input string) ([]byte, error) {
	if input == "-" {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return data, nil
	}
	if strings.HasPrefix(input, "@") {
		data, err := os.ReadFile(strings.TrimPrefix(input, "@"))
		if err != nil {
			return nil, fmt.Errorf("read request file: %w", err)
		}
		return data, nil
	}
	return []byte(input), nil
}

// call owns one connection and delegates the request through its exact upstream path.
func (o *options) call(ctx context.Context, method protoreflect.MethodDescriptor, request, response proto.Message) (err error) {
	path := "/" + string(method.Parent().FullName()) + "/" + string(method.Name())
	if o.target != "" {
		connection, err := grpc.NewClient(o.target, grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
		if err != nil {
			return fmt.Errorf("connect to gRPC target: %w", err)
		}
		defer func() {
			if closeErr := connection.Close(); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close gRPC connection: %w", closeErr))
			}
		}()
		if err := o.invoke(ctx, connection, path, request, response); err != nil {
			return fmt.Errorf("invoke target: %w", err)
		}
		return nil
	}
	binary, err := workerPath(o.worker)
	if err != nil {
		return fmt.Errorf("resolve managed worker: %w", err)
	}
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	var collections ankicol.BackendCollectionServiceClient
	opened := false
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if opened {
			if _, closeErr := collections.CloseCollection(cleanup, &ankicol.CloseCollectionRequest{}); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close local collection: %w", closeErr))
			}
		}
		if stopErr := worker.Stop(cleanup); stopErr != nil {
			err = errors.Join(err, fmt.Errorf("stop managed worker: %w", stopErr))
		}
	}()
	if err := worker.Start(ctx); err != nil {
		return fmt.Errorf("start managed worker: %w", err)
	}
	connection, err := worker.Connection()
	if err != nil {
		return fmt.Errorf("get managed connection: %w", err)
	}
	collections = ankicol.NewBackendCollectionServiceClient(connection)
	if o.collection != "" {
		database, err := filepath.Abs(o.collection)
		if err != nil {
			return fmt.Errorf("resolve local collection: %w", err)
		}
		mediaFolder := o.mediaFolder
		if mediaFolder == "" {
			mediaFolder = database + ".media"
		}
		mediaDB := o.mediaDB
		if mediaDB == "" {
			mediaDB = database + ".media.db"
		}
		if _, err := collections.OpenCollection(ctx, &ankicol.OpenCollectionRequest{CollectionPath: database, MediaFolderPath: mediaFolder, MediaDbPath: mediaDB}); err != nil {
			return fmt.Errorf("open local collection: %w", err)
		}
		opened = true
	}
	// These native operations take the open collection even when they fail.
	if path == "/anki.collection.BackendCollectionService/CloseCollection" || path == "/anki.import_export.BackendImportExportService/ExportCollectionPackage" {
		opened = false
	}
	if err := o.invoke(ctx, connection, path, request, response); err != nil {
		return fmt.Errorf("invoke managed worker: %w", err)
	}
	if path == "/anki.collection.BackendCollectionService/OpenCollection" {
		opened = true
	}
	return nil
}

// invoke preserves transport errors and presents native error details on any transport.
func (o *options) invoke(ctx context.Context, connection grpc.ClientConnInterface, path string, request, response proto.Message) error {
	var trailers metadata.MD
	if err := connection.Invoke(ctx, path, request, response, grpc.Trailer(&trailers)); err != nil {
		details := trailers.Get("anki-error-bin")
		if o.target != "" && len(details) > 0 {
			messageType, lookupErr := o.types.FindMessageByName("anki.backend.BackendError")
			if lookupErr == nil {
				native := messageType.New()
				if decodeErr := proto.Unmarshal([]byte(details[0]), native.Interface()); decodeErr == nil {
					field := native.Descriptor().Fields().ByName("message")
					if field != nil {
						return fmt.Errorf("%s: %w", native.Get(field).String(), err)
					}
				}
			}
		}
		return fmt.Errorf("gRPC %s: %w", path, err)
	}
	return nil
}

// workerPath honors an explicit binary, then Bazel runtime data, then PATH.
func workerPath(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	path, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err == nil {
		if _, statErr := os.Stat(path); statErr == nil {
			return path, nil
		}
	}
	path, err = exec.LookPath("anki_grpc")
	if err != nil {
		return "", fmt.Errorf("find anki_grpc; pass --worker outside Bazel: %w", err)
	}
	return path, nil
}

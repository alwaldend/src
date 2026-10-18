package rpc

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Method retains an exposed route and its authoritative request/response schema.
type Method struct {
	// Name is the kebab-case method command.
	Name string
	// Path is the complete gRPC method path, including implicit backend routes.
	Path string
	// Descriptor supplies the request and response types.
	Descriptor protoreflect.MethodDescriptor
}

// Service groups the methods exposed by one upstream service.
type Service struct {
	// Name is the kebab-case service command without the Service suffix.
	Name string
	// FullName is the original fully qualified protobuf service name.
	FullName protoreflect.FullName
	// Methods contains sorted declared and implicitly forwarded operations.
	Methods []Method
}

var (
	// acronymBoundary separates an uppercase acronym from the following word.
	acronymBoundary = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	// wordBoundary separates ordinary camel-case words.
	wordBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// commandName converts protobuf identifiers into stable CLI names.
func commandName(name string) string {
	name = acronymBoundary.ReplaceAllString(name, "${1}-${2}")
	name = wordBoundary.ReplaceAllString(name, "${1}-${2}")
	return strings.ToLower(strings.ReplaceAll(name, "_", "-"))
}

// Services reads Anki descriptors registered by the generated schema imports.
func Services() ([]Service, error) {
	descriptors := map[protoreflect.FullName]protoreflect.ServiceDescriptor{}
	protoregistry.GlobalFiles.RangeFiles(func(file protoreflect.FileDescriptor) bool {
		if strings.HasPrefix(file.Path(), "anki/") {
			for index := 0; index < file.Services().Len(); index++ {
				service := file.Services().Get(index)
				descriptors[service.FullName()] = service
			}
		}
		return true
	})
	var services []Service
	names := map[string]bool{}
	for fullName, descriptor := range descriptors {
		service := Service{Name: commandName(strings.TrimSuffix(string(descriptor.Name()), "Service")), FullName: fullName}
		if names[service.Name] {
			return nil, fmt.Errorf("duplicate Anki service command %q", service.Name)
		}
		names[service.Name] = true
		methods := map[protoreflect.Name]protoreflect.MethodDescriptor{}
		if strings.HasPrefix(string(descriptor.Name()), "Backend") {
			collectionName := protoreflect.FullName(string(descriptor.ParentFile().Package()) + "." + strings.TrimPrefix(string(descriptor.Name()), "Backend"))
			if collection := descriptors[collectionName]; collection != nil {
				for index := 0; index < collection.Methods().Len(); index++ {
					method := collection.Methods().Get(index)
					methods[method.Name()] = method
				}
			}
		}
		for index := 0; index < descriptor.Methods().Len(); index++ {
			method := descriptor.Methods().Get(index)
			methods[method.Name()] = method
		}
		for name, descriptor := range methods {
			if descriptor.IsStreamingClient() || descriptor.IsStreamingServer() {
				return nil, fmt.Errorf("Anki method %s/%s requires unsupported streaming", fullName, name)
			}
			for _, message := range []protoreflect.MessageDescriptor{descriptor.Input(), descriptor.Output()} {
				if _, err := protoregistry.GlobalTypes.FindMessageByName(message.FullName()); err != nil {
					return nil, fmt.Errorf("resolve Anki message %s: %w", message.FullName(), err)
				}
			}
			service.Methods = append(service.Methods, Method{Name: commandName(string(name)), Path: "/" + string(fullName) + "/" + string(name), Descriptor: descriptor})
		}
		sort.Slice(service.Methods, func(i, j int) bool { return service.Methods[i].Name < service.Methods[j].Name })
		services = append(services, service)
	}
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })
	return services, nil
}

// Resolve finds a supported route, including implicit backend forwarding methods.
func Resolve(path string) (Method, error) {
	services, err := Services()
	if err != nil {
		return Method{}, fmt.Errorf("read Anki catalog: %w", err)
	}
	for _, service := range services {
		for _, method := range service.Methods {
			if method.Path == path {
				return method, nil
			}
		}
	}
	return Method{}, fmt.Errorf("unknown Anki method %q", path)
}

// Prepare validates payload formats and returns a call that borrows a prepared connection.
// Request decoding occurs before the caller starts a worker or opens collection storage.
func Prepare(method Method, input []byte, inputFormat, outputFormat string) (func(context.Context, grpc.ClientConnInterface) ([]byte, error), error) {
	requestType, err := protoregistry.GlobalTypes.FindMessageByName(method.Descriptor.Input().FullName())
	if err != nil {
		return nil, fmt.Errorf("resolve %s request: %w", method.Path, err)
	}
	responseType, err := protoregistry.GlobalTypes.FindMessageByName(method.Descriptor.Output().FullName())
	if err != nil {
		return nil, fmt.Errorf("resolve %s response: %w", method.Path, err)
	}
	request, response := requestType.New().Interface(), responseType.New().Interface()
	if len(strings.TrimSpace(string(input))) == 0 && method.Descriptor.Input().Fields().Len() == 0 {
		if inputFormat == "json" {
			input = []byte("{}")
		}
	}
	if err := codec.Decode(input, codec.Format(inputFormat), request); err != nil {
		return nil, fmt.Errorf("decode %s request: %w", method.Path, err)
	}
	if _, err := codec.Encode(response, codec.Format(outputFormat)); err != nil {
		return nil, fmt.Errorf("validate %s response format: %w", method.Path, err)
	}
	return func(ctx context.Context, connection grpc.ClientConnInterface) ([]byte, error) {
		if err := connection.Invoke(ctx, method.Path, request, response); err != nil {
			return nil, fmt.Errorf("invoke Anki %s: %w", method.Path, err)
		}
		output, err := codec.Encode(response, codec.Format(outputFormat))
		if err != nil {
			return nil, fmt.Errorf("encode %s response: %w", method.Path, err)
		}
		return output, nil
	}, nil
}

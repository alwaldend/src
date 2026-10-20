package ankiclient

import (
	"context"
	"errors"
	"fmt"

	ankibackend "github.com/ankitects/anki/proto/anki/backend"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// ErrorTrailer identifies Anki's serialized BackendError in response metadata.
const ErrorTrailer = "anki-error-bin"

// DecodeCallError preserves gRPC status and native BackendError details for every transport.
func DecodeCallError(ctx context.Context, method string, err error, trailers metadata.MD) error {
	if err == nil {
		return nil
	}
	if ctx.Err() != nil {
		return fmt.Errorf("call Anki %s: %w", method, ctx.Err())
	}
	data := trailers.Get(ErrorTrailer)
	if len(data) != 1 {
		return fmt.Errorf("call Anki %s: %w", method, err)
	}
	var failure ankibackend.BackendError
	if decodeError := proto.Unmarshal([]byte(data[0]), &failure); decodeError != nil {
		return fmt.Errorf("decode Anki %s failure: %w", method, errors.Join(err, decodeError))
	}
	return fmt.Errorf("call Anki %s: %s: %s: %w", method, failure.Kind, failure.Message, err)
}

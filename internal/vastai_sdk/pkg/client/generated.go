package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	generatedapi "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/internal/generated"
)

// generatedTransport adapts schema paths to the configured API root and bounds
// response reads before the generated client buffers or decodes their bodies.
type generatedTransport struct {
	// path is the configured API root, including any test endpoint prefix.
	path string
	// escapedPath preserves encoded segment boundaries in the configured API root.
	escapedPath string
	// key is the credential for this independently configured client.
	key string
}

// RoundTrip preserves the configured API root, authentication, and canonical paths.
func (t *generatedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request = request.Clone(request.Context())
	escapedPath := t.escapedPath + strings.TrimPrefix(request.URL.EscapedPath(), "/api/v0")
	request.URL.Path = t.path + strings.TrimPrefix(request.URL.Path, "/api/v0")
	request.URL.RawPath = escapedPath
	if !strings.HasSuffix(escapedPath, "/") {
		request.URL.Path += "/"
		request.URL.RawPath += "/"
	}
	request.Header.Set("Authorization", "Bearer "+t.key)
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, fmt.Errorf("send generated API request: %w", err)
	}
	response.Body = &limitedBody{reader: io.LimitReader(response.Body, (8<<20)+1), body: response.Body}
	return response, nil
}

// limitedBody caps buffering while retaining the original response's close operation.
type limitedBody struct {
	// reader stops after the size limit and one byte used to detect overflow.
	reader io.Reader
	// body owns the underlying connection's response stream.
	body io.ReadCloser
}

// Read supplies only the bounded prefix to generated response decoders.
func (b *limitedBody) Read(data []byte) (int, error) {
	n, err := b.reader.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {
		return n, fmt.Errorf("read bounded API response: %w", err)
	}
	return n, err
}

// Close releases the underlying response stream.
func (b *limitedBody) Close() error {
	if err := b.body.Close(); err != nil {
		return fmt.Errorf("close bounded API response: %w", err)
	}
	return nil
}

// generatedError retains the original cause without rendering response-derived text.
type generatedError struct {
	// cause preserves transport or generation failure identity for callers.
	cause error
}

// Error supplies a fixed credential-safe diagnostic for generated failures.
func (e *generatedError) Error() string {
	return "generated API request failed; inspect the account before retrying a mutation"
}

// Unwrap preserves errors.Is and errors.As without disclosing the cause's message.
func (e *generatedError) Unwrap() error { return e.cause }

// generated uses schema-generated request serialization and endpoint selection.
// The upstream client retains raw response bodies; the adapter validates those
// because documented response variants differ from some generated model shapes.
func generated[T any](c *Client, method, path string, result any, requireSuccess bool, execute func() (T, *http.Response, error)) error {
	_, response, err := execute()
	var apiError *generatedapi.GenericOpenAPIError
	if response == nil || (err != nil && !errors.As(err, &apiError)) {
		if err == nil {
			return errors.New("generated API request returned no response")
		}
		return fmt.Errorf("execute %s %s: %w", method, path, &generatedError{cause: err})
	}
	if err := c.decodeResponse(method, path, response, result, requireSuccess); err != nil {
		return fmt.Errorf("validate generated %s %s response: %w", method, path, err)
	}
	return nil
}

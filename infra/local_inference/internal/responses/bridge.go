package responses

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// maxBody bounds a single request, JSON reply, or SSE event in memory.
const maxBody = 32 << 20

// identity records the original Responses namespace and local function name.
type identity struct {
	// namespace distinguishes functions sharing a local name.
	namespace string
	// name is the function name expected by Codex.
	name string
}

// flatName produces a deterministic name valid within the 64-character limit.
func flatName(namespace, name string) string {
	digest := sha256.Sum256([]byte(namespace + "\x00" + name))
	return "ns_" + hex.EncodeToString(digest[:24])
}

// object validates an external JSON object before adapting its fields.
func object(value any) (map[string]any, error) {
	result, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("expected JSON object")
	}
	return result, nil
}

// decode reads bounded JSON while retaining exact numeric values.
func decode(data []byte) (map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value map[string]any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	if value == nil {
		return nil, errors.New("expected JSON object")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("read trailing JSON: %w", err)
		}
		return nil, errors.New("expected one JSON object")
	}
	return value, nil
}

// flatten translates tool schemas and replayed calls, rejecting silent tool loss.
func flatten(body map[string]any) (map[string]identity, error) {
	aliases := make(map[string]identity)
	used := make(map[string]bool)
	tools := []any{}
	entries, ok := body["tools"].([]any)
	if body["tools"] != nil && !ok {
		return nil, errors.New("tools must be an array")
	}
	for _, entry := range entries {
		tool, err := object(entry)
		if err != nil {
			return nil, fmt.Errorf("read tool: %w", err)
		}
		if tool["type"] == "function" {
			name, _ := tool["name"].(string)
			if name == "" || used[name] {
				return nil, errors.New("empty or duplicate function name")
			}
			used[name] = true
			tools = append(tools, tool)
			continue
		}
		if tool["type"] != "namespace" {
			return nil, errors.New("only function tools and function namespaces are supported")
		}
		namespace, _ := tool["name"].(string)
		nested, ok := tool["tools"].([]any)
		if namespace == "" || !ok {
			return nil, errors.New("namespace requires a name and tool array")
		}
		for _, entry := range nested {
			function, err := object(entry)
			if err != nil {
				return nil, fmt.Errorf("read namespace tool: %w", err)
			}
			name, _ := function["name"].(string)
			if function["type"] != "function" || name == "" {
				return nil, errors.New("namespace contains an unsupported function schema")
			}
			alias := flatName(namespace, name)
			if used[alias] {
				return nil, errors.New("duplicate or colliding function identity")
			}
			used[alias] = true
			aliases[alias] = identity{namespace: namespace, name: name}
			function["name"] = alias
			description, _ := function["description"].(string)
			groupDescription, _ := tool["description"].(string)
			function["description"] = strings.TrimSpace(namespace + "." + name + "\n" + groupDescription + "\n" + description)
			tools = append(tools, function)
		}
	}
	if body["tools"] != nil {
		body["tools"] = tools
	}
	if input, ok := body["input"].([]any); ok {
		for _, entry := range input {
			if call, ok := entry.(map[string]any); ok && call["type"] == "function_call" {
				flattenCall(call)
			}
		}
	}
	if choice, ok := body["tool_choice"].(map[string]any); ok {
		if choice["type"] == "function" {
			flattenCall(choice)
		}
		if allowed, ok := choice["tools"].([]any); ok {
			for _, entry := range allowed {
				if call, ok := entry.(map[string]any); ok {
					flattenCall(call)
				}
			}
		}
	}
	return aliases, nil
}

// flattenCall translates namespaced call history and forced tool selections.
func flattenCall(call map[string]any) {
	namespace, _ := call["namespace"].(string)
	name, _ := call["name"].(string)
	if namespace != "" && name != "" {
		call["name"] = flatName(namespace, name)
		delete(call, "namespace")
	}
}

// restore walks protocol containers without altering user text or arguments.
func restore(value any, aliases map[string]identity) {
	switch value := value.(type) {
	case []any:
		for _, item := range value {
			restore(item, aliases)
		}
	case map[string]any:
		typ, _ := value["type"].(string)
		if typ == "function_call" || strings.HasPrefix(typ, "response.function_call_arguments.") {
			name, _ := value["name"].(string)
			if id, ok := aliases[name]; ok {
				value["name"] = id.name
				value["namespace"] = id.namespace
			}
		}
		for _, key := range []string{"item", "output", "response"} {
			if child, ok := value[key]; ok {
				restore(child, aliases)
			}
		}
	}
}

// streamBody closes both sides of a streaming translation when the client leaves.
type streamBody struct {
	// ReadCloser is the translated stream read by ReverseProxy.
	io.ReadCloser
	// upstream is the original inference response body.
	upstream io.Closer
}

// Close cancels the original body and releases a blocked translation writer.
func (body *streamBody) Close() error {
	if err := errors.Join(body.ReadCloser.Close(), body.upstream.Close()); err != nil {
		return fmt.Errorf("close response stream: %w", err)
	}
	return nil
}

// translateStream forwards complete SSE events, preserving event/id metadata.
func translateStream(source io.ReadCloser, dest *io.PipeWriter, aliases map[string]identity) {
	defer source.Close()
	err := rewriteEvents(source, dest, aliases)
	_ = dest.CloseWithError(err)
}

// rewriteEvents transforms JSON data fields with bounded event buffering.
func rewriteEvents(source io.Reader, dest io.Writer, aliases map[string]identity) error {
	scanner := bufio.NewScanner(source)
	scanner.Buffer(make([]byte, 4096), maxBody)
	var fields, data []string
	size := 0
	flush := func() error {
		if len(data) > 0 {
			payload := strings.Join(data, "\n")
			if payload != "[DONE]" {
				value, err := decode([]byte(payload))
				if err != nil {
					return fmt.Errorf("decode inference event: %w", err)
				}
				restore(value, aliases)
				encoded, err := json.Marshal(value)
				if err != nil {
					return fmt.Errorf("encode inference event: %w", err)
				}
				payload = string(encoded)
			}
			fields = append(fields, "data: "+payload)
		}
		if _, err := io.WriteString(dest, strings.Join(fields, "\n")+"\n\n"); err != nil {
			return fmt.Errorf("write inference event: %w", err)
		}
		fields = nil
		data = nil
		size = 0
		return nil
	}
	for scanner.Scan() {
		line := scanner.Text()
		size += len(line)
		if size > maxBody {
			return errors.New("inference event exceeds size limit")
		}
		if line == "" {
			if err := flush(); err != nil {
				return fmt.Errorf("forward inference event: %w", err)
			}
			continue
		}
		if strings.HasPrefix(line, "data:") {
			data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		} else {
			fields = append(fields, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read inference stream: %w", err)
	}
	if len(fields)+len(data) > 0 {
		if err := flush(); err != nil {
			return fmt.Errorf("forward final inference event: %w", err)
		}
	}
	return nil
}

// Handler constructs a fixed-upstream loopback proxy with request-local mappings.
func Handler(upstream *url.URL) http.Handler {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true
	transport.ResponseHeaderTimeout = 10 * time.Minute
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" || !((r.Method == http.MethodPost && r.URL.Path == "/v1/responses") || (r.Method == http.MethodGet && (r.URL.Path == "/health" || r.URL.Path == "/v1/models"))) {
			http.NotFound(w, r)
			return
		}
		aliases := make(map[string]identity)
		if r.Method == http.MethodPost {
			data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
			if err != nil {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			body, err := decode(data)
			if err == nil {
				aliases, err = flatten(body)
			}
			if err != nil {
				http.Error(w, "unsupported Responses request: "+err.Error(), http.StatusBadRequest)
				return
			}
			data, err = json.Marshal(body)
			if err != nil {
				http.Error(w, "cannot encode Responses request", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(data))
			r.ContentLength = int64(len(data))
		}
		proxy := httputil.NewSingleHostReverseProxy(upstream)
		proxy.Transport = transport
		proxy.FlushInterval = -1
		proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
			http.Error(w, "local inference unavailable", http.StatusBadGateway)
		}
		proxy.ModifyResponse = func(response *http.Response) error {
			if r.Method != http.MethodPost || response.StatusCode >= 300 {
				return nil
			}
			response.Header.Del("Content-Length")
			response.ContentLength = -1
			if strings.HasPrefix(response.Header.Get("Content-Type"), "text/event-stream") {
				reader, writer := io.Pipe()
				original := response.Body
				response.Body = &streamBody{ReadCloser: reader, upstream: original}
				go translateStream(original, writer, aliases)
				return nil
			}
			defer response.Body.Close()
			data, err := io.ReadAll(io.LimitReader(response.Body, maxBody+1))
			if err != nil {
				return fmt.Errorf("read inference JSON: %w", err)
			}
			if len(data) > maxBody {
				return errors.New("inference JSON exceeds size limit")
			}
			body, err := decode(data)
			if err != nil {
				return fmt.Errorf("decode inference JSON: %w", err)
			}
			restore(body, aliases)
			data, err = json.Marshal(body)
			if err != nil {
				return fmt.Errorf("encode inference JSON: %w", err)
			}
			response.Body = io.NopCloser(bytes.NewReader(data))
			return nil
		}
		proxy.ServeHTTP(w, r)
	})
}

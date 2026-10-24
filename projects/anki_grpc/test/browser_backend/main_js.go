//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/browser"
	"github.com/ankitects/anki/proto/anki/backend"
	"github.com/ankitects/anki/proto/anki/card_rendering"
	"github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/generic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// verification captures observable behavior for the runner's persistent report.
type verification struct {
	// Storage records the independent Chromium shared-filesystem workflow.
	Storage *storageVerification `json:"storage,omitempty"`
	// Cases names the checks completed by the actual generated-client workflow.
	Cases []string `json:"cases"`
	// Unicode is Anki's transformed response, retaining non-ASCII text.
	Unicode string `json:"unicode"`
	// PayloadBytes records the repeated request's UTF-8 byte length.
	PayloadBytes int `json:"payloadBytes"`
	// RepeatedCalls counts successful large request and response transfers.
	RepeatedCalls int `json:"repeatedCalls"`
	// MemoryBefore records Rust memory after warming its allocator.
	MemoryBefore int `json:"memoryBefore"`
	// MemoryAfter records Rust memory after repeated commands release their buffers.
	MemoryAfter int `json:"memoryAfter"`
	// NativeError contains the decoded native failure from a generated call.
	NativeError string `json:"nativeError"`
	// Error records a failure while retaining already completed evidence.
	Error string `json:"error,omitempty"`
}

// managedClient adds an observable cleanup boundary to the shared client contract.
type managedClient interface {
	ankiclient.StartableClient
	Done() <-chan struct{}
}

// main retains Go callbacks until the JavaScript runner has inspected the results.
func main() {
	report := verification{}
	finished := make(chan struct{})
	finish := js.FuncOf(func(js.Value, []js.Value) any {
		close(finished)
		return nil
	})
	defer finish.Release()
	js.Global().Set("ankiBrowserBackendFinish", finish)
	if js.Global().Get("ankiBrowserBackendStorageMode").Truthy() {
		if err := verifyStorage(&report); err != nil {
			report.Error = err.Error()
		}
	} else if err := verify(browser.NewClient(js.Global().Get("ankiBrowserBackendFactory")), &report); err != nil {
		report.Error = err.Error()
	} else if err := verifyStartup(&report); err != nil {
		report.Error = err.Error()
	} else if err := verifyTrap(&report); err != nil {
		report.Error = err.Error()
	}
	data, err := json.Marshal(report)
	if err != nil {
		panic(fmt.Errorf("encode browser backend verification: %w", err))
	}
	js.Global().Set("ankiBrowserBackendReport", string(data))
	js.Global().Call("ankiBrowserBackendVerified")
	<-finished
}

// verify exercises only collection-free Anki operations through the production adapter.
func verify(client managedClient, report *verification) (result error) {
	ctx, finish := context.WithTimeout(context.Background(), 30*time.Second)
	defer finish()
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := client.Connection(); err == nil {
		return fmt.Errorf("connection available before startup")
	}
	startup, finishStartup := context.WithCancel(ctx)
	defer finishStartup()
	if err := client.Start(startup); err != nil {
		return fmt.Errorf("start direct browser client: %w", err)
	}
	defer func() {
		if err := client.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("release direct browser client: %w", err))
		}
	}()
	finishStartup()
	connection, err := client.Connection()
	if err != nil {
		return fmt.Errorf("borrow started backend connection: %w", err)
	}
	if err := client.Start(ctx); err == nil {
		return fmt.Errorf("second startup succeeded")
	}
	report.Cases = append(report.Cases, "explicit startup and borrowed connection")

	rendering := card_rendering.NewBackendCardRenderingServiceClient(connection)
	text := "你好, café 🗂️"
	response, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "<b>" + text + "</b>"})
	if err != nil {
		return fmt.Errorf("strip Unicode HTML through generated client: %w", err)
	}
	if response.GetVal() != text {
		return fmt.Errorf("Unicode response changed: got %q", response.GetVal())
	}
	report.Unicode = response.GetVal()
	report.Cases = append(report.Cases, "generated unary call transforms Unicode")

	response, err = rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{})
	if err != nil {
		return fmt.Errorf("transfer empty request and response: %w", err)
	}
	if response.GetVal() != "" {
		return fmt.Errorf("empty response gained content")
	}
	report.Cases = append(report.Cases, "empty protobuf buffers retain their meaning")

	payload := strings.Repeat(text, 49152)
	report.PayloadBytes = len(payload)
	for index := range 20 {
		response, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "<i>" + payload + "</i>"})
		if err != nil {
			return fmt.Errorf("transfer large Unicode payload %d: %w", index, err)
		}
		if response.GetVal() != payload {
			return fmt.Errorf("large response %d was truncated or changed", index)
		}
		if index == 3 {
			report.MemoryBefore = js.Global().Call("ankiBrowserBackendMemoryBytes").Int()
		}
		if index >= 4 {
			report.RepeatedCalls++
		}
	}
	report.MemoryAfter = js.Global().Call("ankiBrowserBackendMemoryBytes").Int()
	if report.MemoryAfter-report.MemoryBefore > 4*report.PayloadBytes {
		return fmt.Errorf("Rust memory grew by %d bytes across %d repeated calls", report.MemoryAfter-report.MemoryBefore, report.RepeatedCalls)
	}
	report.Cases = append(report.Cases, "repeated large buffers remain bounded")

	if err := connection.Invoke(ctx, "/unknown.Service/Method", &generic.Empty{}, &generic.Empty{}); status.Code(err) != codes.Unimplemented {
		return fmt.Errorf("unknown route lost its status: %v", err)
	}
	var trailers metadata.MD
	_, err = config.NewConfigServiceClient(connection).GetConfigJson(ctx, &generic.String{Val: "absent"}, grpc.Trailer(&trailers))
	if status.Code(err) != codes.Unknown || len(trailers.Get(ankiclient.ErrorTrailer)) != 1 {
		return fmt.Errorf("native failure lost its status or trailer: %v", err)
	}
	var failure backend.BackendError
	if err := proto.Unmarshal([]byte(trailers.Get(ankiclient.ErrorTrailer)[0]), &failure); err != nil {
		return fmt.Errorf("decode native failure trailer: %w", err)
	}
	if failure.GetMessage() == "" || !strings.Contains(err.Error(), failure.GetMessage()) {
		return fmt.Errorf("native failure details missing from returned error: %v", err)
	}
	report.NativeError = failure.GetMessage()
	report.Cases = append(report.Cases, "unknown routes and native error details")

	if _, err := rendering.StripHtml(cancelled, &card_rendering.StripHtmlRequest{Text: "unused"}); !errors.Is(err, context.Canceled) {
		return fmt.Errorf("cancelled command lost its context error: %v", err)
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "<b>still ready</b>"}); err != nil {
		return fmt.Errorf("valid call after rejected commands: %w", err)
	}
	report.Cases = append(report.Cases, "rejected commands retain a usable backend")

	var headers metadata.MD
	trailers = metadata.Pairs("stale", "metadata")
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "<b>options</b>"}, grpc.Header(&headers), grpc.Trailer(&trailers)); err != nil {
		return fmt.Errorf("use supported direct call options: %w", err)
	}
	if len(headers) != 0 || len(trailers) != 0 {
		return fmt.Errorf("successful direct call retained stale metadata")
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "oversized"}, grpc.MaxCallSendMsgSize(1)); status.Code(err) != codes.ResourceExhausted {
		return fmt.Errorf("request size option was ignored: %v", err)
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "oversized"}, grpc.MaxCallRecvMsgSize(1)); status.Code(err) != codes.ResourceExhausted {
		return fmt.Errorf("response size option was ignored: %v", err)
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{}, grpc.UseCompressor("gzip")); status.Code(err) != codes.Unimplemented {
		return fmt.Errorf("unsupported compressor option was silently accepted: %v", err)
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{}, grpc.WaitForReady(false)); status.Code(err) != codes.Unimplemented {
		return fmt.Errorf("unsupported transport readiness option was silently accepted: %v", err)
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{}); err != nil {
		return fmt.Errorf("valid call after option rejections: %w", err)
	}
	report.Cases = append(report.Cases, "call options preserve metadata and enforce supported behavior")

	if err := client.Stop(ctx); err != nil {
		return fmt.Errorf("stop initialized backend: %w", err)
	}
	if _, err := client.Connection(); err == nil {
		return fmt.Errorf("connection available after shutdown")
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{Text: "unused"}); err == nil {
		return fmt.Errorf("retained connection accepted a call after shutdown")
	}
	if err := client.Start(ctx); err == nil {
		return fmt.Errorf("startup succeeded after shutdown")
	}
	if err := client.Stop(ctx); err != nil {
		return fmt.Errorf("repeat backend shutdown: %w", err)
	}
	if err := waitDone(ctx, client); err != nil {
		return fmt.Errorf("observe normal backend cleanup: %w", err)
	}
	report.Cases = append(report.Cases, "terminal client and repeatable cleanup")
	return nil
}

// verifyStartup observes rejected factories and joined initialization after cancellation.
func verifyStartup(report *verification) (failure error) {
	ctx, finish := context.WithTimeout(context.Background(), 30*time.Second)
	defer finish()
	failed := browser.NewClient(js.Global().Get("ankiBrowserBackendFailedFactory"))
	if err := failed.Start(ctx); err == nil || !strings.Contains(err.Error(), "fixture startup failure") {
		return fmt.Errorf("factory failure was not observable: %v", err)
	}
	if _, err := failed.Connection(); err == nil {
		return fmt.Errorf("connection available after factory failure")
	}
	// Stop may preserve the startup error, but it must remain terminal and complete.
	_ = failed.Stop(ctx)
	if err := waitDone(ctx, failed); err != nil {
		return fmt.Errorf("observe failed startup cleanup: %w", err)
	}
	report.Cases = append(report.Cases, "factory rejection remains observable and releases startup resources")

	started := make(chan struct{})
	notify := js.FuncOf(func(js.Value, []js.Value) any {
		close(started)
		return nil
	})
	defer notify.Release()
	js.Global().Set("ankiBrowserBackendFactoryStarted", notify)
	delayed := browser.NewClient(js.Global().Get("ankiBrowserBackendDelayedFactory"))
	startup, cancel := context.WithCancel(ctx)
	defer cancel()
	result := make(chan error, 1)
	startDone := make(chan struct{})
	// Drain the fixture even when an assertion fails before factory release, so
	// its Promise callbacks cannot resume Go after the runner finishes main.
	defer func() {
		cancel()
		js.Global().Call("ankiBrowserBackendReleaseFactory")
		cleanup, finishCleanup := context.WithTimeout(context.Background(), 10*time.Second)
		defer finishCleanup()
		select {
		case <-startDone:
		case <-cleanup.Done():
			failure = errors.Join(failure, fmt.Errorf("join delayed startup after verification: %w", cleanup.Err()))
		}
		if err := delayed.Stop(cleanup); err != nil && !errors.Is(err, context.Canceled) {
			failure = errors.Join(failure, fmt.Errorf("stop delayed startup after verification: %w", err))
		}
		if err := waitDone(cleanup, delayed); err != nil {
			failure = errors.Join(failure, fmt.Errorf("drain delayed startup after verification: %w", err))
		}
	}()
	go func() {
		defer close(startDone)
		result <- delayed.Start(startup)
	}()
	select {
	case <-started:
	case <-ctx.Done():
		return fmt.Errorf("wait for delayed factory: %w", ctx.Err())
	}
	cancel()
	// Give cancellation a scheduling interval while the JavaScript Promise stays
	// unresolved. Start must join that factory rather than return on ctx.Done.
	pending := time.NewTimer(50 * time.Millisecond)
	defer pending.Stop()
	select {
	case err := <-result:
		return fmt.Errorf("startup returned before its canceled factory settled: %v", err)
	case <-pending.C:
	case <-ctx.Done():
		return fmt.Errorf("observe pending canceled startup: %w", ctx.Err())
	}
	if _, err := delayed.Connection(); err == nil {
		return fmt.Errorf("unresolved canceled factory published a connection")
	}
	if err := delayed.Stop(startup); !errors.Is(err, context.Canceled) {
		return fmt.Errorf("pending factory shutdown lost its context error: %v", err)
	}
	select {
	case err := <-result:
		return fmt.Errorf("concurrent Stop released Start before its factory settled: %v", err)
	default:
	}
	select {
	case <-delayed.Done():
		return fmt.Errorf("cleanup completed before pending factory settled")
	default:
	}
	js.Global().Call("ankiBrowserBackendReleaseFactory")
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			return fmt.Errorf("completed canceled startup lost its context error: %v", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("wait for completed canceled startup: %w", ctx.Err())
	}
	if err := waitDone(ctx, delayed); err != nil {
		return fmt.Errorf("join canceled factory cleanup: %w", err)
	}
	if js.Global().Call("ankiBrowserBackendDelayedCloses").Int() != 1 {
		return fmt.Errorf("canceled factory backend was not closed exactly once")
	}
	if _, err := delayed.Connection(); err == nil {
		return fmt.Errorf("canceled factory published a backend after completion")
	}
	_ = delayed.Stop(ctx)
	report.Cases = append(report.Cases, "canceled startup waits for factory completion and joins cleanup with concurrent Stop")
	return nil
}

// verifyTrap proves a thrown backend call permanently invalidates borrowed RPC clients.
func verifyTrap(report *verification) error {
	ctx, finish := context.WithTimeout(context.Background(), 30*time.Second)
	defer finish()
	client := browser.NewClient(js.Global().Get("ankiBrowserBackendTrapFactory"))
	if err := client.Start(ctx); err != nil {
		return fmt.Errorf("start trap fixture: %w", err)
	}
	connection, err := client.Connection()
	if err != nil {
		return fmt.Errorf("borrow trap fixture connection: %w", err)
	}
	rendering := card_rendering.NewBackendCardRenderingServiceClient(connection)
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{}); err == nil {
		return fmt.Errorf("trapped backend command succeeded")
	}
	if _, err := client.Connection(); err == nil {
		return fmt.Errorf("connection available after backend trap")
	}
	if _, err := rendering.StripHtml(ctx, &card_rendering.StripHtmlRequest{}); err == nil {
		return fmt.Errorf("retained client accepted a command after backend trap")
	}
	if err := client.Stop(ctx); err == nil {
		return fmt.Errorf("shutdown hid the terminal backend trap")
	}
	if err := waitDone(ctx, client); err != nil {
		return fmt.Errorf("observe trapped backend cleanup: %w", err)
	}
	report.Cases = append(report.Cases, "backend trap invalidates retained clients and survives shutdown")
	return nil
}

// waitDone bounds observation of descriptor and JavaScript callback cleanup.
func waitDone(ctx context.Context, client managedClient) error {
	select {
	case <-client.Done():
		return nil
	case <-ctx.Done():
		return fmt.Errorf("wait for backend cleanup: %w", ctx.Err())
	}
}

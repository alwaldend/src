//go:build js && wasm

package browser

import (
	"fmt"
	"strconv"
	"sync"
	"syscall/js"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/timezone"
)

// timezoneHost lends the browser's named Go timezone to one Rust backend.
// Filesystem operations never traverse these callbacks.
type timezoneHost struct {
	// endpoint exposes UTC and local timestamp operations to the factory.
	endpoint js.Value
	// utc owns the UTC offset callback until backend destruction.
	utc js.Func
	// local owns the packed local offsets callback until backend destruction.
	local js.Func
	// closeOnce prevents releasing either callback twice.
	closeOnce sync.Once
}

// newTimezoneHost resolves the browser zone using the shared Go timezone data.
func newTimezoneHost() (*timezoneHost, error) {
	zone := js.Global().Get("Intl").Get("DateTimeFormat").New().Call("resolvedOptions").Get("timeZone").String()
	location, err := time.LoadLocation(zone)
	if err != nil {
		return nil, fmt.Errorf("load browser timezone %q: %w", zone, err)
	}
	host := &timezoneHost{endpoint: js.Global().Get("Object").New()}
	host.utc = js.FuncOf(func(_ js.Value, arguments []js.Value) any {
		if len(arguments) != 1 || arguments[0].Type() != js.TypeString {
			return nil
		}
		seconds, err := strconv.ParseInt(arguments[0].String(), 10, 64)
		if err != nil {
			return nil
		}
		return int(timezone.OffsetForUTC(location, seconds))
	})
	host.local = js.FuncOf(func(_ js.Value, arguments []js.Value) any {
		if len(arguments) != 1 || arguments[0].Type() != js.TypeString {
			return nil
		}
		seconds, err := strconv.ParseInt(arguments[0].String(), 10, 64)
		if err != nil {
			return nil
		}
		return strconv.FormatUint(timezone.OffsetsForLocal(location, seconds), 10)
	})
	host.endpoint.Set("offsetForUTC", host.utc)
	host.endpoint.Set("offsetsForLocal", host.local)
	return host, nil
}

// Close releases callbacks only after Rust destruction and active commands finish.
func (host *timezoneHost) Close() {
	host.closeOnce.Do(func() {
		host.utc.Release()
		host.local.Release()
	})
}

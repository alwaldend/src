package ankigrpc

import (
	"fmt"
	"maps"
	"net"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
)

// Factory routes addresses through private registrations and an explicit default.
// Use NewFactory; do not copy a Factory after use. Registration and construction
// are concurrency safe, and constructor calls run outside the registry lock.
type Factory struct {
	// defaultAddress substitutes only an empty operation address.
	defaultAddress string
	// mutex protects constructor registrations.
	mutex sync.RWMutex
	// constructors retains factories, never started clients.
	constructors map[string]ClientFactory
}

// defaultFactory supplies built-in network constructors and opt-in package registrations.
var defaultFactory = &Factory{constructors: map[string]ClientFactory{
	"tcp":  ToClientFactory(ankiclient.NewGRPC),
	"unix": ToClientFactory(ankiclient.NewGRPC),
}}

// NewFactory snapshots built-in and opt-in package registrations. Subsequent
// registrations are independent. An empty default requires explicit addresses.
// Executables supply defaults and register their supported worker constructors.
func NewFactory(defaultAddress string) *Factory {
	defaultFactory.mutex.RLock()
	defer defaultFactory.mutex.RUnlock()
	return &Factory{defaultAddress: defaultAddress, constructors: maps.Clone(defaultFactory.constructors)}
}

// RegisterClientFactory installs a constructor for a lowercase URL scheme in the
// package default. Executables can instead register on their own Factory.
// Invalid schemes, nil functions and duplicate registrations panic.
func RegisterClientFactory(scheme string, factory ClientFactory) {
	defaultFactory.RegisterClientFactory(scheme, factory)
}

// RegisterClientFactory installs a constructor in this factory without changing
// other factories. Invalid schemes, nil functions and duplicates panic.
func (factory *Factory) RegisterClientFactory(scheme string, constructor ClientFactory) {
	parsed, err := url.Parse(scheme + ":")
	if constructor == nil || err != nil || scheme == "" || parsed.Scheme != scheme {
		panic(fmt.Sprintf("register Anki client factory: invalid scheme or factory %q", scheme))
	}
	factory.mutex.Lock()
	defer factory.mutex.Unlock()
	if _, exists := factory.constructors[scheme]; exists {
		panic(fmt.Sprintf("register Anki client factory: scheme %q is already registered", scheme))
	}
	factory.constructors[scheme] = constructor
}

// NewClient constructs an unstarted client through the package default factory.
// Empty addresses fail because the library owns no default worker entrypoint.
// The caller owns startup and cleanup of the returned client.
func NewClient(address string) (ankiclient.StartableClient, error) {
	client, err := defaultFactory.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("construct default Anki client: %w", err)
	}
	return client, nil
}

// NewClient selects a constructor, substituting this factory's configured default
// only for empty input. Constructors receive the selected address unchanged.
// No client is started by routing.
func (factory *Factory) NewClient(address string) (ankiclient.StartableClient, error) {
	if address == "" {
		address = factory.defaultAddress
	}
	scheme, err := clientScheme(address)
	if err != nil {
		return nil, fmt.Errorf("select Anki client: %w", err)
	}
	factory.mutex.RLock()
	constructor := factory.constructors[scheme]
	factory.mutex.RUnlock()
	if constructor == nil {
		return nil, fmt.Errorf("Anki client scheme %q is unavailable; its owner must register a client factory", scheme)
	}
	client, err := constructor(address)
	if err != nil {
		return nil, fmt.Errorf("construct Anki client for scheme %q: %w", scheme, err)
	}
	return client, nil
}

// clientScheme recognizes routing shorthand without interpreting transport data.
// Each constructor owns parsing its original address exactly once.
func clientScheme(address string) (string, error) {
	if address == "" {
		return "", fmt.Errorf("Anki client address requires an explicit endpoint or registered scheme")
	}
	if filepath.IsAbs(address) {
		return "unix", nil
	}
	if host, port, err := net.SplitHostPort(address); err == nil && !strings.ContainsAny(port, "/?#") && (strings.HasPrefix(address, "[") || !strings.ContainsAny(host, "/?#+")) {
		return "tcp", nil
	}
	scheme, _, explicit := strings.Cut(address, ":")
	if !explicit {
		return "tcp", nil
	}
	parsed, err := url.Parse(scheme + ":")
	if err != nil {
		return "", fmt.Errorf("parse Anki client scheme: %w", err)
	}
	if parsed.Scheme == "" {
		return "", fmt.Errorf("Anki client address has an invalid scheme")
	}
	return parsed.Scheme, nil
}

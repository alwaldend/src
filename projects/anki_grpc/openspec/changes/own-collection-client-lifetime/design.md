## Context

Concrete process, network, native and WASM clients already implement startup,
connection access and shutdown. The generic client wrapper added a second
lifecycle solely to manage reusable collection phases. Production collection
lifecycle usage is confined to Anki as Code's Store; MCP never reopens collection
phases on the same transport.

## Goals / Non-Goals

Expose connection access and shutdown through Client; StartableClient embeds
Client and adds Start for callers that own startup. Keep collection storage
policy with the application and retain actual transport termination
handling. Do not introduce a replacement generic wrapper, dependency, protocol
or archive format.

## Decisions

Use Client for Connection and Stop and StartableClient for all three lifecycle
methods. Factories return StartableClient directly. Keeping the expanded wrapper
would retain duplicate lifecycle state; restoring the old
nested collection-client implementation would retain two client owners. Neither
is needed by the observed callers.

ClientFactory is a constructor function returning StartableClient. The concrete
Factory retains registration and default-address behavior; consumers inject its
NewClient method value. ToClientFactory adapts typed constructors with
contextual error wrapping and maps a nil constructor to nil. WithClient retains
default selection for a nil factory; registration rejects nil factories.

Concrete owned transports retain Done for native completion and process-reaping
evidence. The application collection package observes it through a private
Done-only capability. No public AsyncClient interface is needed. Raw network
clients close their connection; server and collection lifetime remain caller-owned.

Anki as Code returns an already-open Client with Connection and Stop, sharing
the generic capability interface instead of declaring a separate Collection
interface.
Its private single-use scope owns the native open, temporary restore, optional
publication, active-call drain and one terminal cleanup attempt. Preparation
finishes before the handle is exposed. Remove reusable phases, a second startup
state, generic transport wrapping and unused restore-directory configuration.
Keep completion private to source reservations rather than expanding the public
handle. Native consuming RPCs remain observable where application cleanup needs
to avoid closing consumed storage or reopen it for enabled write-back.

The native Rust server retains synchronous backend ownership until its serving
runtime has finished. This independent shutdown correction remains necessary
for successful media operations followed by normal stdio exit.

The public-surface audit covers every handwritten exported symbol, not only Go
interface declarations. Native C bridge types, functions, fields and errors
used only inside their Go package remain private. Rust CLI parsing and help
declarations used only inside their module remain private as well. Preserve
exports required by actual consumers, the C ABI or framework integration.
Visibility reductions do not change native execution or the wire protocol.

## Risks / Trade-offs

Cancellation does not establish native completion. The application must retain
its owned restores and source reservations when completion cannot be established.
This responsibility does not belong to a raw network client's Stop method.
Transport fixtures now use generated native open/close RPCs explicitly; package
export consumes native storage, so cleanup must not blindly close it again.

## Migration Plan

Adapt existing real-worker cases before production edits. Return factory results
directly, migrate application API/MCP handles and fixture helpers, and place
managed archive acceptance tests with the application scope. Validate the store,
transport, native/WASM application and MCP workflows with retained artifacts.

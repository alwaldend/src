// Package ankiascodeapi defines shared Go contracts for Anki as code callers.
// App describes configuration, collection, state, generic API and server workflows without
// importing their concrete implementation or native worker.
// Server contexts control startup; the injected runner owns native hosting.
// Collection sources describe either an application-owned path or a borrowed,
// already-open gRPC connection. Creating a source performs no I/O and owns no
// native resources. The concrete application owns workflow execution.
package ankiascodeapi

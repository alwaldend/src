// Package ankiascodeapi defines shared Go contracts for Anki as code callers.
// App describes configuration, collection, state, generic API and server workflows
// without importing their concrete implementation or native worker.
// Server contexts control startup; the injected runner owns native hosting.
// CollectionOptions combines managed path, external address, borrowed connection
// and namespace selection. Borrowed connections retain caller-owned lifetime and
// publication. Creating options performs no I/O; App validates access combinations
// and owns workflow execution.
package ankiascodeapi

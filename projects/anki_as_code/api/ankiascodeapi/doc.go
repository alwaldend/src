// Package ankiascodeapi defines shared Go contracts for Anki as code callers.
// App describes configuration, collection, state, generic API and server workflows
// without importing their concrete implementation or native worker.
// Server contexts control startup; the injected runner owns native hosting.
// MCP is a separate hosting contract injected into CLI beside App.
// CollectionOptions combines managed path, backend address, borrowed connection
// and namespace selection. Borrowed connections retain caller-owned lifetime and
// publication. Creating options performs no I/O; App validates access combinations
// and owns workflow execution.
// ColOpenOptions selects path, address and archive write-back for an explicit session.
// ColOpen returns the existing ankiclient.Client contract, already started;
// callers borrow Connection and Stop the owner to publish and release resources.
package ankiascodeapi

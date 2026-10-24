import "/wasm_exec.js";
import { createBrowserServer } from "/wasi.js";
import { installGoFilesystem } from "/go.js";
import { runStorage } from "/storage_runtime.js";

// The test runner must supply a caller-owned provider module. No default
// browser backing store is selected or implemented by the reusable backend.
try {
    const { fs, root } = await import("/filesystem.js");
    const evidence = await runStorage({
        fs,
        root,
        module: await WebAssembly.compileStreaming(fetch("/guest.wasm")),
        application: await (await fetch("/application.wasm")).arrayBuffer(),
        createBrowserServer,
        installGoFilesystem,
    });
    postMessage(evidence);
} catch (error) {
    postMessage({ error: error.stack || error.message });
}

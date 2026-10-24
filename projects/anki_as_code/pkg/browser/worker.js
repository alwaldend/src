import "./wasm_exec.js";
import { loadWorkspace, saveWorkspace } from "./storage.js";

let initialized = false;
let ready = false;
let failed = false;
let requests = Promise.resolve();

function message(error) {
    return error instanceof Error ? error.message : String(error);
}

function fail(error) {
    if (failed) return;
    failed = true;
    postMessage({ type: "fatal", error: message(error) });
}

async function initialize(wasmURL, serverWasmURL) {
    let serverModule;
    globalThis.ankiascodeCreateServer = async (filesystem) => {
        serverModule ??= (async () => {
            const response = await fetch(serverWasmURL);
            if (!response.ok) {
                throw new Error(
                    `Fetch server WebAssembly: HTTP ${response.status}`,
                );
            }
            return WebAssembly.compile(await response.arrayBuffer());
        })().catch((error) => {
            serverModule = undefined;
            throw error;
        });
        const [{ createBrowserServer }, module] = await Promise.all([
            import("./wasi.js"),
            serverModule,
        ]);
        return createBrowserServer(module, filesystem);
    };
    globalThis.ankiascodeReady = () => {
        if (failed) return;
        ready = true;
        postMessage({ type: "ready" });
    };
    const response = await fetch(wasmURL);
    if (!response.ok) {
        throw new Error(
            `Fetch application WebAssembly: HTTP ${response.status}`,
        );
    }
    const go = new globalThis.Go();
    // ArrayBuffer instantiation supports servers without a WASM MIME type.
    const { instance } = await WebAssembly.instantiate(
        await response.arrayBuffer(),
        go.importObject,
    );
    // Go.run lasts until exit; Go signals readiness separately.
    void go
        .run(instance)
        .then(() => fail(new Error("Go browser application exited")), fail);
}

addEventListener("messageerror", () => {
    fail(new Error("Cannot decode browser application request"));
});

addEventListener("message", ({ data }) => {
    if (failed) return;
    if (!data || typeof data !== "object") {
        fail(new Error("Invalid browser application request"));
        return;
    }
    if (data.type === "init") {
        if (initialized) {
            fail(new Error("Browser application is already initialized"));
            return;
        }
        initialized = true;
        void initialize(data.wasmURL, data.serverWasmURL).catch(fail);
        return;
    }
    if (data.type !== "request") {
        fail(new Error("Unknown browser application request"));
        return;
    }
    // One queue covers Go calls, scoped server cleanup and IndexedDB transactions.
    // A rejected request is reported without preventing subsequent requests.
    requests = requests
        .then(async () => {
            if (failed) return;
            try {
                if (!ready)
                    throw new Error("Browser application is not ready");
                const value = await dispatch(data.method, data.args);
                postMessage({ type: "result", id: data.id, value });
            } catch (error) {
                postMessage({
                    type: "result",
                    id: data.id,
                    error: message(error),
                });
            }
        })
        .catch(fail);
});

async function invoke(method, args) {
    const result = await globalThis.ankiascodeInvoke(method, args);
    if (result.error !== undefined) throw new Error(result.error);
    return result.value;
}

async function dispatch(method, args) {
    if (method === "saveSession") {
        const entries = await invoke("snapshot", {});
        return saveWorkspace(args.key, args.expectedRevision, entries);
    }
    if (method === "restoreSession") {
        const snapshot = await loadWorkspace(args.key);
        if (snapshot === null) return null;
        await invoke("restoreSnapshot", { entries: snapshot.entries });
        return { revision: snapshot.revision, savedAt: snapshot.savedAt };
    }
    return invoke(method, args);
}

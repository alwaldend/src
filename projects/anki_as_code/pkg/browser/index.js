// Run the Go application and its memory workspace in a dedicated worker.
export async function createAnkiApp(options = {}) {
    requireObject(options, "options");
    const wasmURL = resourceURL(options.wasmURL, "app.wasm");
    const workerURL = resourceURL(options.workerURL, "worker.js");
    const serverWasmURL = resourceURL(options.serverWasmURL, "server.wasm");
    const worker = new Worker(workerURL, { type: "module" });
    const pending = new Map();
    let nextID = 0;
    let failure;
    let resolveReady;
    let rejectReady;
    const ready = new Promise((resolve, reject) => {
        resolveReady = resolve;
        rejectReady = reject;
    });
    const startupTimer = setTimeout(() => {
        stop(new Error("Browser application startup timed out after 30s"));
    }, 30000);

    function stop(error) {
        if (failure) return;
        failure = error;
        clearTimeout(startupTimer);
        worker.terminate();
        rejectReady(error);
        for (const request of pending.values()) request.reject(error);
        pending.clear();
    }

    worker.addEventListener("error", (event) => {
        event.preventDefault();
        stop(new Error(event.message || "Browser application worker failed"));
    });
    worker.addEventListener("messageerror", () => {
        stop(new Error("Cannot decode browser application worker response"));
    });
    worker.addEventListener("message", ({ data }) => {
        if (failure) return;
        if (!data || typeof data !== "object") {
            stop(new Error("Invalid browser application worker response"));
            return;
        }
        switch (data.type) {
            case "ready":
                clearTimeout(startupTimer);
                resolveReady();
                return;
            case "fatal":
                stop(new Error(data.error || "Browser application exited"));
                return;
            case "result": {
                const request = pending.get(data.id);
                if (!request) {
                    stop(new Error("Unexpected browser application result"));
                    return;
                }
                pending.delete(data.id);
                if (data.error !== undefined) {
                    request.reject(new Error(data.error));
                } else {
                    request.resolve(data.value);
                }
                return;
            }
            default:
                stop(new Error("Unknown browser application worker response"));
        }
    });

    function request(method, args) {
        if (failure) return Promise.reject(failure);
        const id = ++nextID;
        return new Promise((resolve, reject) => {
            pending.set(id, { resolve, reject });
            try {
                // Clone bytes without detaching the caller's input buffer.
                worker.postMessage({ type: "request", id, method, args });
            } catch (error) {
                pending.delete(id);
                reject(error);
            }
        });
    }

    try {
        worker.postMessage({ type: "init", wasmURL, serverWasmURL });
    } catch (error) {
        stop(error);
    }
    await ready;

    return {
        app: {
            async configExport(options) {
                requireObject(options, "configExport options");
                requirePath(options.output);
                requireObject(options.configs, "configs");
                const configs = JSON.stringify(options.configs);
                if (typeof configs !== "string") {
                    throw new TypeError("configs must serialize to JSON");
                }
                await request("configExport", {
                    output: options.output,
                    configs,
                });
            },
            async configDump(options) {
                requireObject(options, "configDump options");
                requirePath(options.text);
                const { format = "json" } = options;
                if (typeof format !== "string") {
                    throw new TypeError("format must be a string");
                }
                return request("configDump", { text: options.text, format });
            },
            async collectionExport(options) {
                const args = collectionOptions(options);
                requirePath(options.output);
                const stateOutput = options.stateOutput ?? "";
                if (stateOutput !== "") requirePath(stateOutput);
                await request("collectionExport", {
                    ...args,
                    output: options.output,
                    stateOutput,
                });
            },
            async collectionStateUpdate(options) {
                const args = collectionOptions(options);
                requirePath(options.state);
                await request("collectionStateUpdate", {
                    ...args,
                    state: options.state,
                });
            },
            async collectionPlan(options) {
                const args = collectionOptions(options);
                requirePath(options.text);
                return JSON.parse(
                    await request("collectionPlan", {
                        ...args,
                        text: options.text,
                    }),
                );
            },
            async collectionApply(options) {
                const args = collectionOptions(options);
                requirePath(options.text);
                await request("collectionApply", {
                    ...args,
                    text: options.text,
                });
            },
            async collectionBuild(options) {
                const args = collectionOptions(options, "base");
                requirePath(options.text);
                requirePackage(options.output);
                await request("collectionBuild", {
                    ...args,
                    text: options.text,
                    output: options.output,
                });
            },
        },
        session: {
            async save(key, expectedRevision = 0) {
                requireSessionKey(key);
                if (
                    !Number.isSafeInteger(expectedRevision) ||
                    expectedRevision < 0
                ) {
                    throw new TypeError(
                        "expectedRevision must be a nonnegative safe integer",
                    );
                }
                return request("saveSession", { key, expectedRevision });
            },
            async restore(key) {
                requireSessionKey(key);
                return request("restoreSession", { key });
            },
        },
        workspace: {
            async writeFile(path, data) {
                requirePath(path);
                if (!(data instanceof Uint8Array)) {
                    throw new TypeError("data must be a Uint8Array");
                }
                await request("writeFile", { path, data });
            },
            async readFile(path) {
                requirePath(path);
                return request("readFile", { path });
            },
            async listFiles(path) {
                requirePath(path);
                return request("listFiles", { path });
            },
            async remove(path) {
                requirePath(path);
                await request("remove", { path });
            },
            async rename(source, destination) {
                requirePath(source);
                requirePath(destination);
                await request("rename", { source, destination });
            },
        },
        async dispose() {
            stop(new Error("Browser application has been disposed"));
        },
    };
}

function collectionOptions(options, source = "collection") {
    requireObject(options, "collection options");
    requirePackage(options[source]);
    const stateKey = options.stateKey ?? "";
    if (typeof stateKey !== "string") {
        throw new TypeError("stateKey must be a string");
    }
    return { [source]: options[source], stateKey };
}

function requirePackage(path) {
    requirePath(path);
    if (!path.endsWith(".colpkg")) {
        throw new TypeError("collection paths must identify .colpkg packages");
    }
}

function requireSessionKey(key) {
    if (typeof key !== "string" || key.length === 0) {
        throw new TypeError("session key must be a nonempty string");
    }
}

function requireObject(value, name) {
    if (!value || typeof value !== "object" || Array.isArray(value)) {
        throw new TypeError(`${name} must be an object`);
    }
}

function requirePath(path) {
    if (typeof path !== "string" || !path.startsWith("/")) {
        throw new TypeError("path must be an absolute virtual path");
    }
    if (path.includes("\0")) {
        throw new TypeError("path must not contain a null character");
    }
}

function resourceURL(value, fallback) {
    if (value === undefined) return new URL(fallback, import.meta.url).href;
    if (typeof value !== "string" && !(value instanceof URL)) {
        throw new TypeError("resource URLs must be strings or URL objects");
    }
    return new URL(value, import.meta.url).href;
}

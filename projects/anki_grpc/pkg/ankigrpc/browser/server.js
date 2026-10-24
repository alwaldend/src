/** Wrap an initialized WASI reactor with an owned, synchronous command API. */
export function createServer(instance) {
    let exports = instance.exports;
    let memory = exports.memory;
    for (const name of [
        "anki_backend_new",
        "anki_backend_command",
        "anki_backend_free",
        "anki_browser_buffer_new",
        "anki_browser_buffer_data",
        "anki_browser_buffer_len",
        "anki_browser_buffer_free",
    ]) {
        if (typeof exports[name] !== "function") {
            throw new TypeError(`Missing Anki backend export: ${name}`);
        }
    }
    if (!(memory instanceof WebAssembly.Memory)) {
        throw new TypeError("Missing Anki backend memory");
    }
    const encoder = new TextEncoder();
    let handle = 0;
    let failure;
    let closed = false;

    // Once an export traps, Rust state may be inconsistent. Drop the instance
    // without entering its allocator or destructors again.
    function call(name, ...args) {
        try {
            return exports[name](...args);
        } catch (error) {
            failure = error;
            throw error;
        }
    }

    function withBuffers(operation) {
        const buffers = [];
        function allocate(bytes = new Uint8Array()) {
            const buffer = call("anki_browser_buffer_new", bytes.length);
            buffers.push(buffer);
            const pointer = call("anki_browser_buffer_data", buffer) >>> 0;
            new Uint8Array(memory.buffer, pointer, bytes.length).set(bytes);
            return { buffer, pointer, length: bytes.length };
        }
        function read(buffer) {
            const pointer = call("anki_browser_buffer_data", buffer) >>> 0;
            const length = call("anki_browser_buffer_len", buffer) >>> 0;
            // Each export may grow memory. Acquire a new view, then copy before free.
            return new Uint8Array(memory.buffer, pointer, length).slice();
        }
        try {
            return operation(allocate, read);
        } catch (error) {
            failure = error;
            throw error;
        } finally {
            if (!failure) {
                for (const buffer of buffers.reverse()) {
                    call("anki_browser_buffer_free", buffer);
                }
            }
        }
    }

    const initializationError = withBuffers((allocate, read) => {
        const output = allocate();
        handle = call("anki_backend_new", output.buffer);
        return handle
            ? undefined
            : new TextDecoder().decode(read(output.buffer));
    });
    if (!handle) {
        throw new Error(
            initializationError || "Initialize Anki backend failed",
        );
    }

    return {
        command(method, request) {
            if (closed) throw new Error("Anki backend is closed");
            if (failure) throw failure;
            if (
                typeof method !== "string" ||
                !(request instanceof Uint8Array)
            ) {
                throw new TypeError(
                    "Anki command requires a method and Uint8Array",
                );
            }
            return withBuffers((allocate, read) => {
                const path = allocate(encoder.encode(method));
                const input = allocate(request);
                const output = allocate();
                const code = call(
                    "anki_backend_command",
                    handle,
                    path.pointer,
                    path.length,
                    input.pointer,
                    input.length,
                    output.buffer,
                );
                if (code !== 0 && code !== 1 && code !== 2) {
                    throw new Error(`Invalid Anki backend status: ${code}`);
                }
                return { code, data: read(output.buffer) };
            });
        },
        close() {
            if (closed) return;
            closed = true;
            try {
                if (!failure) call("anki_backend_free", handle);
            } finally {
                handle = 0;
                exports = undefined;
                memory = undefined;
            }
        },
    };
}

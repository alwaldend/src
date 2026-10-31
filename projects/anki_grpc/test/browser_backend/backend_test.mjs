import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import fs from "node:fs/promises";
import path from "node:path";
import { pathToFileURL } from "node:url";
import { WASI } from "node:wasi";

const [applicationPath, backendPath, runtimePath, linkerPath] = process.argv
    .slice(2)
    .map((filename) => path.resolve(filename));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const root = await fs.mkdtemp(path.join(process.env.TEST_TMPDIR, "backend-"));
const report = {
    node: process.version,
    host: "Node built-in WASI reactor with UTC fixture timezone",
    scope: "Generated Go clients call the Rust backend through its JavaScript API; shared storage is verified by the independent production-host E2Es",
    phase: "initializing",
    cases: [],
    backends: [],
};
let goRun;
let goExited = false;
let goFinished = false;
let failure;
let deadline;
const timeout = new Promise((_, reject) => {
    deadline = setTimeout(
        () => reject(new Error("Browser backend E2E timed out")),
        90000,
    );
});

try {
    const [applicationBytes, backendBytes] = await Promise.all([
        fs.readFile(applicationPath),
        fs.readFile(backendPath),
    ]);
    report.applicationSHA256 = createHash("sha256")
        .update(applicationBytes)
        .digest("hex");
    report.backendSHA256 = createHash("sha256")
        .update(backendBytes)
        .digest("hex");
    await import(pathToFileURL(runtimePath).href);
    const { createServer } = await import(pathToFileURL(linkerPath).href);
    const backendModule = await WebAssembly.compile(backendBytes);
    report.backendImports = WebAssembly.Module.imports(backendModule);
    report.backendExports = WebAssembly.Module.exports(backendModule);
    assert.ok(
        report.backendImports.every(({ module }) => module !== "anki_io"),
    );
    assert.ok(!process.execArgv.includes("--experimental-wasm-jspi"));
    assert.ok(
        report.backendExports.some(({ name }) => name === "_initialize"),
    );
    assert.ok(!report.backendExports.some(({ name }) => name === "_start"));
    report.cases.push(
        "reactor runs without pipe imports, CLI entry point or JSPI flags",
    );

    // Each factory initializes the actual Rust reactor, retaining separate ownership.
    async function createBackend(mode) {
        const wasi = new WASI({
            version: "preview1",
            args: [],
            env: { TMPDIR: "/workspace", TZ: "UTC" },
            preopens: { "/workspace": root },
            returnOnExit: true,
        });
        const instance = await WebAssembly.instantiate(backendModule, {
            wasi_snapshot_preview1: wasi.wasiImport,
            anki_timezone: {
                offset_for_utc: () => 0,
                offsets_for_local: () => 0n,
            },
        });
        wasi.initialize(instance);
        const counts = {
            mode,
            commands: 0,
            closes: 0,
            guestTraps: 0,
            guestEntries: 0,
            bytesRequested: 0,
            bytesReturned: 0,
        };
        report.backends.push(counts);
        const samples = [];
        // An invalid descriptor forces a real guest memory trap inside command execution.
        const callable =
            mode === "trap"
                ? {
                      exports: {
                          ...Object.fromEntries(
                              Object.entries(instance.exports).map(
                                  ([name, value]) => [
                                      name,
                                      typeof value === "function"
                                          ? (...args) => {
                                                counts.guestEntries++;
                                                return value(...args);
                                            }
                                          : value,
                                  ],
                              ),
                          ),
                          anki_backend_command() {
                              counts.guestEntries++;
                              counts.guestTraps++;
                              return instance.exports.anki_browser_buffer_data(
                                  0xfffffff8,
                              );
                          },
                      },
                  }
                : instance;
        const server = createServer(callable);
        if (mode === "normal") {
            globalThis.ankiBrowserBackendMemoryBytes = () =>
                instance.exports.memory.buffer.byteLength;
        }
        function checkOwnedBuffers() {
            for (const { data, expected } of samples) {
                assert.deepEqual(
                    data,
                    expected,
                    "later calls changed an earlier response buffer",
                );
            }
        }
        return {
            command(method, bytes) {
                counts.commands++;
                counts.bytesRequested += bytes.length;
                checkOwnedBuffers();
                const request = bytes.slice();
                const result = server.command(method, bytes);
                assert.deepEqual(
                    bytes,
                    request,
                    "command changed its caller-owned request",
                );
                assert.ok(result.data instanceof Uint8Array);
                assert.notEqual(
                    result.data.buffer,
                    instance.exports.memory.buffer,
                    "response borrows live Rust memory",
                );
                counts.bytesReturned += result.data.length;
                checkOwnedBuffers();
                if (samples.length < 3 && result.data.length > 0) {
                    samples.push({
                        data: result.data,
                        expected: result.data.slice(),
                    });
                }
                return result;
            },
            close() {
                counts.closes++;
                checkOwnedBuffers();
                server.close();
                checkOwnedBuffers();
            },
        };
    }

    globalThis.ankiBrowserBackendFactory = (_timezone) =>
        createBackend("normal");
    globalThis.ankiBrowserBackendFailedFactory = async (_timezone) => {
        throw new Error("fixture startup failure");
    };
    let releaseFactory;
    const factoryGate = new Promise((resolve) => {
        releaseFactory = resolve;
    });
    globalThis.ankiBrowserBackendDelayedFactory = async (_timezone) => {
        globalThis.ankiBrowserBackendFactoryStarted();
        await factoryGate;
        return createBackend("delayed");
    };
    globalThis.ankiBrowserBackendReleaseFactory = () => releaseFactory();
    globalThis.ankiBrowserBackendDelayedCloses = () =>
        report.backends.find(({ mode }) => mode === "delayed")?.closes ?? 0;
    globalThis.ankiBrowserBackendTrapFactory = (_timezone) =>
        createBackend("trap");
    const verified = new Promise((resolve) => {
        globalThis.ankiBrowserBackendVerified = () => {
            report.phase = "client workflow completed";
            report.workflow = JSON.parse(globalThis.ankiBrowserBackendReport);
            resolve();
        };
    });
    const go = new globalThis.Go();
    let exitCode;
    go.exit = (code) => {
        exitCode = code;
    };
    const { instance: application } = await WebAssembly.instantiate(
        applicationBytes,
        go.importObject,
    );
    report.phase = "client workflow running";
    goRun = go.run(application).finally(() => {
        goExited = true;
    });
    const earlyExit = goRun.then(() => {
        throw new Error("Go exited before completing the browser workflow");
    });
    await Promise.race([verified, earlyExit, timeout]);
    assert.equal(report.workflow.error, undefined, report.workflow.error);
    assert.equal(report.workflow.repeatedCalls, 16);
    assert.deepEqual(
        report.backends.map(({ mode }) => mode),
        ["normal", "delayed", "trap"],
    );
    assert.ok(
        report.backends.every(({ closes }) => closes === 1),
        "a created backend was not closed exactly once",
    );
    const normal = report.backends.find(({ mode }) => mode === "normal");
    assert.ok(normal.bytesRequested > 16 * report.workflow.payloadBytes);
    assert.ok(normal.bytesReturned > 16 * report.workflow.payloadBytes);
    assert.equal(normal.closes, 1);
    assert.equal(
        report.backends.find(({ mode }) => mode === "delayed").closes,
        1,
    );
    assert.equal(
        report.backends.find(({ mode }) => mode === "trap").commands,
        1,
    );
    assert.equal(
        report.backends.find(({ mode }) => mode === "trap").guestTraps,
        1,
    );
    assert.equal(
        report.backends.find(({ mode }) => mode === "trap").closes,
        1,
    );
    report.cases.push(
        "caller-owned buffers survive later commands and backend cleanup",
    );
    report.cases.push(
        "each normal, canceled and trapped backend closes exactly once",
    );

    // Exercise the JS wrapper directly too: poisoning must reject a second call before entering Wasm.
    const trapped = await createBackend("trap");
    assert.throws(
        () =>
            trapped.command(
                "/anki.card_rendering.BackendCardRenderingService/StripHtml",
                new Uint8Array(),
            ),
        WebAssembly.RuntimeError,
    );
    const guestEntriesAfterTrap = report.backends.at(-1).guestEntries;
    assert.throws(() =>
        trapped.command(
            "/anki.card_rendering.BackendCardRenderingService/StripHtml",
            new Uint8Array(),
        ),
    );
    assert.equal(
        report.backends.at(-1).guestTraps,
        1,
        "poisoned handle entered Wasm a second time",
    );
    assert.equal(
        report.backends.at(-1).guestEntries,
        guestEntriesAfterTrap,
        "poisoned handle called another Wasm export",
    );
    try {
        trapped.close();
    } catch (error) {
        assert.ok(error instanceof Error);
    }
    report.cases.push(
        "actual guest memory trap permanently poisons the JavaScript handle",
    );

    globalThis.ankiBrowserBackendFinish();
    goFinished = true;
    await Promise.race([goRun, timeout]);
    assert.equal(exitCode, 0);
    report.phase = "complete";
} catch (error) {
    failure = error;
    report.error = { name: error.name, message: error.message };
} finally {
    clearTimeout(deadline);
    globalThis.ankiBrowserBackendReleaseFactory?.();
    if (report.workflow && goRun && !goExited && !goFinished) {
        globalThis.ankiBrowserBackendFinish();
        await goRun;
    }
    await fs.writeFile(
        path.join(output, "browser-backend.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
    await fs.rm(root, { recursive: true, force: true });
}

// Go's runtime may retain timers after main exits; all checks and artifacts are complete.
if (failure) console.error(failure);
process.exit(failure ? 1 : 0);

import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import fs from "node:fs/promises";
import http from "node:http";
import { createRequire } from "node:module";
import path from "node:path";

// Reuse the repository's pinned browser tools without an application SDK.
const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const [
    application,
    guest,
    wasmExec,
    wasi,
    bridge,
    worker,
    executablePath,
    goHost,
    runtime,
    filesystem,
    emscripten,
    memfs,
] = process.argv.slice(2).map((file) => path.resolve(file));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must supply the retained artifact directory");
const report = { cases: [], artifacts: {}, requests: [], console: [] };
const resources = new Map([
    ["/application.wasm", application],
    ["/guest.wasm", guest],
    ["/wasm_exec.js", wasmExec],
    ["/wasi.js", wasi],
    ["/server.js", bridge],
    ["/storage_worker.js", worker],
    ["/go.js", goHost],
    ["/storage_runtime.js", runtime],
]);
let browser;
let server;
let deadline;
try {
    assert.ok(
        filesystem,
        "Chromium acceptance requires a caller-supplied filesystem module exporting fs and root",
    );
    resources.set("/filesystem.js", filesystem);
    if (emscripten) resources.set("/emscripten.js", emscripten);
    if (memfs) resources.set("/emscripten_memfs.js", memfs);
    for (const [url, file] of resources) {
        const bytes = await fs.readFile(file);
        report.artifacts[url] = {
            bytes: bytes.length,
            sha256: createHash("sha256").update(bytes).digest("hex"),
        };
        resources.set(url, bytes);
    }
    const compiled = await WebAssembly.compile(resources.get("/guest.wasm"));
    report.imports = WebAssembly.Module.imports(compiled);
    report.exports = WebAssembly.Module.exports(compiled);
    assert.ok(report.imports.some((item) => item.module === "anki_timezone"));
    assert.ok(report.exports.some((item) => item.name === "_initialize"));
    assert.ok(!report.exports.some((item) => item.name === "_start"));
    server = http.createServer((request, response) => {
        const url = new URL(request.url, "http://localhost").pathname;
        report.requests.push(url);
        if (url === "/") {
            response.setHeader("Content-Type", "text/html");
            response.end(
                "<!doctype html><title>Reusable Anki browser backend</title>",
            );
        } else if (resources.has(url)) {
            response.setHeader(
                "Content-Type",
                url.endsWith(".wasm") ? "application/wasm" : "text/javascript",
            );
            response.end(resources.get(url));
        } else {
            response.writeHead(404).end();
        }
    });
    await new Promise((resolve, reject) => {
        server.once("error", reject);
        server.listen(0, "127.0.0.1", resolve);
    });
    browser = await puppeteer.launch({
        executablePath,
        headless: "shell",
        userDataDir: await fs.mkdtemp(
            path.join(process.env.TEST_TMPDIR, "backend-chromium-"),
        ),
        env: { ...process.env, TZ: "America/New_York" },
    });
    report.browser = await browser.version();
    const page = await browser.newPage();
    const errors = [];
    page.on("pageerror", (error) => errors.push(error.message));
    page.on("console", (message) =>
        report.console.push({ type: message.type(), text: message.text() }),
    );
    await page.goto(`http://127.0.0.1:${server.address().port}/`);
    const executionResult = page.evaluate(
        () =>
            new Promise((resolve, reject) => {
                const worker = new Worker("/storage_worker.js", {
                    type: "module",
                });
                const timeout = setTimeout(() => {
                    worker.terminate();
                    reject(
                        new Error(
                            "browser storage workflow timed out after 120 seconds",
                        ),
                    );
                }, 120000);
                worker.onmessage = ({ data }) => {
                    clearTimeout(timeout);
                    worker.terminate();
                    resolve(data);
                };
                worker.onerror = (error) => {
                    clearTimeout(timeout);
                    worker.terminate();
                    reject(
                        new Error(`browser storage worker: ${error.message}`),
                    );
                };
            }),
    );
    report.execution = await Promise.race([
        executionResult,
        new Promise((_, reject) => {
            deadline = setTimeout(
                () =>
                    reject(
                        new Error(
                            "Chromium workflow exceeded its 150-second Node deadline",
                        ),
                    ),
                150000,
            );
        }),
    ]);
    const execution = report.execution;
    assert.equal(execution.error, undefined, execution.error);
    assert.equal(execution.exitCode, 0, "Go harness exit code");
    assert.equal(
        execution.workflow.error,
        undefined,
        execution.workflow.error,
    );
    assert.equal(execution.timezone, "America/New_York");
    assert.equal(
        execution.instances,
        2,
        "independent Rust instances must share the caller-owned JavaScript store",
    );
    assert.deepEqual(execution.timezoneOffsets, [
        [-18000, -14400],
        [-18000, -14400],
    ]);
    for (const operation of [
        "openSync",
        "readSync",
        "writeSync",
        "closeSync",
        "readdirSync",
        "timezoneUTC",
    ]) {
        assert.ok(
            execution.hostCalls[operation] > 0,
            `production Rust host did not call ${operation}`,
        );
    }
    assert.equal(execution.remainingDescriptors, 0);
    assert.ok(execution.workflow.storage.noteId > 0);
    assert.ok(execution.workflow.storage.collectionBytes > 4096);
    const media = Buffer.alloc(7 * 16384);
    for (let offset = 0; offset < media.length; offset += 7)
        media.set([0, 255, 128, 13, 10, 1, 254], offset);
    assert.equal(
        execution.workflow.storage.mediaSHA256,
        createHash("sha256").update(media).digest("hex"),
    );
    assert.deepEqual(errors, []);
    report.cases.push(
        ...execution.workflow.cases,
        "named timezone supplies winter/summer offsets and real Rust callbacks",
        "production WASI host runs in Chromium without JSPI",
    );
} catch (error) {
    report.error = error.stack || error.message;
    process.exitCode = 1;
} finally {
    clearTimeout(deadline);
    try {
        if (browser) await browser.close();
    } finally {
        if (server?.listening)
            await new Promise((resolve) => server.close(resolve));
        await fs.writeFile(
            path.join(output, "browser-storage.json"),
            `${JSON.stringify(report, null, 2)}\n`,
        );
    }
}
assert.equal(report.error, undefined, report.error);

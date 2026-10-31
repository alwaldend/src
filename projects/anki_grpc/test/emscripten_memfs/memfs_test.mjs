import assert from "node:assert/strict";
import crypto from "node:crypto";
import fs from "node:fs/promises";
import path from "node:path";
import vm from "node:vm";
import { isMainThread, parentPort, Worker } from "node:worker_threads";
import zlib from "node:zlib";

async function exercise() {
    const artifact = await fs.readFile(process.argv[2]);
    const context = vm.createContext({
        TextDecoder,
        TextEncoder,
        URL,
        console,
    });
    vm.runInContext(
        `globalThis.self = globalThis;
   globalThis.WorkerGlobalScope = class WorkerGlobalScope {};
   globalThis.location = new URL("https://memfs.invalid/worker.js");
   globalThis.WebAssembly = undefined;`,
        context,
    );
    const absentGlobals = [
        "process",
        "Buffer",
        "require",
        "global",
        "window",
        "document",
        "fetch",
        "WebAssembly",
    ];
    const assertAbsentGlobals = () => {
        for (const name of absentGlobals) {
            assert.equal(
                vm.runInContext(`typeof ${name}`, context),
                "undefined",
                name,
            );
        }
    };
    assertAbsentGlobals();
    const module = new vm.SourceTextModule(artifact.toString("utf8"), {
        context,
        initializeImportMeta(meta) {
            meta.url = "https://memfs.invalid/emscripten_memfs.js";
        },
    });
    await module.link((specifier) => {
        throw new Error(`Unexpected runtime import: ${specifier}`);
    });
    await module.evaluate({ timeout: 10000 });
    assert.equal(typeof module.namespace.default, "function");

    const firstPending = module.namespace.default();
    assert.equal(
        typeof firstPending.then,
        "function",
        "factory must be async",
    );
    const [first, second] = await Promise.all([
        firstPending,
        module.namespace.default({}),
    ]);
    assertAbsentGlobals();
    assert.notEqual(first, second);
    const { FS, ERRNO_CODES } = first;
    assert.notEqual(FS, second.FS);
    assert.equal(typeof ERRNO_CODES.ENOENT, "number");
    const primitives = [
        "getStream",
        "getStreamChecked",
        "stat",
        "fstat",
        "fchmod",
        "fchown",
        "utime",
        "ftruncate",
        "read",
        "write",
        "open",
        "close",
        "readdir",
        "rename",
        "symlink",
        "readlink",
    ];
    for (const primitive of primitives) {
        assert.equal(typeof FS[primitive], "function", primitive);
    }
    for (const fd of [0, 1, 2]) assert.equal(FS.getStreamChecked(fd).fd, fd);
    assert.throws(
        () => FS.stat("/missing"),
        (e) => e.errno === ERRNO_CODES.ENOENT,
    );
    assert.throws(
        () => FS.open("/missing/child", "w"),
        (e) => e.errno === ERRNO_CODES.ENOENT,
    );
    assert.throws(
        () => FS.getStreamChecked(9999),
        (e) => e.errno === ERRNO_CODES.EBADF,
    );

    const directory = "/資料-é";
    const original = `${directory}/данные.bin`;
    const renamed = `${directory}/改名.bin`;
    const bytes = Uint8Array.from(
        { length: 65539 },
        (_, index) => index % 256,
    );
    const expected = bytes.slice();
    FS.mkdir(directory);
    FS.chdir(directory);
    assert.equal(second.FS.cwd(), "/");
    assert.throws(
        () => second.FS.stat(directory),
        (e) => e.errno === ERRNO_CODES.ENOENT,
    );
    const stream = FS.open(original, "w+");
    assert.equal(stream.fd, 3, "stdio must retain descriptors 0, 1 and 2");
    assert.equal(second.FS.getStream(3), undefined);
    assert.equal(FS.write(stream, bytes, 0, bytes.length), bytes.length);
    assert.equal(stream.position, bytes.length);
    const patch = Uint8Array.of(255, 0, 128, 13, 10);
    assert.equal(FS.write(stream, patch, 0, patch.length, 7), patch.length);
    assert.equal(
        stream.position,
        bytes.length,
        "positioned write retains offset",
    );
    expected.set(patch, 7);
    const actual = new Uint8Array(bytes.length);
    assert.equal(FS.read(stream, actual, 0, actual.length, 0), actual.length);
    assert.deepEqual(actual, expected);
    assert.equal(FS.read(stream, actual, 0, 1, bytes.length), 0);
    FS.rename(original, renamed);
    assert.throws(
        () => FS.stat(original),
        (e) => e.errno === ERRNO_CODES.ENOENT,
    );
    assert.equal(FS.fstat(stream.fd).size, bytes.length);
    FS.ftruncate(stream.fd, 10);
    FS.ftruncate(stream.fd, 15);
    assert.deepEqual(Array.from(FS.readFile(renamed)), [
        ...expected.slice(0, 10),
        0,
        0,
        0,
        0,
        0,
    ]);
    FS.close(stream);
    assert.throws(
        () => FS.getStreamChecked(stream.fd ?? 3),
        (e) => e.errno === ERRNO_CODES.EBADF,
    );

    const text = "Zażółć gęślą jaźń — 日本語 — Привет";
    FS.writeFile("text.txt", text);
    assert.equal(FS.readFile("text.txt", { encoding: "utf8" }), text);
    FS.symlink("text.txt", "link");
    assert.equal(FS.readlink("link"), "text.txt");
    assert.equal(FS.readFile("link", { encoding: "utf8" }), text);
    assert.ok(Array.from(FS.readdir(directory)).includes("改名.bin"));
    second.FS.mkdir(directory);
    second.FS.writeFile(renamed, Uint8Array.of(42));
    assert.equal(FS.stat(renamed).size, 15);
    assert.deepEqual(Array.from(FS.readFile(renamed)), [
        ...expected.slice(0, 10),
        0,
        0,
        0,
        0,
        0,
    ]);
    assert.deepEqual(Array.from(second.FS.readFile(renamed)), [42]);

    return {
        status: "passed",
        artifact: {
            name: path.basename(process.argv[2]),
            sha256: crypto.createHash("sha256").update(artifact).digest("hex"),
            bytes: artifact.length,
            gzipBytes: zlib.gzipSync(artifact).length,
        },
        environment:
            "isolated ES module realm without Node, DOM, fetch or WebAssembly",
        instances: 2,
        binaryBytes: bytes.length,
        verified: [
            "async factory",
            "independent file, descriptor and working-directory state",
            "stdio descriptors",
            "missing paths and invalid descriptors",
            "binary stream reads, writes, offsets and EOF",
            "rename with open descriptor",
            "truncate and zero extension",
            "Unicode paths and text",
            "symlinks and directory entries",
        ],
        primitives,
    };
}

if (!isMainThread) {
    try {
        parentPort.postMessage(await exercise());
    } catch (error) {
        parentPort.postMessage({
            status: "failed",
            error: String(error.stack ?? error),
        });
    }
} else {
    const outputDirectory = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
    assert.ok(outputDirectory, "Bazel must supply the report directory");
    await fs.mkdir(outputDirectory, { recursive: true });
    const worker = new Worker(new URL(import.meta.url), {
        argv: process.argv.slice(2),
    });
    const report = await new Promise((resolve) => {
        const deadline = setTimeout(() => {
            void worker.terminate();
            resolve({
                status: "failed",
                error: "Generated module E2E exceeded 15 seconds",
            });
        }, 15000);
        worker.once("message", (result) => {
            clearTimeout(deadline);
            resolve(result);
        });
        worker.once("error", (error) => {
            clearTimeout(deadline);
            resolve({ status: "failed", error: String(error) });
        });
        worker.once("exit", (code) => {
            clearTimeout(deadline);
            resolve({
                status: "failed",
                error: `Worker exited before report: ${code}`,
            });
        });
    });
    await fs.writeFile(
        path.join(outputDirectory, "emscripten_memfs_report.json"),
        `${JSON.stringify(report, null, 2)}\n`,
    );
    console.log(JSON.stringify(report));
    assert.equal(report.status, "passed", report.error);
}

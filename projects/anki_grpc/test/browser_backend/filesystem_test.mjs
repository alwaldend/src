import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import fs from "node:fs";
import path from "node:path";
import { pathToFileURL } from "node:url";

const [guestPath, hostPath] = process.argv
    .slice(2)
    .map((value) => path.resolve(value));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
const root = fs.mkdtempSync(path.join(process.env.TEST_TMPDIR, "shared-fs-"));
const report = {
    cases: [],
    host: "production WASI with caller-owned node:fs",
};
const owned = new Set();
const closeAttempts = new Map();
let fault;
let beforeOpen;
let commandHook;
const behaviorFailures = [];
report.behaviors = [];
function behavior(name, operation) {
    const evidence = { name };
    report.behaviors.push(evidence);
    try {
        operation(evidence);
        evidence.passed = true;
    } catch (error) {
        evidence.error = { name: error.name, message: error.message };
        behaviorFailures.push(error);
    } finally {
        fault = undefined;
        beforeOpen = undefined;
    }
}
function inject(method, at, code, after = false) {
    fault = {
        method,
        at,
        calls: 0,
        after,
        error: Object.assign(new Error(`fixture ${code}`), { code }),
    };
}
const provider = new Proxy(fs, {
    get(target, name) {
        const value = target[name];
        if (typeof value !== "function") return value;
        return function (...args) {
            assert.equal(
                this,
                provider,
                `provider receiver for ${String(name)}`,
            );
            if (name === "openSync") beforeOpen?.(...args);
            if (name === "closeSync")
                closeAttempts.set(
                    args[0],
                    (closeAttempts.get(args[0]) || 0) + 1,
                );
            let delayedError;
            if (fault?.method === name && ++fault.calls === fault.at) {
                const { error, after } = fault;
                fault = undefined;
                if (!after) throw error;
                delayedError = error;
            }
            const result = Reflect.apply(value, target, args);
            if (name === "openSync") owned.add(result);
            if (name === "closeSync") owned.delete(args[0]);
            if (delayedError) throw delayedError;
            return result;
        };
    },
});
const timezone = { offsetForUTC: () => 0, offsetsForLocal: () => "0" };
let failure;
let server;
const instantiate = WebAssembly.instantiate;
try {
    const guest = fs.readFileSync(guestPath);
    report.guestSHA256 = createHash("sha256").update(guest).digest("hex");
    report.hostSHA256 = createHash("sha256")
        .update(fs.readFileSync(hostPath))
        .digest("hex");
    const module = await WebAssembly.compile(guest);
    const { createBrowserServer } = await import(pathToFileURL(hostPath).href);
    let instance;
    let wasi;
    let mode = "normal";
    WebAssembly.instantiate = async (compiled, imports) => {
        instance = await instantiate(compiled, imports);
        wasi = imports.wasi_snapshot_preview1;
        if (mode === "initialize-failure") {
            return {
                exports: {
                    ...instance.exports,
                    _initialize() {
                        throw new Error("fixture initialization failure");
                    },
                },
            };
        }
        if (mode === "trap") {
            return {
                exports: {
                    ...instance.exports,
                    anki_backend_command() {
                        return instance.exports.anki_browser_buffer_data(
                            0xfffffff8,
                        );
                    },
                },
            };
        }
        const actual = instance;
        return {
            exports: {
                ...actual.exports,
                anki_backend_command(...args) {
                    commandHook?.();
                    return actual.exports.anki_backend_command(...args);
                },
            },
        };
    };
    fs.writeFileSync(path.join(root, "borrowed.txt"), "caller data");
    server = await createBrowserServer(module, provider, { root, timezone });
    assert.ok(owned.size > 0, "host must own a provider root descriptor");
    const allocation = instance.exports.anki_browser_buffer_new(
        4 * 65536 + 1024,
    );
    const pointer =
        instance.exports.anki_browser_buffer_data(allocation) >>> 0;
    const memory = () => new DataView(instance.exports.memory.buffer);
    const bytes = () => new Uint8Array(instance.exports.memory.buffer);
    const filename = pointer + 128;
    const result = pointer + 64;
    const rights = (1n << 28n) - 1n;
    function open(name, flags = 0, base = rights, fdFlags = 0, lookup = 0) {
        const encoded = new TextEncoder().encode(name);
        bytes().set(encoded, filename);
        const errno = wasi.path_open(
            3,
            lookup,
            filename,
            encoded.length,
            flags,
            base,
            0n,
            fdFlags,
            result,
        );
        return { errno, fd: memory().getUint32(result, true) };
    }
    assert.equal(open("absent/file", 1).errno, 44, "creation requires parent");
    assert.equal(open("borrowed.txt", 5).errno, 20, "exclusive creation");
    assert.equal(open("borrowed.txt", 2).errno, 54, "directory flag");
    assert.equal(
        open("borrowed.txt", 8, 2n).errno,
        2,
        "truncate needs write access",
    );
    const file = open("owned.bin", 1).fd;
    const payload = pointer + 256;
    bytes().set([0, 255, 128, 7], payload);
    memory().setUint32(pointer, payload, true);
    memory().setUint32(pointer + 4, 4, true);
    assert.equal(wasi.fd_write(file, pointer, 1, result), 0);
    assert.equal(memory().getUint32(result, true), 4);
    assert.equal(wasi.fd_seek(file, 1n, 0, result), 0);
    bytes()[payload] = 42;
    memory().setUint32(pointer + 4, 1, true);
    assert.equal(wasi.fd_write(file, pointer, 1, result), 0);
    assert.deepEqual(
        fs.readFileSync(path.join(root, "owned.bin")),
        Buffer.from([0, 42, 128, 7]),
    );
    assert.equal(wasi.fd_fdstat_set_flags(file, 1), 0);
    assert.equal(wasi.fd_seek(file, 0n, 0, result), 0);
    assert.equal(wasi.fd_write(file, pointer, 1, result), 0);
    assert.deepEqual(
        fs.readFileSync(path.join(root, "owned.bin")),
        Buffer.from([0, 42, 128, 7, 42]),
    );
    assert.equal(wasi.fd_seek(file, -1n, 0, result), 28);
    assert.equal(
        wasi.fd_seek(file, 1n << 54n, 0, result),
        61,
        "unsafe JS offset rejected",
    );
    assert.equal(wasi.fd_seek(file, 100n, 0, result), 0);
    assert.equal(wasi.fd_read(file, pointer, 1, result), 0);
    assert.equal(memory().getUint32(result, true), 0, "read beyond EOF");
    assert.equal(
        wasi.fd_read(file, 0xfffffff0, 16, result),
        21,
        "out-of-bounds iovec",
    );
    const readonly = open("owned.bin", 0, 2n).fd;
    assert.equal(
        wasi.fd_write(readonly, pointer, 1, result),
        76,
        "rights prevent write",
    );
    assert.equal(wasi.fd_close(readonly), 0);
    assert.equal(wasi.fd_close(readonly), 8);
    const encoded = new TextEncoder().encode("borrowed.txt");
    bytes().set(encoded, filename);
    assert.equal(wasi.path_remove_directory(3, filename, encoded.length), 54);
    assert.equal(
        wasi.path_filestat_get(3, 0, filename, encoded.length, result),
        0,
    );
    assert.equal(
        memory().getBigUint64(result + 8, true),
        fs.statSync(path.join(root, "borrowed.txt"), { bigint: true }).ino,
        "provider inode retained",
    );
    // Exercise actual host imports against real files, including progress before
    // provider faults. Collect each result so one old-source bug cannot hide another.
    const largePayload = pointer + 512;
    function vectors(lengths, skip = 0) {
        let offset = 0;
        let count = 0;
        for (const length of lengths) {
            const consumed = Math.min(skip, length);
            skip -= consumed;
            if (consumed < length) {
                memory().setUint32(
                    pointer + count * 8,
                    largePayload + offset + consumed,
                    true,
                );
                memory().setUint32(
                    pointer + count * 8 + 4,
                    length - consumed,
                    true,
                );
                count++;
            }
            offset += length;
        }
        return count;
    }
    function position(fd) {
        assert.equal(wasi.fd_tell(fd, result), 0);
        return Number(memory().getBigUint64(result, true));
    }
    for (const operation of ["read", "write"]) {
        for (const [layout, lengths, failAt, completed] of [
            ["chunk", [131073], 2, 65536],
            ["iovec", [70000, 131073], 3, 70000],
        ]) {
            behavior(
                `${operation} partial progress across ${layout} boundary`,
                (evidence) => {
                    const name = `${operation}-${layout}.bin`;
                    const total = lengths.reduce(
                        (sum, length) => sum + length,
                        0,
                    );
                    const expected = Buffer.from(
                        Array.from(
                            { length: total },
                            (_, index) => index % 251,
                        ),
                    );
                    fs.writeFileSync(
                        path.join(root, name),
                        operation === "read" ? expected : Buffer.alloc(0),
                    );
                    const fd = open(name).fd;
                    bytes().fill(0, largePayload, largePayload + total);
                    if (operation === "write")
                        bytes().set(expected, largePayload);
                    const count = vectors(lengths);
                    memory().setUint32(result, 0xdead, true);
                    inject(`${operation}Sync`, failAt, "EINTR");
                    const errno = wasi[`fd_${operation}`](
                        fd,
                        pointer,
                        count,
                        result,
                    );
                    const transferred = memory().getUint32(result, true);
                    evidence.first = {
                        errno,
                        transferred,
                        position: position(fd),
                    };
                    assert.equal(
                        errno,
                        0,
                        "completed bytes must return success",
                    );
                    assert.equal(transferred, completed);
                    assert.equal(evidence.first.position, completed);
                    assert.deepEqual(
                        operation === "write"
                            ? fs.readFileSync(path.join(root, name))
                            : Buffer.from(
                                  bytes().subarray(
                                      largePayload,
                                      largePayload + completed,
                                  ),
                              ),
                        expected.subarray(0, completed),
                    );
                    assert.equal(
                        wasi[`fd_${operation}`](
                            fd,
                            pointer,
                            vectors(lengths, transferred),
                            result,
                        ),
                        0,
                    );
                    assert.equal(
                        memory().getUint32(result, true),
                        total - transferred,
                    );
                    assert.equal(position(fd), total);
                    const actual =
                        operation === "write"
                            ? fs.readFileSync(path.join(root, name))
                            : Buffer.from(
                                  bytes().subarray(
                                      largePayload,
                                      largePayload + total,
                                  ),
                              );
                    assert.deepEqual(
                        actual,
                        expected,
                        "resume must neither duplicate nor skip bytes",
                    );
                    evidence.sha256 = createHash("sha256")
                        .update(actual)
                        .digest("hex");
                    evidence.finalPosition = position(fd);
                    assert.equal(wasi.fd_seek(fd, 0n, 0, result), 0);
                    inject(`${operation}Sync`, 1, "EINTR");
                    assert.equal(
                        wasi[`fd_${operation}`](
                            fd,
                            pointer,
                            vectors(lengths),
                            result,
                        ),
                        27,
                        "zero-progress failure remains errno",
                    );
                    assert.equal(position(fd), 0);
                    assert.equal(wasi.fd_close(fd), 0);
                },
            );
        }
    }
    for (const next of ["fd_write", "fd_sync", "fd_datasync", "fd_close"]) {
        behavior(
            `sync failure preserves progress and reaches ${next}`,
            (evidence) => {
                const name = `sync-${next}.bin`;
                const expected = Buffer.alloc(131073, 163);
                bytes().set(expected, largePayload);
                const fd = open(name, 9, rights, 16).fd;
                inject("fsyncSync", 1, "EIO");
                const errno = wasi.fd_write(
                    fd,
                    pointer,
                    vectors([expected.length]),
                    result,
                );
                const transferred = memory().getUint32(result, true);
                evidence.first = {
                    errno,
                    transferred,
                    position: position(fd),
                    fileBytes: fs.statSync(path.join(root, name)).size,
                };
                assert.equal(errno, 0);
                assert.equal(transferred, 65536);
                assert.equal(position(fd), transferred);
                assert.deepEqual(
                    fs.readFileSync(path.join(root, name)),
                    expected.subarray(0, transferred),
                );
                const args =
                    next === "fd_write"
                        ? [
                              fd,
                              pointer,
                              vectors([expected.length], transferred),
                              result,
                          ]
                        : [fd];
                assert.equal(
                    wasi[next](...args),
                    29,
                    "pending durability error must remain observable",
                );
                assert.equal(
                    fs.statSync(path.join(root, name)).size,
                    transferred,
                    "error delivery must not write more data",
                );
                if (next === "fd_close") {
                    assert.equal(wasi.fd_close(fd), 8);
                } else {
                    assert.equal(
                        wasi.fd_write(
                            fd,
                            pointer,
                            vectors([expected.length], transferred),
                            result,
                        ),
                        0,
                    );
                    assert.equal(wasi.fd_sync(fd), 0);
                    assert.equal(position(fd), expected.length);
                    assert.deepEqual(
                        fs.readFileSync(path.join(root, name)),
                        expected,
                    );
                    assert.equal(wasi.fd_close(fd), 0);
                }
            },
        );
    }
    for (const dangling of [false, true]) {
        behavior(
            `no-follow ${dangling ? "dangling" : "regular-target"} final symlink`,
            (evidence) => {
                const target = dangling ? "dangling-target" : "link-target";
                const name = `${target}.link`;
                if (!dangling)
                    fs.writeFileSync(
                        path.join(root, target),
                        "preserve target",
                    );
                fs.symlinkSync(target, path.join(root, name));
                const denied = open(name, 9, rights, 0, 0);
                evidence.noFollowErrno = denied.errno;
                evidence.target = fs.existsSync(path.join(root, target))
                    ? fs.readFileSync(path.join(root, target), "utf8")
                    : null;
                assert.equal(denied.errno, 32);
                assert.equal(
                    evidence.target,
                    dangling ? null : "preserve target",
                );
                const followed = open(name, 9, rights, 0, 1);
                assert.equal(
                    followed.errno,
                    0,
                    "explicit follow preserves truncation and creation",
                );
                assert.equal(wasi.fd_close(followed.fd), 0);
                assert.equal(
                    fs.readFileSync(path.join(root, target)).length,
                    0,
                );
                assert.ok(
                    fs.lstatSync(path.join(root, name)).isSymbolicLink(),
                );
            },
        );
    }
    behavior(
        "native no-follow flag rejects final symlink replacement before open",
        (evidence) => {
            fs.writeFileSync(path.join(root, "race-path"), "regular");
            fs.writeFileSync(
                path.join(root, "race-target"),
                "preserve race target",
            );
            beforeOpen = (name, flags) => {
                if (name !== path.join(root, "race-path")) return;
                beforeOpen = undefined;
                evidence.nativeNoFollow = Boolean(
                    flags & fs.constants.O_NOFOLLOW,
                );
                fs.unlinkSync(name);
                fs.symlinkSync("race-target", name);
            };
            const denied = open("race-path", 8);
            evidence.errno = denied.errno;
            evidence.target = fs.readFileSync(
                path.join(root, "race-target"),
                "utf8",
            );
            assert.equal(denied.errno, 32);
            assert.equal(evidence.nativeNoFollow, true);
            assert.equal(evidence.target, "preserve race target");
        },
    );
    instance.exports.anki_browser_buffer_free(allocation);
    server.close();
    server.close();
    server = undefined;
    assert.equal(
        owned.size,
        0,
        "normal close closes leaked guest descriptors",
    );
    assert.equal(
        fs.readFileSync(path.join(root, "borrowed.txt"), "utf8"),
        "caller data",
    );
    report.cases.push(
        "real reactor host preserves root, provider identity, flags, positional I/O, rights, memory checks and borrowed files",
    );
    mode = "initialize-failure";
    await assert.rejects(
        createBrowserServer(module, provider, { root, timezone }),
        /fixture initialization failure/,
    );
    assert.equal(owned.size, 0, "initialization failure closes preopen");
    report.cases.push(
        "initialization failure releases owned provider descriptors",
    );
    mode = "trap";
    server = await createBrowserServer(module, provider, { root, timezone });
    assert.throws(
        () =>
            server.command(
                "/anki.card_rendering.BackendCardRenderingService/StripHtml",
                new Uint8Array(),
            ),
        WebAssembly.RuntimeError,
    );
    assert.equal(owned.size, 0, "trap releases descriptors immediately");
    assert.throws(() => server.command("unused", new Uint8Array()));
    server.close();
    assert.equal(
        fs.readFileSync(path.join(root, "borrowed.txt"), "utf8"),
        "caller data",
    );
    report.cases.push(
        "actual Rust trap closes descriptors without destroying borrowed provider",
    );
    mode = "normal";
    for (const code of ["EACCES", "PROGRAMMING_ERROR"]) {
        for (const released of [false, true]) {
            server = await createBrowserServer(module, provider, {
                root,
                timezone,
            });
            const buffer = instance.exports.anki_browser_buffer_new(1024);
            const start =
                instance.exports.anki_browser_buffer_data(buffer) >>> 0;
            const view = new DataView(instance.exports.memory.buffer);
            const source = new Uint8Array(instance.exports.memory.buffer);
            const name = `close-${code}-${released}.bin`;
            const encoded = new TextEncoder().encode(name);
            source.set(encoded, start + 64);
            const before = new Set(owned);
            assert.equal(
                wasi.path_open(
                    3,
                    0,
                    start + 64,
                    encoded.length,
                    1,
                    rights,
                    0n,
                    16,
                    start + 32,
                ),
                0,
            );
            const fd = view.getUint32(start + 32, true);
            const providerFD = [...owned].find((value) => !before.has(value));
            assert.notEqual(providerFD, undefined);
            const attemptsBefore = closeAttempts.get(providerFD) || 0;
            source.set([0, 255, 128, 7], start + 256);
            view.setUint32(start, start + 256, true);
            view.setUint32(start + 4, 4, true);
            inject("fsyncSync", 1, "EIO");
            assert.equal(wasi.fd_write(fd, start, 1, start + 32), 0);
            assert.equal(view.getUint32(start + 32, true), 4);
            behavior(
                `deferred sync plus ${code} close ${released ? "after" : "before"} provider release`,
                (evidence) => {
                    commandHook = () => {
                        evidence.errno = wasi.fd_close(fd);
                    };
                    inject("closeSync", 1, code, released);
                    const command = () =>
                        server.command(
                            "/anki.card_rendering.BackendCardRenderingService/StripHtml",
                            new Uint8Array(),
                        );
                    try {
                        if (code === "EACCES") {
                            assert.equal(command().code, 0);
                            assert.equal(
                                evidence.errno,
                                2,
                                "close errno takes precedence over deferred sync errno",
                            );
                            commandHook = undefined;
                            assert.equal(
                                command().code,
                                0,
                                "recognized close error leaves backend usable",
                            );
                            instance.exports.anki_browser_buffer_free(buffer);
                        } else {
                            assert.throws(command);
                            assert.throws(
                                () =>
                                    server.command("unused", new Uint8Array()),
                                /closed/,
                            );
                        }
                    } finally {
                        commandHook = undefined;
                        assert.equal(
                            wasi.fd_close(fd),
                            8,
                            "guest descriptor is invalid after one close attempt",
                        );
                        server.close();
                        evidence.closeAttempts =
                            (closeAttempts.get(providerFD) || 0) -
                            attemptsBefore;
                        evidence.providerStillOpen = owned.has(providerFD);
                        assert.equal(
                            evidence.closeAttempts,
                            1,
                            "shutdown must not retry an uncertain close",
                        );
                        assert.equal(evidence.providerStillOpen, !released);
                        // Only the fault fixture knows whether the provider actually closed.
                        // Release its deliberately retained fd; the host must not guess.
                        if (owned.has(providerFD)) {
                            fs.closeSync(providerFD);
                            owned.delete(providerFD);
                        }
                    }
                    assert.deepEqual(
                        fs.readFileSync(path.join(root, name)),
                        Buffer.from([0, 255, 128, 7]),
                    );
                },
            );
            assert.equal(owned.size, 0, `${code} close fixture cleanup`);
        }
    }
    for (const [code, expected] of [
        ["EDQUOT", 19],
        ["ESTALE", 72],
        ["ENXIO", 60],
        ["ETIMEDOUT", 73],
        ["ECANCELED", 11],
        ["EREMOTEIO", 29],
        ["PROGRAMMING_ERROR", undefined],
    ]) {
        server = await createBrowserServer(module, provider, {
            root,
            timezone,
        });
        const buffer = instance.exports.anki_browser_buffer_new(132000);
        const start = instance.exports.anki_browser_buffer_data(buffer) >>> 0;
        const view = new DataView(instance.exports.memory.buffer);
        const source = new Uint8Array(instance.exports.memory.buffer);
        const name = `errno-${code}.bin`;
        const encoded = new TextEncoder().encode(name);
        source.set(encoded, start + 64);
        assert.equal(
            wasi.path_open(
                3,
                0,
                start + 64,
                encoded.length,
                1,
                rights,
                0n,
                0,
                start + 32,
            ),
            0,
        );
        const fd = view.getUint32(start + 32, true);
        source.fill(175, start + 256, start + 256 + 131073);
        view.setUint32(start, start + 256, true);
        view.setUint32(start + 4, 131073, true);
        behavior(
            `provider ${code} ${expected === undefined ? "is terminal after progress" : "is recoverable"}`,
            (evidence) => {
                commandHook = () => {
                    evidence.errno = wasi.fd_write(fd, start, 1, start + 32);
                };
                inject("writeSync", expected === undefined ? 2 : 1, code);
                try {
                    if (expected === undefined) {
                        assert.throws(
                            () =>
                                server.command(
                                    "/anki.card_rendering.BackendCardRenderingService/StripHtml",
                                    new Uint8Array(),
                                ),
                            /fixture PROGRAMMING_ERROR/,
                        );
                        assert.equal(
                            fs.statSync(path.join(root, name)).size,
                            65536,
                        );
                        assert.equal(owned.size, 0);
                        assert.throws(
                            () => server.command("unused", new Uint8Array()),
                            /closed/,
                        );
                    } else {
                        assert.equal(
                            server.command(
                                "/anki.card_rendering.BackendCardRenderingService/StripHtml",
                                new Uint8Array(),
                            ).code,
                            0,
                        );
                        assert.equal(evidence.errno, expected);
                        commandHook = undefined;
                        assert.equal(
                            wasi.fd_write(fd, start, 1, start + 32),
                            0,
                        );
                        assert.equal(
                            new DataView(
                                instance.exports.memory.buffer,
                            ).getUint32(start + 32, true),
                            131073,
                        );
                        assert.equal(wasi.fd_close(fd), 0);
                        assert.deepEqual(
                            fs.readFileSync(path.join(root, name)),
                            Buffer.alloc(131073, 175),
                        );
                        instance.exports.anki_browser_buffer_free(buffer);
                    }
                } finally {
                    commandHook = undefined;
                }
            },
        );
        server.close();
        assert.equal(owned.size, 0, `${code} descriptor cleanup`);
    }
    if (behaviorFailures.length)
        throw new AggregateError(
            behaviorFailures,
            `${behaviorFailures.length} filesystem behavior failures`,
        );
} catch (error) {
    failure = error;
    report.error = {
        name: error.name,
        message: error.message,
        stack: error.stack,
    };
} finally {
    WebAssembly.instantiate = instantiate;
    try {
        server?.close();
    } catch (error) {
        report.cleanupError = error.message;
    }
    report.remainingProviderDescriptors = owned.size;
    fs.writeFileSync(
        path.join(output, "filesystem.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
    fs.rmSync(root, { recursive: true, force: true });
}
if (failure) throw failure;

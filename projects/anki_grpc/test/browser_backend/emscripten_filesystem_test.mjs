import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import host from "node:fs/promises";
import path from "node:path";
import { pathToFileURL } from "node:url";

const [adapterPath] = process.argv
    .slice(2)
    .map((value) => path.resolve(value));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const modulePath = path.join(path.dirname(adapterPath), "emscripten_memfs.js");
const report = {
    scope: "Real generated classic MEMFS provider; not Go/Rust Chromium acceptance",
    node: process.version,
    cases: [],
    limitations: [],
    corrections: [],
};
const sha256 = (bytes) => createHash("sha256").update(bytes).digest("hex");
let failure;

async function check(name, operation) {
    await operation();
    report.cases.push(name);
}

// Each correction uses independent files/instances so one old-adapter failure
// does not hide evidence for the other review findings.
async function correction(id, operation) {
    try {
        await operation();
        report.corrections.push({ id, passed: true });
    } catch (error) {
        report.corrections.push({
            id,
            passed: false,
            message: error.message,
            code: error.code,
            stack: error.stack,
        });
    }
}

// Preserve the actual Node callback tuple, including the caller's buffer.
function call(fs, name, ...args) {
    return new Promise((resolve, reject) => {
        let returned = false;
        fs[name](...args, (error, ...values) => {
            try {
                assert.ok(returned, `${name} callback must be asynchronous`);
                if (error) reject(error);
                else resolve(values);
            } catch (failure) {
                reject(failure);
            }
        });
        returned = true;
    });
}

try {
    report.adapterSHA256 = sha256(await host.readFile(adapterPath));
    report.moduleSHA256 = sha256(await host.readFile(modulePath));
    const { createEmscriptenFilesystem } = await import(
        pathToFileURL(adapterPath).href
    );
    let initialized = false;
    let upstream;
    const initialization = createEmscriptenFilesystem({
        onRuntimeInitialized() {
            initialized = true;
            upstream = this.FS;
        },
    });
    assert.equal(typeof initialization.then, "function");
    const fs = await initialization;
    const other = await createEmscriptenFilesystem();
    const c = fs.constants;
    const root = "/provider-e2e";
    const filename = `${root}/雪\\literal.bin`;
    const bytes = Uint8Array.from(
        { length: 196613 },
        (_, i) => (i * 73 + (i >>> 8) * 19) & 255,
    );
    bytes.set([0, 255, 128, 92, 0, 254], 65534);

    await check(
        "initialized fresh factories reserve standard descriptors",
        async () => {
            assert.ok(initialized);
            assert.notEqual(fs, other);
            fs.mkdirSync(root);
            assert.throws(() => other.statSync(root), { code: "ENOENT" });
            for (const fd of [0, 1, 2]) {
                assert.equal(typeof fs.fstatSync(fd).ino, "number");
            }
            const [fd] = await call(
                fs,
                "open",
                filename,
                c.O_CREAT | c.O_EXCL | c.O_RDWR,
                0o640,
            );
            assert.ok(fd >= 3);
            fs.closeSync(fd);
            assert.equal(fs.statSync(filename).mode & 0o777, 0o640);
            assert.equal(
                fs.statSync(filename).ino,
                upstream.stat(filename).ino,
            );
            assert.equal(
                fs.statSync(filename).dev,
                upstream.stat(filename).dev,
            );
            assert.throws(() => other.openSync(filename, "r"), {
                code: "ENOENT",
            });
        },
    );

    await check(
        "callback and Sync positional binary I/O exceeds 64 KiB",
        async () => {
            const fd = fs.openSync(filename, "r+");
            const source = new Uint8Array(bytes.length + 11).fill(0xee);
            source.set(bytes, 7);
            const [written, returned] = await call(
                fs,
                "write",
                fd,
                source,
                7,
                bytes.length,
                0,
            );
            assert.equal(written, bytes.length);
            assert.equal(returned, source);
            source.fill(0); // MEMFS must not take ownership of caller memory.
            const first = new Uint8Array(9);
            assert.equal(
                fs.readSync(fd, first, 0, first.length, null),
                first.length,
            );
            assert.deepEqual(first, bytes.slice(0, first.length));
            const target = new Uint8Array(bytes.length + 19).fill(0xcc);
            const [read, returnedTarget] = await call(
                fs,
                "read",
                fd,
                target,
                13,
                bytes.length,
                0,
            );
            assert.equal(read, bytes.length);
            assert.equal(returnedTarget, target);
            assert.deepEqual(target.subarray(13, 13 + bytes.length), bytes);
            assert.ok(target.subarray(0, 13).every((value) => value === 0xcc));
            assert.ok(
                target
                    .subarray(13 + bytes.length)
                    .every((value) => value === 0xcc),
            );
            assert.equal(
                fs.readSync(fd, first, 0, first.length, null),
                first.length,
            );
            assert.deepEqual(first, bytes.slice(9, 18));
            assert.equal(fs.readSync(fd, first, 0, 1, bytes.length), 0);
            await call(fs, "close", fd);
            assert.deepEqual(upstream.readFile(filename), bytes);
            report.binarySHA256 = sha256(bytes);
            report.binarySize = bytes.length;
            await host.writeFile(
                path.join(output, "memfs-binary.bin"),
                target.subarray(13, 13 + bytes.length),
            );
        },
    );

    await check(
        "append, exclusive creation and truncate zero-fill",
        async () => {
            assert.throws(() => fs.openSync(filename, "wx"), {
                code: "EEXIST",
            });
            const append = fs.openSync(filename, "a+");
            assert.equal(
                fs.writeSync(append, new Uint8Array([231, 0]), 0, 2, 0),
                2,
            );
            assert.equal(fs.fstatSync(append).size, bytes.length + 2);
            const tail = new Uint8Array(2);
            fs.readSync(append, tail, 0, 2, bytes.length);
            assert.deepEqual(tail, new Uint8Array([231, 0]));
            fs.ftruncateSync(append, 5);
            await call(fs, "ftruncate", append, 70000);
            const grown = new Uint8Array(70000);
            fs.readSync(append, grown, 0, grown.length, 0);
            assert.deepEqual(grown.slice(0, 5), bytes.slice(0, 5));
            assert.ok(grown.subarray(5).every((value) => value === 0));
            fs.closeSync(append);
            await call(fs, "truncate", filename, 3);
            assert.equal(fs.statSync(filename).size, 3);
            fs.closeSync(fs.openSync(filename, "w"));
            assert.equal(fs.statSync(filename).size, 0);
        },
    );

    await check(
        "rename replacement and unlink preserve open inode identity",
        async () => {
            const source = fs.openSync(`${root}/source`, "w+");
            const replaced = fs.openSync(`${root}/target`, "w+");
            fs.writeSync(source, new Uint8Array([11, 22, 33]), 0, 3, 0);
            fs.writeSync(replaced, new Uint8Array([44]), 0, 1, 0);
            const sourceInode = fs.fstatSync(source).ino;
            const replacedInode = fs.fstatSync(replaced).ino;
            assert.notEqual(sourceInode, replacedInode);
            await call(fs, "rename", `${root}/source`, `${root}/target`);
            assert.equal(fs.statSync(`${root}/target`).ino, sourceInode);
            assert.equal(fs.fstatSync(replaced).ino, replacedInode);
            assert.equal(fs.fstatSync(replaced).size, 1);
            fs.ftruncateSync(replaced, 2);
            assert.equal(fs.statSync(`${root}/target`).size, 3);
            await call(fs, "unlink", `${root}/target`);
            assert.throws(() => fs.statSync(`${root}/target`), {
                code: "ENOENT",
            });
            assert.equal(fs.fstatSync(source).ino, sourceInode);
            const detached = new Uint8Array(3);
            fs.readSync(source, detached, 0, 3, 0);
            assert.deepEqual(detached, new Uint8Array([11, 22, 33]));
            fs.futimesSync(source, 1700000000.125, new Date(1700000100250));
            assert.equal(fs.fstatSync(source).mtimeMs, 1700000100250);
            await call(fs, "fchmod", source, 0o600);
            assert.equal(fs.fstatSync(source).mode & 0o777, 0o600);
            await call(fs, "fsync", source);
            fs.closeSync(source);
            fs.closeSync(replaced);
            report.limitations.push(
                "MEMFS reports upstream nlink, including on unlinked open nodes; fsync completes in memory only",
            );
        },
    );

    await check(
        "numeric and bigint Stats retain metadata and dates",
        async () => {
            await call(fs, "utimes", filename, 1700000000.125, 1700000100.25);
            const [numeric] = await call(fs, "stat", filename);
            const numericAgain = fs.statSync(filename);
            const big = fs.statSync(filename, { bigint: true });
            for (const key of [
                "dev",
                "ino",
                "mode",
                "nlink",
                "uid",
                "gid",
                "rdev",
                "size",
                "blksize",
                "blocks",
            ]) {
                assert.equal(typeof numeric[key], "number", key);
                assert.equal(big[key], BigInt(numeric[key]), key);
            }
            assert.equal(numeric.ino, numericAgain.ino);
            assert.equal(numeric.atimeMs, 1700000000125);
            assert.equal(numeric.mtimeMs, 1700000100250);
            assert.equal(big.mtimeNs, 1700000100250000000n);
            assert.equal(big.mtimeMs, 1700000100250n);
            assert.ok(numeric.atime instanceof Date);
            assert.ok(big.mtime instanceof Date);
            assert.ok(numeric.isFile());
            assert.ok(!numeric.isDirectory());
            assert.ok(!numeric.isSymbolicLink());
            assert.ok(!numeric.isSocket());
            await call(fs, "chmod", filename, 0o600);
            assert.equal(fs.statSync(filename).mode & 0o777, 0o600);
        },
    );

    await check(
        "symlinks preserve literal targets and no-follow rejects links",
        async () => {
            const target = "雪\\literal.bin";
            await call(fs, "symlink", target, `${root}/relative-link`);
            assert.equal(fs.readlinkSync(`${root}/relative-link`), target);
            assert.deepEqual(
                await call(fs, "readlink", `${root}/relative-link`),
                [target],
            );
            assert.ok(fs.lstatSync(`${root}/relative-link`).isSymbolicLink());
            assert.ok(
                (
                    await call(fs, "lstat", `${root}/relative-link`)
                )[0].isSymbolicLink(),
            );
            assert.equal(
                fs.statSync(`${root}/relative-link`).ino,
                fs.statSync(filename).ino,
            );
            assert.throws(
                () =>
                    fs.openSync(
                        `${root}/relative-link`,
                        c.O_RDONLY | c.O_NOFOLLOW,
                    ),
                { code: "ELOOP" },
            );
            fs.symlinkSync("missing", `${root}/dangling`);
            assert.ok(fs.lstatSync(`${root}/dangling`).isSymbolicLink());
            assert.throws(() => fs.statSync(`${root}/dangling`), {
                code: "ENOENT",
            });
            assert.throws(() => fs.readlinkSync(filename), { code: "EINVAL" });
        },
    );

    await check(
        "directory names, descriptor errors and unsupported operations",
        async () => {
            assert.equal(
                fs.mkdirSync(`${root}/nested/child`, {
                    recursive: true,
                    mode: 0o700,
                }),
                `${root}/nested`,
            );
            assert.equal(
                fs.mkdirSync(`${root}/nested/child`, { recursive: true }),
                undefined,
            );
            assert.throws(() => fs.mkdirSync(`${root}/nested/child`), {
                code: "EEXIST",
            });
            fs.mkdirSync(`${root}/nested/options`, { mode: 0o700 });
            assert.equal(
                fs.statSync(`${root}/nested/options`).mode & 0o777,
                0o700,
            );
            assert.throws(
                () => fs.mkdirSync(`${filename}/child`, { recursive: true }),
                { code: "ENOTDIR" },
            );
            await call(fs, "mkdir", `${root}/empty`, 0o750);
            assert.ok(fs.statSync(`${root}/empty`).isDirectory());
            const [names] = await call(fs, "readdir", root);
            assert.ok(names.includes("雪\\literal.bin"));
            assert.ok(!names.includes(".") && !names.includes(".."));
            assert.deepEqual(names.sort(), fs.readdirSync(root).sort());
            const dir = fs.openSync(root, c.O_RDONLY | c.O_DIRECTORY);
            assert.ok(fs.fstatSync(dir).isDirectory());
            assert.throws(() => fs.readSync(dir, new Uint8Array(1), 0, 1, 0), {
                code: "EISDIR",
            });
            fs.fsyncSync(dir);
            fs.closeSync(dir);
            assert.throws(() => fs.rmdirSync(root), { code: "ENOTEMPTY" });
            await call(fs, "rmdir", `${root}/empty`);
            assert.throws(() => fs.statSync(`${root}/empty`), {
                code: "ENOENT",
            });
            assert.throws(() => fs.readdirSync(filename), { code: "ENOTDIR" });
            assert.throws(
                () => fs.openSync(filename, c.O_RDONLY | c.O_DIRECTORY),
                { code: "ENOTDIR" },
            );
            const readOnly = fs.openSync(filename, "r");
            assert.throws(
                () => fs.writeSync(readOnly, new Uint8Array(1), 0, 1, null),
                { code: "EBADF" },
            );
            assert.throws(
                () => fs.readSync(readOnly, new Uint8Array(1), 0, 1, -1),
                { code: "EINVAL" },
            );
            assert.throws(
                () => fs.readSync(readOnly, new Uint8Array(1), 0, 2, 0),
                { code: "EINVAL" },
            );
            assert.throws(() => fs.ftruncateSync(readOnly, -1), {
                code: "EINVAL",
            });
            fs.closeSync(readOnly);
            await assert.rejects(call(fs, "fstat", readOnly), {
                code: "EBADF",
            });
            assert.throws(() => fs.closeSync(readOnly), { code: "EBADF" });
            assert.throws(() => fs.fsyncSync(readOnly), { code: "EBADF" });
            await assert.rejects(
                call(fs, "open", `${root}/absent`, c.O_RDONLY, 0),
                { code: "ENOENT" },
            );
            await assert.rejects(call(fs, "chown", filename, 1, 1), {
                code: "ENOTSUP",
            });
            const owned = fs.openSync(filename, "r+");
            await assert.rejects(call(fs, "fchown", owned, 1, 1), {
                code: "ENOTSUP",
            });
            fs.closeSync(owned);
            await assert.rejects(
                call(fs, "link", filename, `${root}/hard-link`),
                { code: "EMLINK" },
            );
            assert.throws(() => fs.statSync(`${root}/hard-link`), {
                code: "ENOENT",
            });
            report.limitations.push(
                "Classic MEMFS hard links fail EMLINK; ownership changes fail ENOTSUP",
            );
        },
    );

    await check(
        "factory initialization failure rejects without poisoning later factories",
        async () => {
            const sentinel = new Error(
                "provider E2E initialization rejection",
            );
            await assert.rejects(
                createEmscriptenFilesystem({
                    preRun: [
                        () => {
                            throw sentinel;
                        },
                    ],
                }),
                (error) => error === sentinel,
            );
            const after = await createEmscriptenFilesystem();
            after.mkdirSync("/ready");
            assert.ok(after.statSync("/ready").isDirectory());
        },
    );
    await correction("exclusive-create-existing-symlink", async () => {
        for (const [index, flags] of [
            "wx",
            c.O_WRONLY | c.O_CREAT | c.O_EXCL,
        ].entries()) {
            const target = `${root}/exclusive-target-${index}`;
            const link = `${root}/exclusive-link-${index}`;
            fs.symlinkSync(target, link);
            const inode = fs.lstatSync(link).ino;
            assert.throws(() => fs.openSync(link, flags), { code: "EEXIST" });
            assert.throws(() => fs.statSync(target), { code: "ENOENT" });
            assert.equal(fs.lstatSync(link).ino, inode);
            const fd = fs.openSync(target, "w");
            fs.writeSync(fd, new Uint8Array([0, 255, 19]), 0, 3, null);
            fs.closeSync(fd);
            await assert.rejects(call(fs, "open", link, flags, 0o600), {
                code: "EEXIST",
            });
            assert.deepEqual(
                upstream.readFile(target),
                new Uint8Array([0, 255, 19]),
            );
        }
    });

    await correction("concurrent-shared-factory-options", async () => {
        const modules = [];
        const preRuns = [];
        const postRuns = [];
        let preInitializations = 0;
        const initialize = () => {
            preInitializations++;
        };
        const preRun = (module) => {
            preRuns.push(module);
        };
        const postRun = (module) => {
            postRuns.push(module);
        };
        const options = {
            preInit: [initialize],
            preRun: [preRun],
            postRun: [postRun],
            onRuntimeInitialized() {
                modules.push(this);
            },
        };
        const before = Object.getOwnPropertyDescriptors(options);
        const [first, second] = await Promise.all([
            createEmscriptenFilesystem(options),
            createEmscriptenFilesystem(options),
        ]);
        first.mkdirSync("/first-only");
        assert.throws(() => second.statSync("/first-only"), {
            code: "ENOENT",
        });
        second.mkdirSync("/second-only");
        assert.throws(() => first.statSync("/second-only"), {
            code: "ENOENT",
        });
        assert.equal(modules.length, 2);
        assert.notEqual(modules[0], modules[1]);
        assert.notEqual(modules[0].FS, modules[1].FS);
        assert.notEqual(
            modules[0].FS.lookupPath("/").node,
            modules[1].FS.lookupPath("/").node,
        );
        assert.equal(preInitializations, 2);
        assert.deepEqual(preRuns, modules);
        assert.deepEqual(postRuns, modules);
        assert.deepEqual(Object.getOwnPropertyDescriptors(options), before);
        assert.deepEqual(options.preInit, [initialize]);
        assert.deepEqual(options.preRun, [preRun]);
        assert.deepEqual(options.postRun, [postRun]);
        const frozen = await createEmscriptenFilesystem(
            Object.freeze({
                preInit: Object.freeze([initialize]),
            }),
        );
        frozen.mkdirSync("/frozen-options-ready");
        assert.equal(preInitializations, 3);
    });

    await correction("unrepresentable-write-extent", async () => {
        const name = `${root}/range-error`;
        const fd = fs.openSync(name, "w+");
        const original = new Uint8Array([17, 28, 39]);
        const source = new Uint8Array([7]);
        fs.writeSync(fd, original, 0, original.length, 0);
        fs.readSync(fd, new Uint8Array(1), 0, 1, null);
        const opened = upstream.getStreamChecked(fd);
        // Both extents exceed the process address space and reject without a
        // large physical allocation. The second still fits safe JS arithmetic.
        report.writeExtentErrors = [];
        for (const position of [
            Number.MAX_SAFE_INTEGER,
            Number.MAX_SAFE_INTEGER - 1,
        ]) {
            assert.throws(
                () => fs.writeSync(fd, source, 0, 1, position),
                (error) => {
                    report.writeExtentErrors.push({
                        position,
                        code: error.code,
                        name: error.name,
                    });
                    return position === Number.MAX_SAFE_INTEGER
                        ? error.code === "EFBIG"
                        : ["EFBIG", "ENOMEM"].includes(error.code);
                },
            );
            assert.deepEqual(upstream.readFile(name), original);
            assert.equal(opened.position, 1);
        }
        upstream.llseek(opened, Number.MAX_SAFE_INTEGER, 0);
        await assert.rejects(call(fs, "write", fd, source, 0, 1, null), {
            code: "EFBIG",
        });
        assert.equal(opened.position, Number.MAX_SAFE_INTEGER);
        assert.deepEqual(upstream.readFile(name), original);
        fs.closeSync(fd);
        const append = fs.openSync(name, "a+");
        // For O_APPEND, the requested position is ignored: validate the actual
        // EOF extent, and retain the original stream position for pwrite.
        assert.equal(
            fs.writeSync(append, source, 0, 1, Number.MAX_SAFE_INTEGER),
            1,
        );
        assert.equal(upstream.getStreamChecked(append).position, 0);
        fs.closeSync(append);
        const finalBytes = upstream.readFile(name);
        assert.deepEqual(finalBytes, new Uint8Array([17, 28, 39, 7]));
        report.correctionBinarySHA256 = sha256(finalBytes);
        await host.writeFile(
            path.join(output, "memfs-corrections.bin"),
            finalBytes,
        );
    });

    await correction("read-write-callback-error-tuples", async () => {
        const name = `${root}/callback-errors`;
        const closed = fs.openSync(name, "w");
        fs.closeSync(closed);
        const original = new Uint8Array([11, 22, 33]).subarray(1);
        async function failedCall(method, fd) {
            let returned = false;
            let calls = 0;
            let asynchronous;
            const tuple = await new Promise((resolve) => {
                fs[method](fd, original, 0, 1, null, (...args) => {
                    calls++;
                    asynchronous = returned;
                    resolve(args);
                });
                returned = true;
            });
            await new Promise(setImmediate);
            assert.equal(asynchronous, true);
            assert.equal(calls, 1);
            assert.equal(tuple[0].code, "EBADF");
            assert.equal(tuple.length, 3);
            assert.equal(tuple[1], 0);
            assert.equal(tuple[2], original);
        }
        await failedCall("read", closed);
        await failedCall("write", closed);
        const readOnly = fs.openSync(name, "r");
        await failedCall("write", readOnly);
        fs.closeSync(readOnly);
        const writeOnly = fs.openSync(name, "w");
        await failedCall("read", writeOnly);
        fs.closeSync(writeOnly);
        for (const invalid of [false, true]) {
            const fd = fs.openSync(name, "r");
            if (invalid) fs.closeSync(fd);
            let calls = 0;
            const sentinel = new Error(`caller callback exception ${invalid}`);
            const thrown = await new Promise((resolve) => {
                process.once("uncaughtException", resolve);
                fs.read(fd, original, 0, 1, null, () => {
                    calls++;
                    throw sentinel;
                });
            });
            await new Promise(setImmediate);
            assert.equal(thrown, sentinel);
            assert.equal(calls, 1);
            if (!invalid) fs.closeSync(fd);
        }
    });

    const failedCorrections = report.corrections.filter(
        ({ passed }) => !passed,
    );
    assert.equal(
        failedCorrections.length,
        0,
        `Failed corrections: ${failedCorrections.map(({ id }) => id).join(", ")}`,
    );
    report.passed = true;
} catch (error) {
    failure = error;
    report.passed = false;
    report.error = {
        message: error.message,
        code: error.code,
        stack: error.stack,
    };
} finally {
    await host.writeFile(
        path.join(output, "emscripten-filesystem.json"),
        `${JSON.stringify(report, null, 2)}\n`,
    );
}
if (failure) throw failure;

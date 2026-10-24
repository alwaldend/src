import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { createHash } from "node:crypto";
import fs from "node:fs/promises";
import http from "node:http";
import { createRequire } from "node:module";
import path from "node:path";

// The shared npm workspace owns the repository-pinned browser dependency.
const require = createRequire(
    new URL("../../../../tools/package.json", import.meta.url),
);
const puppeteer = require("puppeteer");
const [archive, executablePath, extractor, fixtureTool] = process.argv
    .slice(2)
    .map((filename) => path.resolve(filename));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
assert.ok(output, "Bazel must provide an artifact directory");
const root = await fs.mkdtemp(path.join(process.env.TEST_TMPDIR, "browser-"));
const unpacked = spawnSync(extractor, [archive, root], {
    encoding: "utf8",
    timeout: 30000,
});
assert.equal(unpacked.status, 0, unpacked.stderr || unpacked.error?.message);
const bundleFiles = (await fs.readdir(root)).sort();
for (const filename of [
    "app.wasm",
    "index.d.ts",
    "index.html",
    "index.js",
    "wasm_exec.js",
    "worker.js",
])
    assert.ok(
        bundleFiles.includes(filename),
        `missing browser resource: ${filename}`,
    );
const fixtureArchive = path.join(output, "native-input.colpkg");
const fixtureBaseline = path.join(output, "native-baseline.json");
const fixture = spawnSync(
    fixtureTool,
    ["create", fixtureArchive, fixtureBaseline],
    {
        encoding: "utf8",
        timeout: 60000,
    },
);
assert.equal(fixture.status, 0, fixture.stderr || fixture.error?.message);

const requested = [];
const server = http.createServer(async (request, response) => {
    const pathname = new URL(request.url, "http://localhost").pathname;
    requested.push(pathname);
    if (pathname === "/test.html") {
        response.setHeader("Content-Type", "text/html");
        response.end("<!doctype html><title>Anki browser E2E</title>");
        return;
    }
    if (pathname === "/invalid.wasm") {
        response.setHeader("Content-Type", "application/wasm");
        response.end("invalid WebAssembly bytes");
        return;
    }
    if (pathname === "/broken-worker.js") {
        response.setHeader("Content-Type", "text/javascript");
        response.end('throw new Error("intentional worker startup failure");');
        return;
    }
    if (pathname === "/native-input.colpkg") {
        response.end(await fs.readFile(fixtureArchive));
        return;
    }
    if (pathname === "/no-jspi-worker.js") {
        response.setHeader("Content-Type", "text/javascript");
        response.end(
            'import "/worker.js"; WebAssembly.Suspending = undefined; WebAssembly.promising = undefined;',
        );
        return;
    }
    if (pathname === "/abort-save-worker.js") {
        response.setHeader("Content-Type", "text/javascript");
        response.end(`import "/worker.js";
const put = IDBObjectStore.prototype.put;
IDBObjectStore.prototype.put = function (...args) {
    IDBObjectStore.prototype.put = put;
    const request = put.apply(this, args);
    request.addEventListener("success", () => this.transaction.abort(), { once: true });
    return request;
};`);
        return;
    }
    const filename = pathname.slice(1);
    if (!bundleFiles.includes(filename)) {
        response.writeHead(404).end();
        return;
    }
    try {
        const types = {
            ".html": "text/html",
            ".js": "text/javascript",
            ".wasm": "application/wasm",
        };
        response.setHeader(
            "Content-Type",
            types[path.extname(filename)] || "application/octet-stream",
        );
        response.end(await fs.readFile(path.join(root, filename)));
    } catch (error) {
        response.writeHead(500).end(error.message);
    }
});
await new Promise((resolve) => server.listen(0, "127.0.0.1", resolve));
const base = `http://127.0.0.1:${server.address().port}`;
const browserErrors = [];
const consoleMessages = [];
const report = {
    archiveSHA256: createHash("sha256")
        .update(await fs.readFile(archive))
        .digest("hex"),
    bundleFiles,
};
let browser;
let page;
try {
    browser = await puppeteer.launch({ executablePath, headless: "shell" });
    report.browser = await browser.version();
    page = await browser.newPage();
    page.on("pageerror", (error) => browserErrors.push(error.message));
    page.on("console", (message) =>
        consoleMessages.push({ type: message.type(), text: message.text() }),
    );
    await page.goto(`${base}/test.html`, { waitUntil: "load" });
    const result = await page.evaluate(async () => {
        const { createAnkiApp } = await import("/index.js");
        const progress = { cases: [], rejected: {}, files: {} };
        globalThis.browserTestReport = progress;
        const check = (condition, message) => {
            if (!condition) throw new Error(message);
        };
        const equal = (actual, expected, message) =>
            check(
                JSON.stringify(actual) === JSON.stringify(expected),
                message,
            );
        const rejects = async (name, operation) => {
            try {
                await operation();
            } catch (error) {
                check(error instanceof Error, `${name}: expected an Error`);
                check(error.message.length > 0, `${name}: empty error`);
                progress.rejected[name] = error.message;
                return;
            }
            throw new Error(`${name}: expected rejection`);
        };
        const workers = [];
        const NativeWorker = globalThis.Worker;
        globalThis.Worker = class extends NativeWorker {
            constructor(...args) {
                super(...args);
                this.testTerminated = false;
                workers.push(this);
            }
            terminate() {
                this.testTerminated = true;
                return super.terminate();
            }
        };
        const instances = [];
        const create = async (options) => {
            const instance = await createAnkiApp(options);
            instances.push(instance);
            return instance;
        };
        const encoder = new TextEncoder();
        const decoder = new TextDecoder();
        const text =
            '<div class="example">你好, café</div>\nSecond line "quoted"';
        const configs = {
            config: { version: "1" },
            note_types: [
                {
                    path: "note_types/Basic.notetype.anki.toml",
                    config: {
                        name: "Basic",
                        fields: [{ name: "Front" }, { name: "Back" }],
                        templates: [
                            {
                                name: "Forward",
                                config: {
                                    q_format: "{{Front}}",
                                    a_format: "{{FrontSide}}\n{{Back}}",
                                },
                            },
                        ],
                    },
                },
            ],
            decks: [
                {
                    path: "decks/漢字/deck.anki.toml",
                    deck: { preset: "Default", normal: { new_limit: 0 } },
                    notes: [
                        {
                            path: "decks/漢字/notes/greeting.note.anki.toml",
                            config: {
                                note_type:
                                    "note_types/Basic.notetype.anki.toml",
                                fields: { Front: text, Back: "" },
                                tags: ["browser", " spaced "],
                                cards: [{}],
                            },
                        },
                    ],
                },
            ],
        };
        try {
            const first = await create();
            const { app, workspace } = first;
            progress.cases.push("module and worker ready");
            const original = JSON.stringify(configs);
            await app.configExport({ output: "/workspace/config", configs });
            equal(
                JSON.stringify(configs),
                original,
                "export changed caller input",
            );
            const names = await workspace.listFiles("/workspace/config");
            equal(
                names,
                [
                    "/workspace/config/collection.anki.toml",
                    "/workspace/config/decks/漢字/deck.anki.toml",
                    "/workspace/config/decks/漢字/notes/greeting.note.anki.toml",
                    "/workspace/config/note_types/Basic.notetype.anki.toml",
                ],
                "recursive file listing must be complete and sorted",
            );
            for (const name of names) {
                const bytes = await workspace.readFile(name);
                check(
                    bytes instanceof Uint8Array,
                    "read must return Uint8Array",
                );
                progress.files[name] = decoder.decode(bytes);
            }
            progress.json = await app.configDump({
                text: "/workspace/config",
            });
            const dumped = JSON.parse(progress.json);
            equal(
                dumped.config.version,
                "1",
                "64-bit version must be a string",
            );
            equal(dumped.decks[0].title, "漢字", "deck title must be derived");
            equal(
                dumped.decks[0].notes[0].config.fields,
                { Back: "", Front: text },
                "Unicode, multiline text, quotes and empty fields must survive",
            );
            equal(
                dumped.decks[0].notes[0].config.cards[0].deck,
                "decks/漢字/deck.anki.toml",
                "implicit card ownership must be resolved",
            );
            equal(
                dumped.decks[0].deck.normal.new_limit,
                0,
                "optional zero lost",
            );
            progress.toml = await app.configDump({
                text: "/workspace/config/collection.anki.toml",
                format: "toml",
            });
            check(
                progress.toml.includes("你好"),
                "TOML dump lost note content",
            );
            for (const name of names)
                equal(
                    decoder.decode(await workspace.readFile(name)),
                    progress.files[name],
                    "dump changed workspace source files",
                );
            progress.cases.push(
                "configuration export and JSON/TOML round trip",
            );

            for (const [name, contents] of Object.entries(progress.files))
                await workspace.writeFile(
                    name.replace("/workspace/config", "/imported"),
                    encoder.encode(contents),
                );
            const notePath =
                "/imported/decks/漢字/notes/greeting.note.anki.toml";
            await workspace.writeFile(
                notePath,
                encoder.encode(
                    'note_type = "note_types/Basic.notetype.anki.toml"\n' +
                        '[fields]\nFront = "edited Ж"\n',
                ),
            );
            const edited = JSON.parse(
                await app.configDump({ text: "/imported" }),
            );
            equal(
                edited.decks[0].notes[0].config.fields.Front,
                "edited Ж",
                "dump must observe newly imported bytes",
            );
            progress.cases.push("import and fresh reads after edits");

            const source = new Uint8Array(270);
            for (let index = 0; index < source.length; index++)
                source[index] = index % 256;
            const bytes = source.subarray(7, 263);
            const expected = Array.from(bytes);
            await workspace.writeFile("/binary/nested/bytes.bin", bytes);
            equal(
                Array.from(bytes),
                expected,
                "write detached or changed input",
            );
            equal(source.buffer.byteLength, 270, "caller buffer was detached");
            const received = await workspace.readFile(
                "/binary/nested/bytes.bin",
            );
            equal(
                Array.from(received),
                expected,
                "binary subarray bytes changed",
            );
            received.fill(0);
            equal(
                Array.from(
                    await workspace.readFile("/binary/nested/bytes.bin"),
                ),
                expected,
                "returned bytes must not alias workspace storage",
            );
            const writes = [
                workspace.writeFile("/ordered/value", encoder.encode("first")),
                workspace.writeFile(
                    "/ordered/value",
                    encoder.encode("second"),
                ),
            ];
            const read = workspace.readFile("/ordered/value");
            await Promise.all(writes);
            equal(
                decoder.decode(await read),
                "second",
                "request order changed",
            );
            await workspace.remove("/binary");
            await rejects("removed subtree", () =>
                workspace.readFile("/binary/nested/bytes.bin"),
            );
            progress.cases.push(
                "binary ownership, request ordering and removal",
            );

            const writeFiles = {
                "/writes/directory/retained.txt": "retained child",
                "/writes/file.txt": "retained file",
            };
            for (const [name, contents] of Object.entries(writeFiles))
                await workspace.writeFile(name, encoder.encode(contents));
            await workspace.writeFile(
                "/writes/empty/removed.txt",
                encoder.encode("removed"),
            );
            await workspace.remove("/writes/empty/removed.txt");
            const invalidWritePaths = [
                ["populated directory", "/writes/directory"],
                ["empty directory", "/writes/empty"],
                ["root directory", "/"],
                ["regular-file parent", "/writes/file.txt/child"],
                ["regular-file ancestor", "/writes/file.txt/missing/child"],
            ];
            for (const [name, pathname] of invalidWritePaths) {
                await rejects(`write to ${name}`, () =>
                    workspace.writeFile(pathname, encoder.encode("invalid")),
                );
                equal(
                    await workspace.listFiles("/writes"),
                    Object.keys(writeFiles).sort(),
                    `write to ${name} changed the file tree`,
                );
                for (const [filename, contents] of Object.entries(writeFiles))
                    equal(
                        decoder.decode(await workspace.readFile(filename)),
                        contents,
                        `write to ${name} changed ${filename}`,
                    );
            }
            const createdFile = "/writes/new/nested/created.txt";
            await workspace.writeFile(
                createdFile,
                encoder.encode("new file contents"),
            );
            await workspace.writeFile(createdFile, encoder.encode("replaced"));
            writeFiles[createdFile] = "replaced";
            const writeSave = await first.session.save(
                "browser-workspace-writes",
            );
            const written = await create();
            equal(
                await written.session.restore("browser-workspace-writes"),
                writeSave,
                "write snapshot metadata changed",
            );
            equal(
                await written.workspace.listFiles("/writes"),
                Object.keys(writeFiles).sort(),
                "save and restore lost files after rejected writes",
            );
            for (const [name, contents] of Object.entries(writeFiles)) {
                equal(
                    decoder.decode(await workspace.readFile(name)),
                    contents,
                    `write changed working file ${name}`,
                );
                equal(
                    decoder.decode(await written.workspace.readFile(name)),
                    contents,
                    `write changed saved file ${name}`,
                );
            }
            for (const [name, pathname] of invalidWritePaths)
                await rejects(`write to restored ${name}`, () =>
                    written.workspace.writeFile(
                        pathname,
                        encoder.encode("invalid"),
                    ),
                );
            written.dispose();
            progress.writeFiles = writeFiles;
            progress.cases.push(
                "rejected directory and ancestor writes preserve complete saved workspaces",
            );

            const renameFiles = {
                "/rename/source/moved.txt": "source contents",
                "/rename/destination/kept.txt": "destination contents",
                "/rename/published.colpkg": "original package",
                "/rename/staged.colpkg": "completed package",
            };
            for (const [name, contents] of Object.entries(renameFiles))
                await workspace.writeFile(name, encoder.encode(contents));
            for (const [name, source, destination] of [
                [
                    "rename onto a nonempty directory",
                    "/rename/source",
                    "/rename/destination",
                ],
                [
                    "rename directory onto its descendant",
                    "/rename/source",
                    "/rename/source/child",
                ],
                [
                    "rename directory onto a file",
                    "/rename/source",
                    "/rename/published.colpkg",
                ],
                [
                    "rename file onto a directory",
                    "/rename/staged.colpkg",
                    "/rename/destination",
                ],
                [
                    "rename beneath a file",
                    "/rename/staged.colpkg",
                    "/rename/published.colpkg/child",
                ],
                [
                    "rename with a missing parent",
                    "/rename/staged.colpkg",
                    "/rename/missing/child",
                ],
            ])
                await rejects(name, () =>
                    workspace.rename(source, destination),
                );
            await workspace.rename(
                "/rename/staged.colpkg",
                "/rename/completed.colpkg",
            );
            await workspace.rename(
                "/rename/completed.colpkg",
                "/rename/published.colpkg",
            );
            renameFiles["/rename/published.colpkg"] = "completed package";
            delete renameFiles["/rename/staged.colpkg"];
            await workspace.rename(
                "/rename/published.colpkg",
                "/rename/published.colpkg",
            );
            equal(
                await workspace.listFiles("/rename"),
                Object.keys(renameFiles).sort(),
                "rename lost or left behind directory entries",
            );
            const renameSave = await first.session.save(
                "browser-workspace-renames",
            );
            for (const revision of [-1, 0.5, Number.MAX_SAFE_INTEGER + 1, "1"])
                await rejects(`invalid save revision ${revision}`, () =>
                    first.session.save("invalid-browser-snapshot", revision),
                );
            equal(
                await first.session.restore("invalid-browser-snapshot"),
                null,
                "invalid revision created a saved workspace",
            );
            const renamed = await create();
            equal(
                await renamed.session.restore("browser-workspace-renames"),
                renameSave,
                "rename snapshot metadata changed",
            );
            equal(
                await renamed.workspace.listFiles("/rename"),
                Object.keys(renameFiles).sort(),
                "save and restore lost files after rejected renames",
            );
            for (const [name, contents] of Object.entries(renameFiles)) {
                equal(
                    decoder.decode(await workspace.readFile(name)),
                    contents,
                    `rename changed working file ${name}`,
                );
                equal(
                    decoder.decode(await renamed.workspace.readFile(name)),
                    contents,
                    `rename changed saved file ${name}`,
                );
            }
            renamed.dispose();
            progress.renameFiles = renameFiles;
            progress.cases.push(
                "file publication and rejected renames preserve complete saved workspaces",
            );

            const { loadWorkspace, saveWorkspace } = await import(
                "/storage.js"
            );
            const restoring = await create();
            const restoreFiles = {
                "/current/retained.bin": [0, 255, 128, 1],
                "/current/edit.txt": Array.from(encoder.encode("saved work")),
            };
            for (const [name, bytes] of Object.entries(restoreFiles))
                await restoring.workspace.writeFile(
                    name,
                    new Uint8Array(bytes),
                );
            const restoreKey = "browser-restore-preservation";
            const restoreSave = await restoring.session.save(restoreKey);
            const savedBeforeRestore = await loadWorkspace(restoreKey);
            restoreFiles["/current/edit.txt"] = Array.from(
                encoder.encode("unsaved edits"),
            );
            await restoring.workspace.writeFile(
                "/current/edit.txt",
                new Uint8Array(restoreFiles["/current/edit.txt"]),
            );
            const parentFile = {
                path: "/a",
                directory: false,
                data: encoder.encode("parent bytes"),
            };
            const childFile = {
                path: "/a/b",
                directory: false,
                data: encoder.encode("child bytes"),
            };
            const childDirectory = { path: "/a/b", directory: true };
            const invalidSnapshots = {
                "file before child file": [parentFile, childFile],
                "file after child file": [childFile, parentFile],
                "file before child directory": [parentFile, childDirectory],
                "file after child directory": [childDirectory, parentFile],
                "file before deeper descendant": [
                    parentFile,
                    { ...childFile, path: "/a/b/c" },
                ],
                "file after deeper descendant": [
                    { ...childFile, path: "/a/b/c" },
                    parentFile,
                ],
                "file at root": [{ ...parentFile, path: "/" }],
            };
            for (const [name, entries] of Object.entries(invalidSnapshots)) {
                const key = `browser-invalid-snapshot-${name}`;
                await saveWorkspace(key, 0, entries);
                const invalidBeforeRestore = await loadWorkspace(key);
                await rejects(`restore ${name}`, () =>
                    restoring.session.restore(key),
                );
                equal(
                    await restoring.workspace.listFiles("/"),
                    Object.keys(restoreFiles).sort(),
                    `restore ${name} replaced the current file tree`,
                );
                for (const [filename, bytes] of Object.entries(restoreFiles))
                    equal(
                        Array.from(
                            await restoring.workspace.readFile(filename),
                        ),
                        bytes,
                        `restore ${name} changed ${filename}`,
                    );
                equal(
                    await loadWorkspace(key),
                    invalidBeforeRestore,
                    `restore ${name} changed the invalid saved record`,
                );
                equal(
                    await loadWorkspace(restoreKey),
                    savedBeforeRestore,
                    `restore ${name} changed the previous saved workspace`,
                );
            }
            const preservedSave = await restoring.session.save(
                restoreKey,
                restoreSave.revision,
            );
            restoring.dispose();
            const restoreReader = await create();
            equal(
                await restoreReader.session.restore(restoreKey),
                preservedSave,
                "saved workspace metadata changed after rejected restores",
            );
            equal(
                await restoreReader.workspace.listFiles("/"),
                Object.keys(restoreFiles).sort(),
                "rejected restores lost files during save and reload",
            );
            for (const [name, bytes] of Object.entries(restoreFiles))
                equal(
                    Array.from(await restoreReader.workspace.readFile(name)),
                    bytes,
                    `rejected restores lost saved bytes at ${name}`,
                );
            const unorderedKey = "browser-unordered-snapshot";
            await saveWorkspace(unorderedKey, 0, [
                childFile,
                { path: "/a", directory: true },
                { path: "/", directory: true },
            ]);
            await restoreReader.session.restore(unorderedKey);
            equal(
                Array.from(
                    await restoreReader.workspace.readFile(childFile.path),
                ),
                Array.from(childFile.data),
                "compatible child-before-parent snapshot lost file bytes",
            );
            restoreReader.dispose();
            progress.restoreFiles = restoreFiles;
            progress.cases.push(
                "conflicting snapshot paths reject without changing saved bytes or unsaved work",
            );

            const second = await create();
            await rejects("separate instance", () =>
                second.workspace.readFile("/ordered/value"),
            );
            await second.workspace.writeFile(
                "/ordered/value",
                encoder.encode("other"),
            );
            equal(
                decoder.decode(await workspace.readFile("/ordered/value")),
                "second",
                "workers must have isolated storage",
            );
            progress.cases.push("independent application workspaces");

            const sessionKey = "browser-collection-round-trip";
            const collection = "/collection/collection.colpkg";
            const configDirectory = "/collection/config";
            const state = "/collection/state.toml";
            const native = await create();
            await native.workspace.writeFile(
                collection,
                new Uint8Array(
                    await (await fetch("/native-input.colpkg")).arrayBuffer(),
                ),
            );
            await native.app.collectionExport({
                collection,
                output: configDirectory,
                stateOutput: state,
            });
            await native.app.collectionStateUpdate({ collection, state });
            const exported = JSON.parse(
                await native.app.configDump({ text: configDirectory }),
            );
            const note = exported.decks
                .flatMap((deck) => deck.notes || [])
                .find(
                    (entry) =>
                        entry.config.fields.Front ===
                        "browser fixture original",
                );
            check(
                note,
                "native fixture note did not export through the browser server",
            );
            note.config.fields.Front = "browser edited 你好";
            await native.app.configExport({
                output: configDirectory,
                configs: exported,
            });
            const plan = await native.app.collectionPlan({
                collection,
                text: configDirectory,
            });
            check(
                plan && typeof plan === "object",
                "collection plan must return protobuf JSON",
            );
            // Queue save before apply resolves: the saved package must include the completed operation.
            const applying = native.app.collectionApply({
                collection,
                text: configDirectory,
            });
            const saving = native.session.save(sessionKey);
            await applying;
            const saved = await saving;
            equal(saved.revision, 1, "first completed snapshot revision");
            check(
                Number.isFinite(Date.parse(saved.savedAt)),
                "snapshot must expose its save time",
            );
            const savedFiles = {};
            for (const filename of await native.workspace.listFiles(
                configDirectory,
            )) {
                savedFiles[filename] = Array.from(
                    await native.workspace.readFile(filename),
                );
            }
            native.dispose();
            const restored = await create();
            await restored.workspace.writeFile(
                "/discard-on-restore",
                encoder.encode("stale"),
            );
            equal(
                await restored.session.restore(sessionKey),
                saved,
                "restore snapshot metadata changed",
            );
            await rejects("restore replaces workspace", () =>
                restored.workspace.readFile("/discard-on-restore"),
            );
            for (const [filename, contents] of Object.entries(savedFiles)) {
                equal(
                    Array.from(await restored.workspace.readFile(filename)),
                    contents,
                    "persisted configuration bytes changed",
                );
            }
            const absent = await restored.session.restore(
                "missing-browser-snapshot",
            );
            equal(absent, null, "missing snapshot must return null");
            equal(
                JSON.parse(
                    await restored.app.configDump({ text: configDirectory }),
                ),
                exported,
                "missing restore changed working configuration",
            );
            // A native verifier opens the downloaded result after this browser round trip.
            progress.collectionPackage = Array.from(
                await restored.workspace.readFile(collection),
            );
            progress.collectionFiles = savedFiles;
            progress.savedSession = saved;
            progress.cases.push(
                "native collection import, edit, apply, ordered save and fresh-worker restore",
            );

            const rival = await create();
            equal(
                await rival.session.restore(sessionKey),
                saved,
                "second worker restore baseline",
            );
            const winner = await restored.session.save(
                sessionKey,
                saved.revision,
            );
            equal(
                winner.revision,
                2,
                "snapshot update must advance the revision",
            );
            await rival.workspace.writeFile(
                "/conflicting-writer",
                encoder.encode("must not persist"),
            );
            await rejects("conflicting snapshot writer", () =>
                rival.session.save(sessionKey, saved.revision),
            );
            await rejects(
                "existing snapshot requires an observed revision",
                () => rival.session.save(sessionKey),
            );
            const reader = await create();
            equal(
                await reader.session.restore(sessionKey),
                winner,
                "conflict changed stored revision",
            );
            await rejects("conflict preserves previous snapshot", () =>
                reader.workspace.readFile("/conflicting-writer"),
            );
            progress.cases.push(
                "IndexedDB optimistic concurrency preserves the committed snapshot",
            );

            const aborted = await create({
                workerURL: new URL("/abort-save-worker.js", location.href),
            });
            equal(
                await aborted.session.restore(sessionKey),
                winner,
                "transaction-abort baseline",
            );
            await aborted.workspace.writeFile(
                "/aborted-write",
                encoder.encode("must not persist"),
            );
            await rejects("aborted IndexedDB transaction", () =>
                aborted.session.save(sessionKey, winner.revision),
            );
            equal(
                await reader.session.restore(sessionKey),
                winner,
                "failed transaction changed stored metadata",
            );
            await rejects(
                "aborted transaction preserves previous snapshot",
                () => reader.workspace.readFile("/aborted-write"),
            );
            equal(
                Array.from(await reader.workspace.readFile(collection)),
                progress.collectionPackage,
                "failed transaction changed saved collection bytes",
            );
            progress.cases.push(
                "failed IndexedDB transaction preserves the complete previous snapshot",
            );

            const withoutJSPI = await create({
                workerURL: new URL("/no-jspi-worker.js", location.href),
            });
            await withoutJSPI.app.configExport({ output: "/config", configs });
            check(
                (
                    await withoutJSPI.app.configDump({ text: "/config" })
                ).includes("你好"),
                "configuration must work without JSPI",
            );
            await withoutJSPI.workspace.writeFile(
                collection,
                new Uint8Array(
                    await (await fetch("/native-input.colpkg")).arrayBuffer(),
                ),
            );
            await withoutJSPI.app.collectionExport({
                collection,
                output: "/export",
            });
            check(
                (
                    await withoutJSPI.app.configDump({ text: "/export" })
                ).includes("browser fixture original"),
                "real collection operations must work without JSPI",
            );
            const unavailable = await create({
                serverWasmURL: new URL("/invalid.wasm", location.href),
            });
            await unavailable.app.configExport({ output: "/config", configs });
            await unavailable.workspace.writeFile(
                collection,
                await withoutJSPI.workspace.readFile(collection),
            );
            await rejects("invalid server WebAssembly", () =>
                unavailable.app.collectionExport({
                    collection,
                    output: "/export",
                }),
            );
            check(
                (
                    await unavailable.app.configDump({ text: "/config" })
                ).includes("你好"),
                "failed server startup must leave configuration editing usable",
            );
            progress.cases.push(
                "collections work without JSPI and failed server startup preserves configuration editing",
            );
            for (const instance of [
                restored,
                rival,
                reader,
                aborted,
                withoutJSPI,
                unavailable,
            ])
                instance.dispose();

            await rejects("unknown schema field", () =>
                app.configExport({
                    output: "/invalid-schema",
                    configs: { config: { version: "1" }, unknown_field: true },
                }),
            );
            check(
                !(await workspace.listFiles("/")).some((name) =>
                    name.startsWith("/invalid-schema/"),
                ),
                "invalid schema wrote output files",
            );
            await rejects("invalid integer", () =>
                app.configExport({
                    output: "/invalid-integer",
                    configs: { config: { version: "not-an-integer" } },
                }),
            );
            await rejects("null configuration", () =>
                app.configExport({ output: "/invalid-null", configs: null }),
            );
            await rejects("relative export path", () =>
                app.configExport({ output: "relative", configs }),
            );
            await rejects("relative dump path", () =>
                app.configDump({ text: "relative" }),
            );
            await rejects("missing config", () =>
                app.configDump({ text: "/missing-config" }),
            );
            await rejects("unsupported format", () =>
                app.configDump({
                    text: "/workspace/config",
                    format: "invalid",
                }),
            );
            await rejects("empty path", () => workspace.readFile(""));
            await rejects("relative write path", () =>
                workspace.writeFile("relative", new Uint8Array([1])),
            );
            await rejects("relative list path", () =>
                workspace.listFiles("relative"),
            );
            await rejects("relative removal path", () =>
                workspace.remove("relative"),
            );
            await rejects("wrong bytes type", () =>
                workspace.writeFile("/wrong", [1]),
            );
            await workspace.writeFile(
                "/malformed/collection.anki.toml",
                encoder.encode("["),
            );
            await rejects("malformed TOML", () =>
                app.configDump({ text: "/malformed" }),
            );
            equal(
                JSON.parse(
                    await app.configDump({ text: "/workspace/config" }),
                ),
                dumped,
                "valid operations must recover after rejected requests",
            );
            progress.cases.push(
                "input errors and recovery in the same application",
            );

            for (const [name, options] of [
                [
                    "missing wasm",
                    { wasmURL: new URL("/missing.wasm", location.href) },
                ],
                [
                    "invalid wasm",
                    { wasmURL: new URL("/invalid.wasm", location.href) },
                ],
                [
                    "failed worker",
                    { workerURL: new URL("/broken-worker.js", location.href) },
                ],
            ]) {
                const before = workers.length;
                await rejects(name, () => create(options));
                check(
                    workers
                        .slice(before)
                        .every((worker) => worker.testTerminated),
                    `${name}: failed startup must release its worker`,
                );
            }
            progress.cases.push(
                "fetch, compilation and worker startup failures",
            );

            const pending = workspace.readFile("/ordered/value");
            first.dispose();
            await rejects("pending request after dispose", () => pending);
            await rejects("future request after dispose", () =>
                workspace.readFile("/ordered/value"),
            );
            first.dispose();
            second.dispose();
            check(
                workers.every((worker) => worker.testTerminated),
                "workers leaked",
            );
            progress.cases.push(
                "pending and future rejection on idempotent disposal",
            );
            progress.workerCount = workers.length;
            return progress;
        } finally {
            for (const instance of instances) instance.dispose();
            globalThis.Worker = NativeWorker;
        }
    });
    Object.assign(report, result);
    assert.deepEqual(
        browserErrors.filter(
            (message) =>
                !message.includes("intentional worker startup failure"),
        ),
        [],
        "unexpected errors escaped the public API",
    );
    for (const name of [
        "/index.js",
        "/worker.js",
        "/wasm_exec.js",
        "/app.wasm",
    ])
        assert.ok(
            requested.includes(name),
            `packaged resource was not loaded: ${name}`,
        );
    for (const [filename, contents] of Object.entries(result.files)) {
        const destination = path.join(output, "workspace", filename.slice(1));
        await fs.mkdir(path.dirname(destination), { recursive: true });
        await fs.writeFile(destination, contents);
    }
    await fs.writeFile(path.join(output, "config-dump.json"), result.json);
    await fs.writeFile(path.join(output, "config-dump.toml"), result.toml);
    const browserPackage = path.join(output, "browser-output.colpkg");
    await fs.writeFile(
        browserPackage,
        new Uint8Array(result.collectionPackage),
    );
    delete report.collectionPackage;
    for (const [filename, contents] of Object.entries(
        result.collectionFiles,
    )) {
        const destination = path.join(output, "restored", filename.slice(1));
        await fs.mkdir(path.dirname(destination), { recursive: true });
        await fs.writeFile(destination, new Uint8Array(contents));
    }
    const verified = spawnSync(
        fixtureTool,
        [
            "verify",
            browserPackage,
            fixtureBaseline,
            path.join(output, "native-verification.json"),
        ],
        {
            encoding: "utf8",
            timeout: 60000,
        },
    );
    assert.equal(
        verified.status,
        0,
        verified.stderr || verified.error?.message,
    );
    console.log(`Verified ${result.cases.length} browser workflow groups`);
} catch (error) {
    report.failure = error.stack;
    if (page)
        Object.assign(
            report,
            await page
                .evaluate(() => globalThis.browserTestReport)
                .catch(() => ({})),
        );
    throw error;
} finally {
    Object.assign(report, { browserErrors, consoleMessages, requested });
    await fs.writeFile(
        path.join(output, "browser-report.json"),
        JSON.stringify(report, null, 2),
    );
    if (browser) await browser.close();
    await new Promise((resolve, reject) =>
        server.close((error) => (error ? reject(error) : resolve())),
    );
    await fs.rm(root, { recursive: true, force: true });
}

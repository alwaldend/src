import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { pathToFileURL } from "node:url";
import { runStorage } from "./storage_runtime.js";

const [application, guest, wasmExec, wasi, goHost] = process.argv
    .slice(2)
    .map((file) => path.resolve(file));
const output = process.env.TEST_UNDECLARED_OUTPUTS_DIR;
const root = fs.mkdtempSync(
    path.join(process.env.TEST_TMPDIR, "backend-node-storage-"),
);
const report = {
    host: "Node production WASI and shared OS filesystem",
    artifacts: {},
};
let deadline;
try {
    for (const file of [application, guest, wasmExec, wasi, goHost]) {
        report.artifacts[path.basename(file)] = createHash("sha256")
            .update(fs.readFileSync(file))
            .digest("hex");
    }
    await import(pathToFileURL(wasmExec).href);
    const { createBrowserServer } = await import(pathToFileURL(wasi).href);
    const { installGoFilesystem } = await import(pathToFileURL(goHost).href);
    report.execution = await Promise.race([
        runStorage({
            fs,
            root,
            module: await WebAssembly.compile(fs.readFileSync(guest)),
            application: fs.readFileSync(application),
            createBrowserServer,
            installGoFilesystem,
        }),
        new Promise((_, reject) => {
            deadline = setTimeout(
                () =>
                    reject(
                        new Error("Node shared storage exceeded 120 seconds"),
                    ),
                120000,
            );
        }),
    ]);
    const result = report.execution;
    assert.equal(result.error, undefined, result.error);
    assert.equal(result.exitCode, 0);
    assert.equal(result.workflow.error, undefined, result.workflow.error);
    assert.equal(result.instances, 2);
    assert.equal(result.remainingDescriptors, 0);
    assert.equal(result.timezone, "America/New_York");
    assert.deepEqual(result.timezoneOffsets, [
        [-18000, -14400],
        [-18000, -14400],
    ]);
    for (const method of [
        "openSync",
        "readSync",
        "writeSync",
        "closeSync",
        "readdirSync",
        "timezoneUTC",
    ])
        assert.ok(result.hostCalls[method] > 0, method);
    assert.ok(result.workflow.storage.noteId > 0);
    assert.ok(result.workflow.storage.collectionBytes > 4096);
    for (let index = 0; index < 2; index++) {
        assert.deepEqual(
            fs.readFileSync(
                path.join(
                    root,
                    `workspace/collection.media/concurrent-${index}.bin`,
                ),
            ),
            Buffer.from([index, 0, 255, 128]),
        );
    }
    fs.copyFileSync(
        path.join(root, "workspace/collection.anki2"),
        path.join(output, "collection.anki2"),
    );
    fs.copyFileSync(
        path.join(root, "workspace/collection.media/rust-media.bin"),
        path.join(output, "rust-media.bin"),
    );
} catch (error) {
    report.error = error.stack || error.message;
} finally {
    clearTimeout(deadline);
    fs.writeFileSync(
        path.join(output, "node-storage.json"),
        JSON.stringify(report, null, 2) + "\n",
    );
    fs.rmSync(root, { recursive: true, force: true });
}
if (report.error) console.error(report.error);
process.exit(report.error ? 1 : 0);

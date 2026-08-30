import assert from "node:assert/strict";
import { mkdtemp, readFile } from "node:fs/promises";
import { join } from "node:path";
import process from "node:process";
import test from "node:test";

import {
    ProcessSupervisor,
    waitForProcessGroup,
} from "../internal/process_supervisor.mjs";

const temporaryRoot = process.env.TEST_TMPDIR;
assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");

async function processIsLive(pid) {
    try {
        const stat = await readFile(`/proc/${pid}/stat`, "utf8");
        const state = stat.slice(stat.lastIndexOf(")") + 2).split(/\s+/u)[0];
        return state !== "Z" && state !== "X";
    } catch (error) {
        if (error?.code === "ENOENT") return false;
        throw error;
    }
}

function execute(supervisor, program, overrides = {}) {
    return supervisor.execute({
        file: process.execPath,
        args: ["-e", program],
        options: {
            cwd: process.cwd(),
            env: { ...process.env },
            maxBytes: 1024,
            timeoutMs: 5_000,
            allowTruncatedOutput: false,
            ...overrides,
        },
    });
}

test("process supervisor returns bounded stdout, stderr, and status", async () => {
    const supervisor = new ProcessSupervisor();
    const result = await execute(
        supervisor,
        "process.stdout.write('ok'); process.stderr.write('warn'); " +
            "process.exitCode = 7;",
    );
    assert.deepEqual(result, {
        code: 7,
        signal: null,
        stdout: "ok",
        stderr: "warn",
        truncated: false,
        outputLimitExceeded: false,
    });
    await supervisor.close();
});

test("process supervisor rejects excess output by default", async () => {
    const supervisor = new ProcessSupervisor();
    await assert.rejects(
        execute(supervisor, "process.stdout.write('abcdefgh');", {
            maxBytes: 4,
        }),
        (error) => error.code === "EXEC_OUTPUT_LIMIT",
    );
    await supervisor.close();
});

test("process supervisor can return a valid UTF-8 truncated prefix", async () => {
    const supervisor = new ProcessSupervisor();
    const result = await execute(
        supervisor,
        "process.stdout.write('a€tail');",
        { maxBytes: 3, allowTruncatedOutput: true },
    );
    assert.equal(result.stdout, "a");
    assert.equal(result.truncated, true);
    assert.equal(result.outputLimitExceeded, true);
    await supervisor.close();
});

test("process supervisor times out and closes its process group", async () => {
    const supervisor = new ProcessSupervisor();
    const directory = await mkdtemp(join(temporaryRoot, "exec-group-"));
    const pidFile = join(directory, "pids");
    const program = `
      const { spawn } = require("node:child_process");
      const { writeFileSync } = require("node:fs");
      const descendant = spawn(process.execPath, [
        "-e", "setInterval(() => {}, 1000)",
      ], { stdio: "ignore" });
      writeFileSync(${JSON.stringify(pidFile)},
        process.pid + " " + descendant.pid);
      setInterval(() => {}, 1000);
    `;
    await assert.rejects(
        execute(supervisor, program, {
            timeoutMs: 5_000,
        }),
        (error) => error.code === "EXEC_TIMEOUT",
    );
    const pids = (await readFile(pidFile, "utf8"))
        .trim()
        .split(" ")
        .map(Number);
    assert.equal(pids.length, 2);
    for (const pid of pids) assert.equal(await processIsLive(pid), false);

    const pending = execute(
        supervisor,
        "setInterval(() => {}, 1000);",
        { timeoutMs: 5_000 },
    );
    await new Promise((resolve) => setTimeout(resolve, 25));
    const disposed = Object.assign(new Error("disposed"), {
        code: "EXEC_DISPOSED",
    });
    await supervisor.close(disposed);
    await assert.rejects(pending, (error) => error.code === "EXEC_DISPOSED");
    await assert.rejects(
        execute(supervisor, "", {}),
        (error) => error.code === "EXEC_DISPOSED",
    );
});

test("process-group verification bounds permanent proc failures", async () => {
    const failure = Object.assign(new Error("proc unavailable"), {
        code: "ENOENT",
    });
    let attempts = 0;
    const error = await waitForProcessGroup(
        1,
        { exitCode: 0, pid: 0, signalCode: null },
        async () => {
            attempts += 1;
            throw failure;
        },
    );

    assert.equal(attempts, 3);
    assert.equal(error.code, "EXEC_CLEANUP");
    assert.match(error.message, /proc unavailable/);
});

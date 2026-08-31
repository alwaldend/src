import assert from "node:assert/strict";
import { mkdir, mkdtemp, readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";
import process from "node:process";
import test from "node:test";

import { Client } from "@modelcontextprotocol/client";
import { InMemoryTransport } from "@modelcontextprotocol/server";
import { createMcpServer } from "../internal/mcp.mjs";
import { CordisRuntime } from "../internal/runtime.mjs";

const temporaryRoot = process.env.TEST_TMPDIR;
assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
const defaultScratch = "out/test-task/mcp_cordis/runs/test-run";

async function workspace(label) {
    const root = await mkdtemp(join(temporaryRoot, `${label}-`));
    const project = join(root, "projects", "mcp_cordis");
    await mkdir(join(project, "plugins"), { recursive: true });
    await writeFile(join(project, "cordis.yaml"), "[]\n");
    return root;
}

function plugin(generation) {
    return `import { readFile, writeFile } from "node:fs/promises";
export default {
  description: ${JSON.stringify(`test plugin ${generation}`)},
  apply(ctx) {
    ctx.tool({
      name: "echo_value",
      description: "Return the active generation.",
      inputSchema: {
        type: "object",
        additionalProperties: false,
        required: ["value"],
        properties: { value: { type: "integer" } }
      }
    }, ({ value }) => ({ generation: ${JSON.stringify(generation)}, value }));
    ctx.tool({
      name: "slow_value",
      description: "Return after a bounded delay.",
      inputSchema: {
        type: "object",
        additionalProperties: false,
        required: ["delay_ms"],
        properties: {
          delay_ms: { type: "integer", minimum: 1, maximum: 1000 },
          completion_path: { type: "string" }
        }
      }
    }, async ({ delay_ms: delayMs, completion_path: completionPath }) => {
      await new Promise((resolve) => setTimeout(resolve, delayMs));
      if (completionPath) {
        await writeFile(ctx.resolveWorkspace(completionPath), "complete");
      }
      return { generation: ${JSON.stringify(generation)} };
    });
  }
};
`;
}

function slowEvaluationPlugin(generation, markerPath, releasePath) {
    return plugin(generation).replace(
        "export default {",
        `await writeFile(${JSON.stringify(markerPath)}, "started");
while (await readFile(${JSON.stringify(releasePath)}, "utf8").catch(() => "") !== "release") {
  await new Promise((resolve) => setTimeout(resolve, 10));
}
export default {`,
    );
}

function gatedApplyPlugin(generation, markerPath, releasePath) {
    return plugin(generation).replace(
        "  apply(ctx) {",
        `  async apply(ctx) {
    await writeFile(${JSON.stringify(markerPath)}, "started");
    while (await readFile(${JSON.stringify(releasePath)}, "utf8").catch(() => "") !== "release") {
      await new Promise((resolve) => setTimeout(resolve, 10));
    }`,
    );
}

function markedApplyPlugin(generation, markerPath) {
    return plugin(generation).replace(
        "  apply(ctx) {",
        `  async apply(ctx) {
    await writeFile(${JSON.stringify(markerPath)}, "started");`,
    );
}

function failingApplyPlugin(generation, markerPath) {
    return plugin(generation).replace(
        "  apply(ctx) {",
        `  async apply() {
    await writeFile(${JSON.stringify(markerPath)}, "started");
    throw new Error("intentional replacement failure");`,
    );
}

async function invoke(runtime, scope, packageName, value) {
    return runtime.invoke({
        scope,
        packageName,
        tool: "echo_value",
        arguments: { value },
    });
}

async function waitFor(callback, timeoutMs = 5_000) {
    const deadline = Date.now() + timeoutMs;
    let lastError;
    while (Date.now() < deadline) {
        try {
            if (await callback()) return;
        } catch (error) {
            lastError = error;
        }
        await new Promise((resolve) => setTimeout(resolve, 20));
    }
    assert.fail(
        `condition did not become true within ${timeoutMs} ms` +
            (lastError ? `: ${lastError.message}` : ""),
    );
}

function nextHmrReload(runtime, label, timeoutMs = 5_000) {
    return new Promise((resolve, reject) => {
        let dispose;
        const timer = setTimeout(() => {
            void dispose();
            reject(
                new Error(
                    `Cordis HMR did not reload ${label} within ${timeoutMs} ms`,
                ),
            );
        }, timeoutMs);
        timer.unref();
        dispose = runtime.root.on("hmr/reload", () => {
            clearTimeout(timer);
            void dispose();
            resolve();
        });
    });
}

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

function processPlugin(activationPidPath = undefined) {
    const childProgram = `
      const { spawn } = require("node:child_process");
      const { writeFileSync } = require("node:fs");
      const descendant = spawn(process.execPath, [
        "-e", "setInterval(() => {}, 1000)",
      ], { stdio: "ignore" });
      writeFileSync(process.argv[1], process.pid + " " + descendant.pid);
      setInterval(() => {}, 1000);
    `;
    const activation =
        activationPidPath === undefined
            ? ""
            : `void launch(${JSON.stringify(activationPidPath)}).catch(() => {});`;
    return `import { readFile } from "node:fs/promises";
export default {
  apply(ctx) {
    const launch = (pidPath) => ctx.exec(process.execPath, [
      "-e", ${JSON.stringify(childProgram)}, ctx.resolveWorkspace(pidPath),
    ], { timeoutMs: 30000 });
    ${activation}
    ctx.tool({
      name: "exec_tree",
      inputSchema: {
        type: "object",
        required: ["pid_path"],
        properties: { pid_path: { type: "string" } },
        additionalProperties: false
      }
    }, ({ pid_path: pidPath }) => launch(pidPath));
    ctx.tool({
      name: "fire_and_forget_tree",
      inputSchema: {
        type: "object",
        required: ["pid_path"],
        properties: { pid_path: { type: "string" } },
        additionalProperties: false
      }
    }, async ({ pid_path: pidPath }) => {
      const execution = launch(pidPath);
      void execution.catch(() => {});
      while (true) {
        try {
          await readFile(ctx.resolveWorkspace(pidPath));
          return { started: true };
        } catch (error) {
          if (error.code !== "ENOENT") throw error;
          await new Promise((resolve) => setTimeout(resolve, 5));
        }
      }
    });
  }
};
`;
}

async function readPids(root, relativePath) {
    return (await readFile(join(root, relativePath), "utf8"))
        .trim()
        .split(" ")
        .map(Number);
}

async function assertProcessesStopped(pids) {
    assert.equal(pids.length, 2);
    for (const pid of pids) assert.equal(await processIsLive(pid), false);
}

test("standard Cordis modules use native eventual HMR", async (t) => {
    const root = await workspace("runtime");
    const runtime = new CordisRuntime({ workspaceRoot: root });
    t.after(() => runtime.shutdown());

    assert.deepEqual(await runtime.initialize(), { loaded: [], errors: [] });
    const created = await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("v1"),
        activate: true,
    });
    assert.equal(created.module, "./plugins/echo.mjs");
    assert.equal(created.description, "");
    assert.equal(created.running, true);
    assert.deepEqual((await invoke(runtime, "scratch", "echo", 1)).value, {
        generation: "v1",
        value: 1,
    });
    const invalidCompletionPath = `${defaultScratch}/invalid_timeout_complete`;
    await assert.rejects(
        runtime.invoke({
            scope: "scratch",
            packageName: "echo",
            tool: "slow_value",
            arguments: {
                delay_ms: 40,
                completion_path: invalidCompletionPath,
            },
            timeoutMs: 0,
        }),
        /timeoutMs must be a positive integer/u,
    );
    await new Promise((resolve) => setTimeout(resolve, 60));
    await assert.rejects(
        readFile(join(root, invalidCompletionPath)),
        (error) => error.code === "ENOENT",
    );

    const completionPath = `${defaultScratch}/slow_value_complete`;
    await assert.rejects(
        runtime.invoke({
            scope: "scratch",
            packageName: "echo",
            tool: "slow_value",
            arguments: { delay_ms: 75, completion_path: completionPath },
            timeoutMs: 10,
        }),
        (error) => error.code === "invoke_timeout",
    );
    let stopped = false;
    const drainingStop = runtime
        .stop({ scope: "scratch", name: "echo" })
        .then(() => {
            stopped = true;
        });
    await new Promise((resolve) => setTimeout(resolve, 20));
    assert.equal(stopped, false);
    await drainingStop;
    assert.equal(
        await readFile(join(root, completionPath), "utf8"),
        "complete",
    );
    await runtime.run({ scope: "scratch", name: "echo" });
    assert.equal((await invoke(runtime, "scratch", "echo", 9)).value.value, 9);
    assert.equal(
        (await runtime.inspect({ scope: "scratch", name: "echo" })).enabled,
        true,
        "run must persist the enabled state",
    );

    const updated = await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("v2"),
    });
    assert.equal(updated.persisted, true);
    assert.equal(updated.sourceChanged, true);
    assert.equal(updated.activation, "pending");
    assert.equal(updated.running, true);
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 2)).value.generation ===
            "v2"
        );
    });
    const scratchSource = join(runtime.scratchRoot, "plugins", "echo.mjs");
    const v2Source = await readFile(scratchSource, "utf8");
    assert.equal(v2Source, plugin("v2"));
    await assert.rejects(
        runtime.define({
            scope: "scratch",
            name: "echo",
            source: "export default {",
        }),
        (error) => error.code === "invalid_module_source",
    );
    assert.equal(await readFile(scratchSource, "utf8"), v2Source);
    assert.deepEqual((await invoke(runtime, "scratch", "echo", 3)).value, {
        generation: "v2",
        value: 3,
    });

    const manualReload = nextHmrReload(runtime, "manual edit");
    await writeFile(scratchSource, plugin("manual"));
    await manualReload;
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "manual"
        );
    });
    assert.equal(
        (await runtime.inspect({ scope: "scratch", name: "echo" })).enabled,
        true,
        "manual reload must preserve the enabled state",
    );

    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("superseded"),
    });
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "superseded"
        );
    });
    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("latest"),
    });
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "latest"
        );
    });
    assert.equal(
        (await runtime.inspect({ scope: "scratch", name: "echo" })).enabled,
        true,
        "ordinary reloads must preserve the enabled state",
    );

    const slowReloadMarker = join(runtime.scratchRoot, "slow_reload_started");
    const slowReloadRelease = join(runtime.scratchRoot, "slow_reload_release");
    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: slowEvaluationPlugin(
            "slow",
            slowReloadMarker,
            slowReloadRelease,
        ),
    });
    const slowReloadRefresh = runtime.root.hmr.refreshFile(scratchSource);
    await waitFor(async () => {
        return (await readFile(slowReloadMarker, "utf8")) === "started";
    });
    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("rapid_latest"),
    });
    const rapidRefresh = runtime.root.hmr.refreshFile(scratchSource);
    await writeFile(slowReloadRelease, "release");
    await Promise.all([slowReloadRefresh, rapidRefresh]);
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "rapid_latest"
        );
    });

    const slowApplyMarker = join(runtime.scratchRoot, "slow_apply_started");
    const slowApplyRelease = join(runtime.scratchRoot, "slow_apply_release");
    const latestApplyMarker = join(runtime.scratchRoot, "latest_apply_started");
    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: gatedApplyPlugin(
            "slow_apply",
            slowApplyMarker,
            slowApplyRelease,
        ),
    });
    const slowApplyRefresh = runtime.root.hmr.refreshFile(scratchSource);
    await waitFor(async () => {
        return (await readFile(slowApplyMarker, "utf8")) === "started";
    });
    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: markedApplyPlugin("apply_latest", latestApplyMarker),
    });
    const applyRefresh = runtime.root.hmr.refreshFile(scratchSource);
    await assert.rejects(
        readFile(latestApplyMarker, "utf8"),
        (error) => error.code === "ENOENT",
    );
    await writeFile(slowApplyRelease, "release");
    await Promise.all([slowApplyRefresh, applyRefresh]);
    assert.equal(await readFile(latestApplyMarker, "utf8"), "started");
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "apply_latest"
        );
    });
    const beforeOverlapStop = await runtime.inspect({
        scope: "scratch",
        name: "echo",
    });
    assert.equal(
        beforeOverlapStop.enabled,
        true,
        JSON.stringify(beforeOverlapStop),
    );

    const failedApplyMarker = join(runtime.scratchRoot, "failed_apply_started");
    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: failingApplyPlugin("failed_apply", failedApplyMarker),
    });
    const failedRefresh = runtime.root.hmr.refreshFile(scratchSource);
    await waitFor(async () => {
        return (await readFile(failedApplyMarker, "utf8")) === "started";
    });
    await failedRefresh;
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "apply_latest"
        );
    });
    const afterFailedReplacement = await runtime.inspect({
        scope: "scratch",
        name: "echo",
    });
    assert.equal(afterFailedReplacement.enabled, true);
    assert.equal(afterFailedReplacement.running, true);

    await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("recovered"),
    });
    await runtime.root.hmr.refreshFile(scratchSource);
    await waitFor(async () => {
        return (
            (await invoke(runtime, "scratch", "echo", 3)).value.generation ===
            "recovered"
        );
    });

    const stoppedAfterOverlap = await runtime.stop({
        scope: "scratch",
        name: "echo",
    });
    assert.equal(
        stoppedAfterOverlap.running,
        false,
        JSON.stringify(stoppedAfterOverlap),
    );
    await assert.rejects(
        invoke(runtime, "scratch", "echo", 4),
        (error) => error.code === "tool_not_found",
    );
    const reactivated = await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("reactivated"),
        activate: true,
    });
    assert.equal(reactivated.activation, undefined);
    assert.equal(reactivated.running, true);
    assert.equal(
        (await invoke(runtime, "scratch", "echo", 4)).value.generation,
        "reactivated",
    );

    await runtime.stop({ scope: "scratch", name: "echo" });
    const disabledUpdate = await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("disabled_update"),
    });
    assert.equal(disabledUpdate.activation, undefined);
    const enabled = await runtime.define({
        scope: "scratch",
        name: "echo",
        source: plugin("disabled_update"),
        activate: true,
    });
    assert.equal(enabled.activation, undefined);
    assert.equal(
        (await invoke(runtime, "scratch", "echo", 5)).value.generation,
        "disabled_update",
    );

    const promoted = await runtime.promote({
        name: "echo",
        targetName: "reusable_echo",
        activate: true,
    });
    assert.equal(promoted.scope, "project");
    assert.equal(
        (await invoke(runtime, "project", "reusable_echo", 6)).value
            .generation,
        "disabled_update",
    );
    assert.match(
        await readFile(
            join(root, "projects", "mcp_cordis", "cordis.yaml"),
            "utf8",
        ),
        /id: reusable_echo[\s\S]*name: \.\/plugins\/reusable_echo\.mjs/u,
    );

    await runtime.remove({ scope: "scratch", name: "echo" });
    assert.equal(
        (await runtime.listPackages({ scope: "scratch" })).packages.length,
        0,
    );

    await runtime.shutdown();
    const restarted = new CordisRuntime({ workspaceRoot: root });
    t.after(() => restarted.shutdown());
    const startup = await restarted.initialize();
    assert.deepEqual(startup.errors, []);
    assert.equal(
        (await invoke(restarted, "project", "reusable_echo", 7)).value
            .generation,
        "disabled_update",
    );
});

test("unfiltered listing preserves a healthy scope after partial startup", async (t) => {
    for (const [label, config] of [
        ["malformed config", "[\n"],
        ["missing module", "- id: missing\n  name: ./plugins/missing.mjs\n"],
    ]) {
        await t.test(label, async (t) => {
            const root = await workspace("partial-startup");
            await writeFile(
                join(root, "projects", "mcp_cordis", "cordis.yaml"),
                config,
            );
            const runtime = new CordisRuntime({ workspaceRoot: root });
            t.after(() => runtime.shutdown());

            const startup = await runtime.initialize();
            assert.equal(startup.errors.length, 1);
            assert.equal(startup.errors[0].scope, "project");
            const originalMessage = startup.errors[0].error.message;
            startup.errors[0].error.message = "caller mutation";
            const repeatedStartup = await runtime.initialize();
            assert.equal(repeatedStartup.errors.length, 1);
            assert.equal(repeatedStartup.errors[0].scope, "project");
            assert.equal(
                repeatedStartup.errors[0].error.message,
                originalMessage,
            );
            await runtime.define({
                scope: "scratch",
                name: "healthy",
                source: plugin("healthy"),
            });

            const listed = await runtime.listPackages();
            assert.deepEqual(
                listed.packages.map(({ scope, name }) => ({ scope, name })),
                [{ scope: "scratch", name: "healthy" }],
            );
            assert.equal(listed.errors.length, 1);
            assert.equal(listed.errors[0].scope, "project");
            await assert.rejects(
                runtime.listPackages({ scope: "project" }),
                (error) => error.code === "scope_unavailable",
            );
        });
    }
});

test("MCP gateway exposes and invokes a runtime plugin", async (t) => {
    const root = await workspace("mcp");
    const runtime = new CordisRuntime({ workspaceRoot: root });
    await runtime.initialize();
    const server = createMcpServer(runtime);
    const client = new Client({ name: "mcp-cordis-test", version: "1.0.0" });
    const [clientTransport, serverTransport] =
        InMemoryTransport.createLinkedPair();
    t.after(async () => {
        await client.close().catch(() => {});
        await server.close().catch(() => {});
        await runtime.shutdown();
    });
    await Promise.all([
        client.connect(clientTransport),
        server.connect(serverTransport),
    ]);

    const result = await client.callTool({
        name: "cordis_define",
        arguments: {
            scope: "scratch",
            name: "echo",
            source: plugin("mcp"),
            activate: true,
        },
    });
    assert.equal(result.isError, undefined);
    const invoked = await client.callTool({
        name: "cordis_invoke",
        arguments: {
            scope: "scratch",
            package: "echo",
            tool: "echo_value",
            arguments: { value: 8 },
        },
    });
    assert.deepEqual(invoked.structuredContent.value, {
        generation: "mcp",
        value: 8,
    });
});

test("standard asynchronous Cordis plugins are supported", async (t) => {
    const root = await workspace("async-activation");
    const runtime = new CordisRuntime({ workspaceRoot: root });
    t.after(() => runtime.shutdown());
    await runtime.initialize();

    const created = await runtime.define({
        scope: "scratch",
        name: "asynchronous",
        source: `await Promise.resolve();
          export default {
            async apply(ctx) {
              await Promise.resolve();
              ctx.tool({
                name: "async_value",
                inputSchema: {
                  type: "object",
                  additionalProperties: false
                }
              }, () => ({ ready: true }));
            }
          };`,
        activate: true,
    });
    assert.equal(created.running, true);
    const invoked = await runtime.invoke({
        scope: "scratch",
        packageName: "asynchronous",
        tool: "async_value",
    });
    assert.deepEqual(invoked.value, { ready: true });
});

test("source at the payload limit round-trips without metadata", async (t) => {
    const root = await workspace("source-limit");
    const runtime = new CordisRuntime({ workspaceRoot: root });
    t.after(() => runtime.shutdown());
    await runtime.initialize();

    const source = `//${"x".repeat(1_999_997)}\n`;
    assert.equal(Buffer.byteLength(source), 2_000_000);
    await runtime.define({ scope: "scratch", name: "limit", source });
    const inspected = await runtime.inspect({
        scope: "scratch",
        name: "limit",
        includeSource: true,
    });
    assert.equal(Buffer.byteLength(inspected.source), 2_000_000);
    assert.equal(inspected.source, source);
    const roundTrip = await runtime.define({
        scope: "scratch",
        name: "limit",
        source: inspected.source,
    });
    assert.equal(roundTrip.updated, true);
    const legacyRoundTrip = await runtime.define({
        scope: "scratch",
        name: "limit",
        source:
            `export const __mcp_cordis_source_sha256 = ` +
            `"${"a".repeat(64)}";\n${source}`,
    });
    assert.equal(legacyRoundTrip.sourceChanged, false);
    assert.equal(
        (
            await runtime.inspect({
                scope: "scratch",
                name: "limit",
                includeSource: true,
            })
        ).source,
        source,
    );
});

test("invocation and Fiber disposal join process-tree cleanup", async (t) => {
    const root = await workspace("process-wiring");
    const runtime = new CordisRuntime({ workspaceRoot: root });
    t.after(() => runtime.shutdown());
    await runtime.initialize();

    await runtime.define({
        scope: "scratch",
        name: "process_tools",
        source: processPlugin(),
        activate: true,
    });
    const timeoutPidPath = `${defaultScratch}/invoke_timeout_pids`;
    const timedInvocation = runtime.invoke({
        scope: "scratch",
        packageName: "process_tools",
        tool: "exec_tree",
        arguments: { pid_path: timeoutPidPath },
        timeoutMs: 300,
    });
    const timedRejection = assert.rejects(
        timedInvocation,
        (error) => error.code === "invoke_timeout",
    );
    await waitFor(async () =>
        Boolean(await readFile(join(root, timeoutPidPath))),
    );
    await timedRejection;
    await assertProcessesStopped(await readPids(root, timeoutPidPath));

    const forgottenPidPath = `${defaultScratch}/fire_and_forget_pids`;
    const forgottenInvocation = runtime.invoke({
        scope: "scratch",
        packageName: "process_tools",
        tool: "fire_and_forget_tree",
        arguments: { pid_path: forgottenPidPath },
    });
    await waitFor(async () =>
        Boolean(await readFile(join(root, forgottenPidPath))),
    );
    const stopped = runtime.stop({ scope: "scratch", name: "process_tools" });
    await Promise.all([forgottenInvocation, stopped]);
    await assertProcessesStopped(await readPids(root, forgottenPidPath));

    const activationPidPath = `${defaultScratch}/activation_pids`;
    await runtime.define({
        scope: "scratch",
        name: "activation_process",
        source: processPlugin(activationPidPath),
        activate: true,
    });
    await waitFor(async () =>
        Boolean(await readFile(join(root, activationPidPath))),
    );
    await runtime.stop({ scope: "scratch", name: "activation_process" });
    await assertProcessesStopped(await readPids(root, activationPidPath));
});

test("task and run ownership isolate concurrent scratch", async (t) => {
    const root = await workspace("scratch-isolation");
    const first = new CordisRuntime({
        workspaceRoot: root,
        taskId: "task-one",
        runId: "run-one",
        workerId: "worker-one",
    });
    const second = new CordisRuntime({
        workspaceRoot: root,
        taskId: "task-two",
        runId: "run-one",
        workerId: "worker-two",
    });
    t.after(() => Promise.all([first.shutdown(), second.shutdown()]));
    await Promise.all([first.initialize(), second.initialize()]);

    assert.notEqual(first.scratchRoot, second.scratchRoot);
    await first.define({
        scope: "scratch",
        name: "isolated",
        source: plugin("task-one"),
        activate: false,
    });
    assert.equal(
        (await second.listPackages({ scope: "scratch" })).packages.length,
        0,
    );

    for (const [runtime, taskId, workerId] of [
        [first, "task-one", "worker-one"],
        [second, "task-two", "worker-two"],
    ]) {
        const manifest = JSON.parse(
            await readFile(join(runtime.scratchRoot, "manifest.json"), "utf8"),
        );
        assert.equal(manifest.taskId, taskId);
        assert.equal(manifest.runId, "run-one");
        assert.equal(manifest.workerId, workerId);
        assert.equal(manifest.lockScope, `${taskId}/run-one`);
    }
});

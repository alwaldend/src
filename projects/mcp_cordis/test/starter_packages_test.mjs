import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { cp, mkdir, mkdtemp, writeFile } from "node:fs/promises";
import { createServer } from "node:http";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";
import test from "node:test";
import { CordisRuntime } from "../internal/runtime.mjs";

const execute = promisify(execFile);
const temporaryRoot = process.env.TEST_TMPDIR;
assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");

async function git(root, ...args) {
    return execute("git", ["-C", root, ...args], {
        encoding: "utf8",
        maxBuffer: 1024 * 1024,
    });
}

async function startHttpServer(t) {
    const server = createServer((request, response) => {
        response.setHeader("Content-Type", "text/plain; charset=utf-8");
        response.setHeader("Set-Cookie", "session=must-not-leak");
        if (request.url === "/utf8") {
            response.end("éx");
            return;
        }
        if (request.url === "/invalid-utf8") {
            response.end(Buffer.from([0xc3]));
            return;
        }
        response.end(`probe:${request.method}:` + "x".repeat(4096));
    });
    await new Promise((resolve, reject) => {
        server.once("error", reject);
        server.listen(0, "127.0.0.1", resolve);
    });
    t.after(() => new Promise((resolve) => server.close(resolve)));
    return server.address().port;
}

test("checked-in starter packages execute through Cordis", async (t) => {
    const root = await mkdtemp(join(temporaryRoot, "starters-"));
    const destination = join(root, "projects", "mcp_cordis");
    await mkdir(destination, { recursive: true });
    const plugins = fileURLToPath(new URL("../plugins", import.meta.url));
    const config = fileURLToPath(new URL("../cordis.yaml", import.meta.url));
    await cp(plugins, join(destination, "plugins"), { recursive: true });
    await cp(config, join(destination, "cordis.yaml"));

    await writeFile(join(root, "AGENTS.md"), "fixture instructions\n");
    await writeFile(join(root, "README.md"), "fixture readme\n");
    await writeFile(join(root, "large.txt"), "base\n");
    await writeFile(join(root, "sample.txt"), "needle_one\n");
    await git(root, "init", "--quiet");
    await git(root, "config", "user.name", "MCP Cordis Test");
    await git(root, "config", "user.email", "mcp-cordis@example.invalid");
    await git(root, "config", "commit.gpgsign", "false");
    await git(root, "add", ".");
    await git(root, "commit", "--quiet", "-m", "first");
    const base = (await git(root, "rev-parse", "HEAD")).stdout.trim();
    await writeFile(
        join(root, "large.txt"),
        Array.from(
            { length: 300 },
            (_, index) => `expanded line ${index.toString().padStart(3, "0")}`,
        ).join("\n") + "\n",
    );
    await writeFile(join(root, "sample.txt"), "needle_one\nneedle_two\n");
    await git(root, "add", "large.txt", "sample.txt");
    await git(root, "commit", "--quiet", "-m", "second");
    const head = (await git(root, "rev-parse", "HEAD")).stdout.trim();
    await writeFile(
        join(root, "sample.txt"),
        "needle_one\nneedle_two\ndirty\n",
    );

    const runtime = new CordisRuntime({ workspaceRoot: root });
    t.after(() => runtime.shutdown());
    const startup = await runtime.initialize();
    assert.equal(startup.errors.length, 0);
    assert.deepEqual(
        startup.loaded.map((item) => item.name).sort(),
        ["git_worktree", "network_probe", "repo_context"],
    );
    const snapshots = await runtime.listPackages({ scope: "project" });
    for (const item of snapshots.packages) {
        assert.match(item.description, /\S/u, item.name);
    }
    assert.equal(runtime.listTools().tools.length, 8);

    const context = await runtime.invoke({
        scope: "project",
        packageName: "repo_context",
        tool: "repo_context_get",
        arguments: { path: ".", max_bytes: 16_384 },
    });
    assert.equal(context.value.path, ".");
    assert.equal(context.value.instructions[0].path, "AGENTS.md");

    const search = await runtime.invoke({
        scope: "project",
        packageName: "repo_context",
        tool: "repo_context_search",
        arguments: {
            query: "needle_two",
            paths: ["sample.txt"],
            max_matches: 5,
        },
    });
    assert.equal(search.value.matchCount, 1);
    assert.equal(search.value.matches[0].line, 2);

    const read = await runtime.invoke({
        scope: "project",
        packageName: "repo_context",
        tool: "repo_context_read",
        arguments: {
            files: [{ path: "sample.txt", start_line: 2, end_line: 3 }],
            max_total_bytes: 32,
        },
    });
    assert.equal(read.value.files[0].content, "needle_two\ndirty");
    await assert.rejects(
        runtime.invoke({
            scope: "project",
            packageName: "repo_context",
            tool: "repo_context_read",
            arguments: { files: [{ path: "../outside" }] },
        }),
        /outside|escape/i,
    );

    const snapshot = await runtime.invoke({
        scope: "project",
        packageName: "git_worktree",
        tool: "git_snapshot",
        arguments: { repo: ".", log_limit: 2 },
    });
    assert.equal(snapshot.value.branch.oid, head);
    assert.ok(snapshot.value.changes.some((item) => {
        return item.path === "sample.txt";
    }));
    assert.equal(snapshot.value.commits.length, 2);

    const comparison = await runtime.invoke({
        scope: "project",
        packageName: "git_worktree",
        tool: "git_compare",
        arguments: { repo: ".", base, head, max_bytes: 32_768 },
    });
    assert.equal(comparison.value.base, base);
    assert.equal(comparison.value.head, head);
    assert.match(comparison.value.diff, /needle_two/);

    const boundedComparison = await runtime.invoke({
        scope: "project",
        packageName: "git_worktree",
        tool: "git_compare",
        arguments: { repo: ".", base, head, max_bytes: 1024 },
    });
    assert.equal(boundedComparison.value.truncated, true);
    assert.equal(boundedComparison.value.filesTruncated, false);
    assert.equal(boundedComparison.value.diffTruncated, true);
    assert.equal(boundedComparison.value.diffBytes, 1024);
    assert.equal(
        Buffer.byteLength(boundedComparison.value.diff, "utf8"),
        1024,
    );
    assert.match(
        boundedComparison.value.diff,
        /^diff --git a\/large\.txt b\/large\.txt/,
    );
    assert.match(boundedComparison.value.diff, /expanded line 000/);
    assert.match(boundedComparison.value.diff, /^[\x00-\x7f]+$/);
    assert.deepEqual(
        boundedComparison.value.files.map((item) => item.path).sort(),
        ["large.txt", "sample.txt"],
    );

    const port = await startHttpServer(t);
    const tcp = await runtime.invoke({
        scope: "project",
        packageName: "network_probe",
        tool: "tcp_probe",
        arguments: { host: "127.0.0.1", port, timeout_ms: 2_000 },
    });
    assert.equal(tcp.value.ok, true);
    assert.equal(tcp.value.remotePort, port);

    const http = await runtime.invoke({
        scope: "project",
        packageName: "network_probe",
        tool: "http_probe",
        arguments: {
            url: `http://127.0.0.1:${port}/health`,
            max_body_bytes: 32,
            timeout_ms: 2_000,
        },
    });
    assert.equal(http.value.ok, true);
    assert.equal(http.value.status, 200);
    assert.equal(http.value.bodyBytes, 32);
    assert.equal(http.value.bodyTruncated, true);
    assert.equal(http.value.headers["set-cookie"], "[redacted]");

    const splitUtf8 = await runtime.invoke({
        scope: "project",
        packageName: "network_probe",
        tool: "http_probe",
        arguments: {
            url: `http://127.0.0.1:${port}/utf8`,
            max_body_bytes: 1,
            timeout_ms: 2_000,
        },
    });
    assert.equal(splitUtf8.value.body, "");
    assert.equal(splitUtf8.value.bodyBytes, 1);
    assert.equal(splitUtf8.value.bodyTruncated, true);

    const invalidUtf8 = await runtime.invoke({
        scope: "project",
        packageName: "network_probe",
        tool: "http_probe",
        arguments: {
            url: `http://127.0.0.1:${port}/invalid-utf8`,
            max_body_bytes: 32,
            timeout_ms: 2_000,
        },
    });
    assert.equal(invalidUtf8.value.body, "\ufffd");
    assert.equal(invalidUtf8.value.bodyBytes, 1);
    assert.equal(invalidUtf8.value.bodyTruncated, false);

    const dns = await runtime.invoke({
        scope: "project",
        packageName: "network_probe",
        tool: "dns_lookup",
        arguments: {
            host: "localhost",
            rrtype: "A",
            timeout_ms: 1_000,
        },
    });
    assert.equal(dns.value.host, "localhost");
    assert.equal(typeof dns.value.ok, "boolean");
});

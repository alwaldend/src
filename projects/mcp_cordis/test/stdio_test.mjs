import assert from "node:assert/strict";
import { mkdir, mkdtemp, writeFile } from "node:fs/promises";
import { isAbsolute, join, resolve } from "node:path";
import test from "node:test";
import { Client } from "@modelcontextprotocol/client";
import { StdioClientTransport } from "@modelcontextprotocol/client/stdio";

const temporaryRoot = process.env.TEST_TMPDIR;
assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");

function source(generation) {
    return `export default {
      description: "stdio ${generation}",
      apply(ctx) {
        console.log("plugin apply log ${generation}");
        ctx.tool({
          name: "stdio_echo",
          description: "Return the stdio test generation.",
          inputSchema: {
            type: "object",
            properties: { value: {} },
            additionalProperties: false
          }
        }, ({ value = null }) => {
          process.stdout.write("plugin invoke log ${generation}\\n");
          return {
            generation: ${JSON.stringify(generation)},
            value
          };
        });
      }
    };
`;
}

async function call(client, name, args) {
    const result = await client.callTool({ name, arguments: args });
    const value = result.structuredContent ??
        JSON.parse(result.content[0].text);
    if (result.isError) {
        throw new Error(`${name} failed: ${JSON.stringify(value)}`);
    }
    return value;
}

async function waitFor(callback, timeoutMs = 5_000) {
    const deadline = Date.now() + timeoutMs;
    let lastError;
    while (Date.now() < deadline) {
        try {
            const value = await callback();
            if (value !== undefined) return value;
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

async function connect(binary, workspaceRoot) {
    const transport = new StdioClientTransport({
        command: binary,
        args: ["--workspace-root", workspaceRoot],
        stderr: "pipe",
    });
    let stderr = "";
    transport.stderr.on("data", (chunk) => {
        if (stderr.length < 65_536) stderr += chunk.toString("utf8");
    });
    const client = new Client({
        name: "mcp-cordis-stdio-test",
        version: "1.0.0",
    });
    try {
        await client.connect(transport);
    } catch (error) {
        throw new Error(
            `stdio connection failed: ${error.message}\n${stderr}`,
        );
    }
    return { client, transport, stderr: () => stderr };
}

test("real stdio connection hot-updates and recovers project state", async (t) => {
    const rawBinary = process.argv[2];
    assert.ok(rawBinary, "the Bazel target must pass the server launcher");
    const binary = isAbsolute(rawBinary) ? rawBinary : resolve(rawBinary);
    const workspaceRoot = await mkdtemp(join(temporaryRoot, "stdio-"));
    const projectRoot = join(workspaceRoot, "projects", "mcp_cordis");
    await mkdir(join(projectRoot, "plugins"), { recursive: true });
    await writeFile(join(projectRoot, "cordis.yaml"), "[]\n");

    const firstConnection = await connect(binary, workspaceRoot);
    let restarted;
    t.after(async () => {
        await restarted?.client.close().catch(() => {});
        await restarted?.transport.close().catch(() => {});
        await firstConnection.client.close().catch(() => {});
        await firstConnection.transport.close().catch(() => {});
    });
    const firstPid = firstConnection.transport.pid;
    assert.ok(firstPid);
    const listed = await firstConnection.client.listTools();
    assert.equal(listed.tools.length, 10);
    const toolsByName = new Map(listed.tools.map((tool) => [tool.name, tool]));
    assert.equal(
        toolsByName.get("cordis_define")?.annotations?.destructiveHint,
        true,
    );
    assert.equal(
        toolsByName.get("cordis_promote")?.annotations?.destructiveHint,
        true,
    );

    await call(firstConnection.client, "cordis_define", {
        scope: "project",
        name: "stdio_package",
        source: source("v1"),
        activate: false,
    });
    const stored = await call(firstConnection.client, "cordis_inspect", {
        scope: "project",
        name: "stdio_package",
        include_source: true,
    });
    assert.equal(stored.running, false);
    assert.match(stored.source, /stdio v1/u);
    await call(firstConnection.client, "cordis_run", {
        scope: "project",
        name: "stdio_package",
    });
    const first = await call(firstConnection.client, "cordis_invoke", {
        scope: "project",
        package: "stdio_package",
        tool: "stdio_echo",
        arguments: { value: 1 },
    });
    assert.deepEqual(first.value, { generation: "v1", value: 1 });

    const update = await call(
        firstConnection.client,
        "cordis_define",
        {
            scope: "project",
            name: "stdio_package",
            source: source("v2"),
        },
    );
    assert.equal(update.persisted, true);
    assert.equal(update.activation, "pending");
    const second = await waitFor(async () => {
        const value = await call(firstConnection.client, "cordis_invoke", {
            scope: "project",
            package: "stdio_package",
            tool: "stdio_echo",
            arguments: { value: 2 },
        });
        return value.value.generation === "v2" ? value : undefined;
    });
    assert.deepEqual(second.value, { generation: "v2", value: 2 });
    assert.equal(firstConnection.transport.pid, firstPid);

    await assert.rejects(
        call(firstConnection.client, "cordis_define", {
            scope: "project",
            name: "stdio_package",
            source: "export default {",
        }),
        /cordis_define failed/u,
    );
    const afterInvalidSource = await call(
        firstConnection.client,
        "cordis_invoke",
        {
        scope: "project",
        package: "stdio_package",
        tool: "stdio_echo",
        arguments: { value: "rollback" },
        },
    );
    assert.equal(afterInvalidSource.value.generation, "v2");

    await call(firstConnection.client, "cordis_stop", {
        scope: "project",
        name: "stdio_package",
    });
    const stopped = await call(firstConnection.client, "cordis_inspect", {
        scope: "project",
        name: "stdio_package",
    });
    assert.equal(stopped.running, false);
    await call(firstConnection.client, "cordis_remove", {
        scope: "project",
        name: "stdio_package",
    });
    const removed = await call(firstConnection.client, "cordis_list", {
        scope: "project",
    });
    assert.equal(removed.packages.length, 0);

    await call(firstConnection.client, "cordis_define", {
        scope: "project",
        name: "persistent_package",
        source: source("v2"),
        activate: true,
    });
    await firstConnection.client.close();
    await firstConnection.transport.close();

    restarted = await connect(binary, workspaceRoot);
    assert.notEqual(restarted.transport.pid, firstPid);
    const recovered = await call(restarted.client, "cordis_invoke", {
        scope: "project",
        package: "persistent_package",
        tool: "stdio_echo",
        arguments: { value: "restart" },
    });
    assert.deepEqual(recovered.value, {
        generation: "v2",
        value: "restart",
    });
    await restarted.client.close();
    await restarted.transport.close();
});

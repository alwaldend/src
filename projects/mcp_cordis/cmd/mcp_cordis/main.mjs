import { realpath } from "node:fs/promises";
import { isAbsolute, resolve } from "node:path";
import { Writable } from "node:stream";
import {
    serveStdio,
    StdioServerTransport,
} from "@modelcontextprotocol/server/stdio";
import { createMcpServer } from "../../internal/mcp.mjs";
import { CordisRuntime } from "../../internal/runtime.mjs";

function usage() {
    return [
        "Usage: mcp_cordis --workspace-root PATH [options]",
        "",
        "Options:",
        "  --invoke-timeout-ms N",
        "  --max-output-bytes N",
        "  --help",
    ].join("\n");
}

function positiveInteger(flag, raw) {
    const value = Number(raw);
    if (!Number.isSafeInteger(value) || value <= 0) {
        throw new Error(`${flag} requires a positive integer`);
    }
    return value;
}

function parseArguments(argv) {
    const options = {};
    for (let index = 0; index < argv.length; index += 1) {
        const flag = argv[index];
        if (flag === "--help") return { help: true };
        const value = argv[index + 1];
        if (value === undefined) throw new Error(`${flag} requires a value`);
        index += 1;
        if (flag === "--workspace-root") {
            options.workspaceRoot = value;
        } else if (flag === "--invoke-timeout-ms") {
            options.invokeTimeoutMs = positiveInteger(flag, value);
        } else if (flag === "--max-output-bytes") {
            options.maxOutputBytes = positiveInteger(flag, value);
        } else {
            throw new Error(`unknown option: ${flag}`);
        }
    }
    return options;
}

async function canonicalWorkspaceRoot(candidate) {
    if (!candidate) {
        throw new Error(
            "--workspace-root is required when " +
                "BUILD_WORKSPACE_DIRECTORY is unavailable",
        );
    }
    const absolute = isAbsolute(candidate) ? candidate : resolve(candidate);
    return realpath(absolute);
}

async function main() {
    const options = parseArguments(process.argv.slice(2));
    if (options.help) {
        process.stdout.write(`${usage()}\n`);
        return;
    }
    options.workspaceRoot = await canonicalWorkspaceRoot(
        options.workspaceRoot ?? process.env.BUILD_WORKSPACE_DIRECTORY,
    );

    // Keep the protocol on a private stream. Runtime plugins share this
    // process, so redirect their accidental stdout writes away from the
    // JSON-RPC wire before any plugin is activated.
    const protocolWrite = process.stdout.write.bind(process.stdout);
    const protocolOutput = new Writable({
        write(chunk, encoding, callback) {
            protocolWrite(chunk, encoding, callback);
        },
    });
    process.stdout.write = process.stderr.write.bind(process.stderr);

    const runtime = new CordisRuntime(options);
    let startup;
    try {
        startup = await runtime.initialize();
    } catch (error) {
        await runtime.shutdown().catch(() => {});
        throw error;
    }
    process.stderr.write(
        `[mcp_cordis] workspace=${options.workspaceRoot} ` +
            `loaded=${startup.loaded.length} failed=${startup.errors.length}\n`,
    );

    const handle = serveStdio(() => createMcpServer(runtime), {
        transport: new StdioServerTransport(process.stdin, protocolOutput),
        onerror: (error) => {
            process.stderr.write(
                `[mcp_cordis] MCP transport error: ${error.message}\n`,
            );
        },
    });

    let closing = false;
    const close = async (signal) => {
        if (closing) return;
        closing = true;
        process.stderr.write(`[mcp_cordis] closing on ${signal}\n`);
        await handle.close().catch((error) => {
            process.stderr.write(
                `[mcp_cordis] transport close failed: ${error.message}\n`,
            );
        });
        await runtime.shutdown();
    };
    process.once("SIGINT", () => void close("SIGINT"));
    process.once("SIGTERM", () => void close("SIGTERM"));
    process.stdin.once("end", () => void close("stdin EOF"));
    process.stdin.once("close", () => void close("stdin close"));
}

main().catch((error) => {
    process.stderr.write(
        `[mcp_cordis] fatal: ${error instanceof Error ? error.stack : error}\n`,
    );
    process.exitCode = 1;
});

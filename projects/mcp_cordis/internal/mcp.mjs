import { McpServer } from "@modelcontextprotocol/server";
import { z } from "zod";

const nameSchema = z.string().regex(/^[a-z][a-z0-9_]{0,63}$/);
const scopeSchema = z.enum(["project", "scratch"]);
const argumentsSchema = z.record(z.string(), z.unknown());

function toolResult(value) {
    return {
        content: [{ type: "text", text: JSON.stringify(value, null, 2) }],
        structuredContent: value,
    };
}

function toolError(error) {
    const value = {
        error: {
            code: error?.code ?? "runtime_error",
            message: error instanceof Error ? error.message : String(error),
            ...(error?.details === undefined
                ? {}
                : { details: error.details }),
        },
    };
    return { ...toolResult(value), isError: true };
}

function handler(callback) {
    return async (input) => {
        try {
            return toolResult(await callback(input));
        } catch (error) {
            return toolError(error);
        }
    };
}

export function createMcpServer(runtime) {
    const server = new McpServer(
        { name: "mcp-cordis", version: "0.1.0" },
        {
            instructions: [
                "Use cordis_list_tools to discover live package handlers.",
                "Call them through cordis_invoke immediately; no MCP reload",
                "or client restart is needed. Every package identity includes",
                "an explicit project or scratch scope.",
            ].join(" "),
        },
    );
    registerMcpTools(server, runtime);
    return server;
}

export function registerMcpTools(server, runtime) {
    server.registerTool(
        "cordis_list",
        {
            description: "List standard Cordis entries and their live state.",
            inputSchema: z.object({
                scope: scopeSchema.optional(),
            }),
            annotations: { readOnlyHint: true },
        },
        handler(({ scope }) => runtime.listPackages({ scope })),
    );

    server.registerTool(
        "cordis_inspect",
        {
            description: "Inspect one scoped Cordis entry and normal module.",
            inputSchema: z.object({
                scope: scopeSchema,
                name: nameSchema,
                include_source: z.boolean().default(false),
            }),
            annotations: { readOnlyHint: true },
        },
        handler(({ scope, name, include_source: includeSource }) => {
            return runtime.inspect({ scope, name, includeSource });
        }),
    );

    server.registerTool(
        "cordis_define",
        {
            description: [
                "Atomically persist a normal ESM Cordis plugin module.",
                "A newly created active entry loads through Cordis Include",
                "before this call returns. Re-enabling a stored entry first",
                "refreshes its cached module. Updates to a running entry",
                "return activation=pending during the live swap.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema.default("scratch"),
                name: nameSchema,
                source: z.string().min(1).max(2_000_256),
                activate: z.boolean().optional(),
            }),
            annotations: { destructiveHint: true },
        },
        handler((input) => runtime.define(input)),
    );

    server.registerTool(
        "cordis_run",
        {
            description: [
                "Enable a stored Cordis entry after refreshing any cached",
                "module through Cordis HMR. Activation finishes before this",
                "call returns.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema,
                name: nameSchema,
            }),
            annotations: { destructiveHint: false },
        },
        handler((input) => runtime.run(input)),
    );

    server.registerTool(
        "cordis_reload",
        {
            description: [
                "Atomically rewrite the entry's current module bytes so",
                "official Cordis HMR performs an eventual live reload.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema,
                name: nameSchema,
            }),
            annotations: { destructiveHint: false },
        },
        handler((input) => runtime.reload(input)),
    );

    server.registerTool(
        "cordis_list_tools",
        {
            description: [
                "List handlers currently registered by live Cordis package",
                "Fibers. Use their scope/package/name with cordis_invoke.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema.optional(),
                package: nameSchema.optional(),
            }),
            annotations: { readOnlyHint: true },
        },
        handler(({ scope, package: packageName }) => {
            return runtime.listTools({ scope, packageName });
        }),
    );

    server.registerTool(
        "cordis_invoke",
        {
            description: [
                "Invoke a live package handler over this same MCP connection.",
                "A catalog version can guard against stale discoveries.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema,
                package: nameSchema,
                tool: nameSchema,
                arguments: argumentsSchema.default({}),
                timeout_ms: z.number().int().min(1).max(300_000).optional(),
                catalog_version: z.number().int().min(1).optional(),
            }),
        },
        handler(
            ({
                scope,
                package: packageName,
                tool,
                arguments: args,
                timeout_ms: timeoutMs,
                catalog_version: catalogVersion,
            }) => {
                return runtime.invoke({
                    scope,
                    packageName,
                    tool,
                    arguments: args,
                    timeoutMs,
                    catalogVersion,
                });
            },
        ),
    );

    server.registerTool(
        "cordis_stop",
        {
            description: [
                "Disable one live Cordis entry while retaining its module.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema,
                name: nameSchema,
            }),
            annotations: { destructiveHint: false },
        },
        handler((input) => runtime.stop(input)),
    );

    server.registerTool(
        "cordis_remove",
        {
            description: [
                "Stop and permanently remove one scoped Cordis entry and",
                "its normal module file.",
            ].join(" "),
            inputSchema: z.object({
                scope: scopeSchema,
                name: nameSchema,
            }),
            annotations: { destructiveHint: true },
        },
        handler((input) => runtime.remove(input)),
    );

    server.registerTool(
        "cordis_promote",
        {
            description: [
                "Copy a scratch Cordis module into the reusable project",
                "config, optionally activating the promoted entry.",
            ].join(" "),
            inputSchema: z.object({
                name: nameSchema,
                target_name: nameSchema.optional(),
                activate: z.boolean().default(false),
            }),
            annotations: { destructiveHint: true },
        },
        handler(({ name, target_name: targetName, activate }) => {
            return runtime.promote({
                name,
                targetName,
                activate,
            });
        }),
    );
}

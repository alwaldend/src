import { createHash, randomUUID } from "node:crypto";
import { AsyncLocalStorage } from "node:async_hooks";
import {
    mkdir,
    open,
    readFile,
    rename,
    rm,
    writeFile,
} from "node:fs/promises";
import * as path from "node:path";
import process from "node:process";
import { fileURLToPath, pathToFileURL } from "node:url";
import { SourceTextModule } from "node:vm";
import { Context } from "@deepseek-ai/cordis";
import Hmr from "@deepseek-ai/cordis-plugin-hmr";
import Include, { entryListSchema } from "@deepseek-ai/cordis-plugin-include";
import Loader from "@deepseek-ai/cordis-plugin-loader";
import Timer from "@deepseek-ai/cordis-plugin-timer";
import Ajv2020 from "ajv/dist/2020.js";
import * as yaml from "js-yaml";
import { ProcessSupervisor } from "./process_supervisor.mjs";

const MAX_SOURCE_BYTES = 2_000_000;
const MAX_READ_BYTES = 64 * 1024 * 1024;
const NAME_PATTERN = /^[a-z][a-z0-9_]{0,63}$/u;
const TOOL_NAME_PATTERN = NAME_PATTERN;
const OWNERSHIP_ID_PATTERN = /^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$/u;
const LEGACY_SOURCE_MARKER_PATTERN =
    /^(#![^\n]*\n)?export const __mcp_cordis_source_sha256 = "[a-f0-9]{64}";\n/u;

export class RuntimeError extends Error {
    constructor(code, message, details = undefined) {
        super(message);
        this.name = "RuntimeError";
        this.code = code;
        this.details = details;
    }
}

function codedError(code, message, ErrorType = Error) {
    const error = new ErrorType(message);
    error.code = code;
    return error;
}

function plainError(error) {
    return {
        code: error?.code ?? "runtime_error",
        message: error instanceof Error ? error.message : String(error),
        ...(error?.details === undefined ? {} : { details: error.details }),
    };
}

function deferred() {
    let resolve;
    let reject;
    const promise = new Promise((resolvePromise, rejectPromise) => {
        resolve = resolvePromise;
        reject = rejectPromise;
    });
    return { promise, reject, resolve };
}

function packageKey(scope, name) {
    return `${validateScope(scope)}:${validateName(name)}`;
}

function validateScope(scope) {
    if (scope !== "project" && scope !== "scratch") {
        throw new TypeError("scope must be project or scratch");
    }
    return scope;
}

function validateName(name) {
    if (typeof name !== "string" || !NAME_PATTERN.test(name)) {
        throw new TypeError("name must match [a-z][a-z0-9_]{0,63}");
    }
    return name;
}

function validateOwnershipID(field, value) {
    if (typeof value !== "string" || !OWNERSHIP_ID_PATTERN.test(value)) {
        throw new TypeError(`${field} must be a lowercase portable identity`);
    }
    return value;
}

function validateSource(source) {
    if (typeof source !== "string" || source.length === 0) {
        throw new TypeError("source must be a non-empty string");
    }
    if (Buffer.byteLength(source) > MAX_SOURCE_BYTES) {
        throw new RangeError(
            `source exceeds the ${MAX_SOURCE_BYTES}-byte limit`,
        );
    }
    return source.endsWith("\n") ? source : `${source}\n`;
}

function validateModuleSource(source, identifier = "runtime module") {
    let module;
    try {
        module = new SourceTextModule(source, { identifier });
    } catch (error) {
        throw new RuntimeError(
            "invalid_module_source",
            `source is not a valid ESM module: ${error.message}`,
        );
    }
}

function prepareSource(source) {
    if (typeof source !== "string" || source.length === 0) {
        throw new TypeError("source must be a non-empty string");
    }
    const legacyMarker = source.match(LEGACY_SOURCE_MARKER_PATTERN);
    if (legacyMarker) {
        source =
            `${legacyMarker[1] ?? ""}` + source.slice(legacyMarker[0].length);
    }
    source = validateSource(source);
    validateModuleSource(source);
    return source;
}

function positiveInteger(value, fallback, label, maximum = Infinity) {
    if (value === undefined) return fallback;
    if (!Number.isSafeInteger(value) || value <= 0 || value > maximum) {
        throw new TypeError(`${label} must be a positive integer`);
    }
    return value;
}

function resolveInsideWorkspace(
    workspaceRoot,
    relativePath = ".",
    { allowAbsolute = false } = {},
) {
    if (typeof relativePath !== "string" || relativePath.includes("\0")) {
        throw new TypeError(
            "workspace path must be a string without NUL bytes",
        );
    }
    if (path.isAbsolute(relativePath) && !allowAbsolute) {
        throw new TypeError("workspace path must be relative");
    }

    const resolved = path.resolve(workspaceRoot, relativePath);
    const relative = path.relative(workspaceRoot, resolved);
    if (
        relative === ".." ||
        relative.startsWith(`..${path.sep}`) ||
        path.isAbsolute(relative)
    ) {
        throw new RangeError("workspace path escapes the workspace root");
    }
    return resolved;
}

function isInside(root, candidate) {
    const relative = path.relative(root, candidate);
    return (
        relative === "" ||
        (relative !== ".." &&
            !relative.startsWith(`..${path.sep}`) &&
            !path.isAbsolute(relative))
    );
}

async function readTextBounded(filename, maxBytes) {
    const handle = await open(filename, "r");
    try {
        const buffer = Buffer.allocUnsafe(maxBytes + 1);
        let offset = 0;
        while (offset < buffer.length) {
            const { bytesRead } = await handle.read(
                buffer,
                offset,
                buffer.length - offset,
                offset,
            );
            if (!bytesRead) break;
            offset += bytesRead;
        }
        if (offset > maxBytes) {
            throw new RangeError(
                `file exceeds the ${maxBytes}-byte read limit`,
            );
        }
        return buffer.subarray(0, offset).toString("utf8");
    } finally {
        await handle.close();
    }
}

function jsonClone(value, label) {
    if (value === undefined) return null;
    let encoded;
    try {
        encoded = JSON.stringify(value, (_key, part) => {
            const type = typeof part;
            if (
                type === "bigint" ||
                type === "function" ||
                type === "symbol" ||
                type === "undefined"
            ) {
                throw new TypeError(`${label} is not JSON-serializable`);
            }
            return part;
        });
    } catch (error) {
        throw new TypeError(
            `${label} is not JSON-serializable: ${error.message}`,
            { cause: error },
        );
    }
    if (encoded === undefined) {
        throw new TypeError(`${label} is not JSON-serializable`);
    }
    return JSON.parse(encoded);
}

function normalizeArguments(args) {
    if (args === undefined || args === null) return {};
    if (typeof args !== "object" || Array.isArray(args)) {
        throw new TypeError("tool arguments must be an object");
    }
    return args;
}

function formatValidationErrors(errors) {
    return (errors ?? [])
        .map((error) => {
            const location = error.instancePath || "/";
            return `${location} ${error.message ?? "is invalid"}`;
        })
        .join("; ");
}

function normalizeToolDefinition(definition, handler, ajv) {
    if (
        !definition ||
        typeof definition !== "object" ||
        Array.isArray(definition)
    ) {
        throw new TypeError("tool definition must be an object");
    }
    if (typeof handler !== "function") {
        throw new TypeError("tool handler must be a function");
    }
    const { name } = definition;
    if (typeof name !== "string" || !TOOL_NAME_PATTERN.test(name)) {
        throw new TypeError(
            "tool definition.name must match [a-z][a-z0-9_]{0,63}",
        );
    }
    if (
        definition.description !== undefined &&
        typeof definition.description !== "string"
    ) {
        throw new TypeError("tool definition.description must be a string");
    }
    if (
        !definition.inputSchema ||
        typeof definition.inputSchema !== "object" ||
        Array.isArray(definition.inputSchema)
    ) {
        throw new TypeError("tool definition.inputSchema must be an object");
    }

    const inputSchema = jsonClone(
        definition.inputSchema,
        `input schema for ${name}`,
    );
    let validate;
    try {
        validate = ajv.compile(inputSchema);
    } catch (error) {
        throw new TypeError(
            `invalid JSON Schema for tool ${name}: ${error.message}`,
            { cause: error },
        );
    }
    return {
        handler,
        validate,
        metadata: {
            name,
            description: definition.description ?? "",
            inputSchema,
        },
    };
}

async function atomicWrite(filename, content) {
    await mkdir(path.dirname(filename), { recursive: true });
    const temporary = path.join(
        path.dirname(filename),
        `.${path.basename(filename)}.${process.pid}.${randomUUID()}.tmp`,
    );
    try {
        await writeFile(temporary, content, {
            encoding: "utf8",
            flag: "wx",
            mode: 0o644,
        });
        await rename(temporary, filename);
    } catch (error) {
        await rm(temporary, { force: true }).catch(() => {});
        throw error;
    }
}

async function readMaybe(filename) {
    try {
        return await readFile(filename, "utf8");
    } catch (error) {
        if (error?.code === "ENOENT") return undefined;
        throw error;
    }
}

function hashSource(source) {
    return createHash("sha256").update(source).digest("hex");
}

function managedSpecifier(name) {
    return `./plugins/${validateName(name)}.mjs`;
}

function parseEntries(content, filename) {
    let value;
    try {
        value = yaml.load(content, { schema: entryListSchema });
    } catch (error) {
        throw new RuntimeError(
            "invalid_cordis_config",
            `failed to parse ${filename}: ${error.message}`,
        );
    }
    if (!Array.isArray(value)) {
        throw new RuntimeError(
            "invalid_cordis_config",
            `${filename} must contain a top-level entry list`,
        );
    }
    return value;
}

function dumpEntries(entries) {
    return yaml.dump(entries, {
        schema: entryListSchema,
        indent: 2,
        lineWidth: 79,
        noRefs: true,
        sortKeys: false,
    });
}

function withTimeout(promise, timeoutMs, error) {
    let timer;
    const timeout = new Promise((_resolve, reject) => {
        timer = setTimeout(() => reject(error), timeoutMs);
        timer.unref();
    });
    return Promise.race([promise, timeout]).finally(() => {
        clearTimeout(timer);
    });
}

export class CordisRuntime {
    #admissions = new Set();
    #ajv = new Ajv2020({ allErrors: true, strict: true });
    #catalogVersion = 1;
    #closed = false;
    #includes = new Map();
    #initialized = false;
    #invocationSupervisors = new AsyncLocalStorage();
    #mutation = Promise.resolve();
    #shutdownPromise;
    #scopeErrors = new Map();
    #supervisors = new WeakMap();
    #tools = new Map();

    constructor({
        workspaceRoot,
        projectRoot = path.join(workspaceRoot, "projects", "mcp_cordis"),
        taskId = "test-task",
        runId = "test-run",
        workerId = "test-worker",
        invokeTimeoutMs = 30_000,
        maxOutputBytes = 1_048_576,
    }) {
        this.workspaceRoot = path.resolve(workspaceRoot);
        this.projectRoot = path.resolve(projectRoot);
        this.taskId = validateOwnershipID("taskId", taskId);
        this.runId = validateOwnershipID("runId", runId);
        this.workerId = validateOwnershipID("workerId", workerId);
        this.scratchRoot = path.join(
            this.workspaceRoot,
            "out",
            this.taskId,
            "mcp_cordis",
            "runs",
            this.runId,
        );
        if (!isInside(this.workspaceRoot, this.projectRoot)) {
            throw new RangeError("projectRoot must be inside workspaceRoot");
        }
        this.invokeTimeoutMs = positiveInteger(
            invokeTimeoutMs,
            30_000,
            "invokeTimeoutMs",
        );
        this.maxOutputBytes = positiveInteger(
            maxOutputBytes,
            1_048_576,
            "maxOutputBytes",
            MAX_READ_BYTES,
        );
        this.configFiles = {
            project: path.join(this.projectRoot, "cordis.yaml"),
            scratch: path.join(this.scratchRoot, "cordis.yaml"),
        };
        this.pluginRoots = {
            project: path.join(this.projectRoot, "plugins"),
            scratch: path.join(this.scratchRoot, "plugins"),
        };
    }

    get catalogVersion() {
        return this.#catalogVersion;
    }

    initialize() {
        this.#assertOpen();
        return this.#admit(() => this.#initializeInternal());
    }

    async #initializeInternal() {
        if (this.#initialized) {
            return {
                loaded: this.#loadedPackages(),
                errors: [...this.#scopeErrors].map(([scope, error]) => ({
                    scope,
                    error: structuredClone(error),
                })),
            };
        }
        await mkdir(this.pluginRoots.scratch, { recursive: true });
        await atomicWrite(
            path.join(this.scratchRoot, "manifest.json"),
            `${JSON.stringify({
                apiVersion: "agents.alwaldend.com/v1alpha1",
                kind: "TaskRunManifest",
                taskId: this.taskId,
                runId: this.runId,
                workerId: this.workerId,
                information: {
                    public: true,
                    secret: true,
                    personal: true,
                },
                budget: {
                    calls: 1,
                    bytes: this.maxOutputBytes,
                    durationMs: this.invokeTimeoutMs,
                    concurrency: 1,
                },
                retention: "task",
                lockScope: `${this.taskId}/${this.runId}`,
                cleanupOwner: "task-owner",
            })}\n`,
        );
        if ((await readMaybe(this.configFiles.scratch)) === undefined) {
            await atomicWrite(this.configFiles.scratch, "[]\n");
        }

        this.root = new Context();
        this.root.baseUrl = pathToFileURL(
            `${this.workspaceRoot}${path.sep}`,
        ).href;
        this.#installHostApi();

        await this.#mount(Loader, { baseUrl: this.root.baseUrl });
        await this.#mount(Timer);
        await this.#mount(Hmr, {
            base: this.workspaceRoot,
            root: [
                path.relative(this.workspaceRoot, this.configFiles.project),
                path.relative(this.workspaceRoot, this.pluginRoots.project),
                path.relative(this.workspaceRoot, this.configFiles.scratch),
                path.relative(this.workspaceRoot, this.pluginRoots.scratch),
            ],
            ignored: ["**/node_modules", "**/.*"],
        });
        this.root.loader.builtins.include = Include;

        const errors = [];
        for (const scope of ["project", "scratch"]) {
            try {
                await this.#validateConfigSources(scope);
                const id = await this.root.loader.create({
                    id: `${scope}_entries`,
                    name: "cordis:include",
                    config: {
                        path: pathToFileURL(this.configFiles[scope]).href,
                        initial: [],
                        enableLogs: true,
                    },
                });
                const include = this.root.loader.resolve(id).subtree;
                if (!(include instanceof Include)) {
                    throw new Error(`${scope} include did not expose a tree`);
                }
                this.#includes.set(scope, include);
            } catch (error) {
                const failure = { scope, error: plainError(error) };
                errors.push(failure);
                this.#scopeErrors.set(scope, structuredClone(failure.error));
                process.stderr.write(
                    `[mcp_cordis] ${scope} Cordis config failed: ` +
                        `${JSON.stringify(failure.error)}\n`,
                );
            }
        }
        await this.root.loader.await();
        this.#initialized = true;
        return { loaded: this.#loadedPackages(), errors };
    }

    async #mount(plugin, config = undefined) {
        const fiber = this.root.plugin(plugin, config);
        await fiber.await();
        return fiber;
    }

    #installHostApi() {
        const runtime = this;
        this.root.workspaceRoot = this.workspaceRoot;
        this.root.resolveWorkspace = function (relativePath = ".") {
            return resolveInsideWorkspace(runtime.workspaceRoot, relativePath);
        };
        this.root.readText = function (relativePath, options = {}) {
            if (
                !options ||
                typeof options !== "object" ||
                Array.isArray(options)
            ) {
                throw new TypeError("readText options must be an object");
            }
            const maxBytes = positiveInteger(
                options.maxBytes,
                runtime.maxOutputBytes,
                "readText options.maxBytes",
                MAX_READ_BYTES,
            );
            return readTextBounded(
                resolveInsideWorkspace(runtime.workspaceRoot, relativePath),
                maxBytes,
            );
        };
        this.root.exec = function (file, args = [], options = {}) {
            return runtime.#execute(this, file, args, options);
        };
        this.root.tool = function (definition, handler) {
            return runtime.#registerTool(this, definition, handler);
        };
    }

    #owningEntry(context) {
        let fiber = context.fiber;
        while (fiber) {
            if (fiber.entry) return fiber.entry;
            const next = fiber.parent?.fiber;
            if (!next || next === fiber) break;
            fiber = next;
        }
        throw new RuntimeError(
            "unmanaged_plugin",
            "host helpers require a Cordis loader-managed plugin",
        );
    }

    #identity(context) {
        const entry = this.#owningEntry(context);
        const base = path.resolve(
            fileURLToPath(entry.parent.tree.ctx.baseUrl),
        );
        let scope;
        if (base === this.projectRoot) {
            scope = "project";
        } else if (base === this.scratchRoot) {
            scope = "scratch";
        } else {
            throw new RuntimeError(
                "unmanaged_plugin",
                `loader entry ${entry.id} is outside managed Cordis configs`,
            );
        }
        // Entry.id is qualified by owning Include entries (for example,
        // "scratch_entries:echo"). The persisted child id is the package
        // identity inside this scoped config.
        const name = validateName(entry.options.id);
        return { entry, key: packageKey(scope, name), name, scope };
    }

    #registerTool(context, definition, handler) {
        const identity = this.#identity(context);
        const record = {
            ...normalizeToolDefinition(definition, handler, this.#ajv),
            context,
            name: identity.name,
            scope: identity.scope,
        };
        let packageTools = this.#tools.get(identity.key);
        if (!packageTools) {
            packageTools = new Map();
            this.#tools.set(identity.key, packageTools);
        }
        const toolName = record.metadata.name;
        context.effect(
            () => {
                if (packageTools.has(toolName)) {
                    throw new Error(`tool ${toolName} is already registered`);
                }
                packageTools.set(toolName, record);
                this.#catalogVersion += 1;
                return () => {
                    if (packageTools.get(toolName) !== record) return;
                    packageTools.delete(toolName);
                    if (packageTools.size === 0) {
                        this.#tools.delete(identity.key);
                    }
                    this.#catalogVersion += 1;
                };
            },
            `mcp_cordis.tool(${JSON.stringify(toolName)})`,
        );
    }

    #supervisor(context) {
        const fiber = context.fiber;
        let supervisor = this.#supervisors.get(fiber);
        if (supervisor) return supervisor;
        supervisor = new ProcessSupervisor();
        this.#supervisors.set(fiber, supervisor);
        context.effect(
            () => () =>
                supervisor.close(
                    codedError(
                        "EXEC_DISPOSED",
                        "ctx.exec was cancelled during Cordis plugin disposal",
                    ),
                ),
            "mcp_cordis.host_processes",
        );
        return supervisor;
    }

    #execute(context, file, args, options) {
        this.#identity(context);
        if (typeof file !== "string" || !file || file.includes("\0")) {
            throw new TypeError("exec file must be a non-empty string");
        }
        if (
            !Array.isArray(args) ||
            args.some((argument) => typeof argument !== "string")
        ) {
            throw new TypeError("exec args must be an array of strings");
        }
        if (
            !options ||
            typeof options !== "object" ||
            Array.isArray(options)
        ) {
            throw new TypeError("exec options must be an object");
        }
        if (
            options.allowTruncatedOutput !== undefined &&
            typeof options.allowTruncatedOutput !== "boolean"
        ) {
            throw new TypeError(
                "exec options.allowTruncatedOutput must be a boolean",
            );
        }
        if (process.platform !== "linux") {
            throw codedError(
                "EXEC_UNSUPPORTED_PLATFORM",
                "ctx.exec requires Linux process-group supervision",
            );
        }

        const cwd = resolveInsideWorkspace(
            this.workspaceRoot,
            options.cwd ?? ".",
            { allowAbsolute: true },
        );
        const maxBytes = positiveInteger(
            options.maxBytes,
            this.maxOutputBytes,
            "exec options.maxBytes",
            MAX_READ_BYTES,
        );
        const timeoutMs = positiveInteger(
            options.timeoutMs,
            30_000,
            "exec options.timeoutMs",
            60 * 60 * 1000,
        );
        if (
            options.env !== undefined &&
            (!options.env ||
                typeof options.env !== "object" ||
                Array.isArray(options.env))
        ) {
            throw new TypeError("exec options.env must be an object");
        }
        const env = { ...process.env };
        for (const [key, value] of Object.entries(options.env ?? {})) {
            if (value === undefined || value === null) {
                delete env[key];
            } else if (typeof value === "string") {
                env[key] = value;
            } else {
                throw new TypeError(
                    `exec environment variable ${key} must be a string`,
                );
            }
        }
        const supervisor =
            this.#invocationSupervisors.getStore() ??
            this.#supervisor(context);
        return supervisor.execute({
            file,
            args,
            options: {
                cwd,
                env,
                maxBytes,
                timeoutMs,
                allowTruncatedOutput: options.allowTruncatedOutput === true,
            },
        });
    }

    async define({ scope = "scratch", name, source, activate = undefined }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(name);
        return this.#mutate(() =>
            this.#defineUnlocked({
                scope,
                name,
                source,
                activate,
            }),
        );
    }

    async #defineUnlocked({ scope, name, source, activate }) {
        source = prepareSource(source);
        const state = await this.#configState(scope);
        const index = state.entries.findIndex((entry) => entry.id === name);
        const existing = index < 0 ? undefined : state.entries[index];
        if (existing && existing.name !== managedSpecifier(name)) {
            throw new RuntimeError(
                "unmanaged_entry",
                `${packageKey(scope, name)} does not use its managed module path`,
            );
        }
        const filename = this.#sourceFile(scope, name);
        let sourceChanged = false;
        let reloadPending = false;

        if (existing) {
            const replacement = await this.#replaceSource(scope, name, source);
            sourceChanged = replacement.changed;
            if (activate === true && existing.disabled) {
                await this.root.hmr.refreshFile(filename);
                const entries = structuredClone(state.entries);
                delete entries[index].disabled;
                await this.#commitConfig(scope, state.content, entries);
            }
            reloadPending = sourceChanged && !existing.disabled;
        } else {
            await atomicWrite(filename, source);
            sourceChanged = true;
            const entries = structuredClone(state.entries);
            entries.push({
                id: name,
                name: managedSpecifier(name),
                ...(activate === true ? {} : { disabled: true }),
            });
            try {
                await this.#commitConfig(scope, state.content, entries);
            } catch (error) {
                try {
                    await rm(filename, { force: true });
                } catch (cleanupError) {
                    throw new RuntimeError(
                        "define_cleanup_failed",
                        `failed to define ${packageKey(scope, name)} and ` +
                            "remove its unreferenced module",
                        {
                            operation: plainError(error),
                            cleanup: plainError(cleanupError),
                        },
                    );
                }
                throw error;
            }
        }

        const snapshot = await this.inspect({
            scope,
            name,
            includeSource: false,
        });
        return {
            created: !existing,
            updated: Boolean(existing),
            persisted: true,
            sourceChanged,
            ...(reloadPending ? { activation: "pending" } : {}),
            ...snapshot,
        };
    }

    async run({ scope, name }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(name);
        return this.#mutate(async () => {
            await this.root.hmr.awaitRefresh();
            const sourceFile = this.#sourceFile(scope, name);
            validateModuleSource(
                await readFile(sourceFile, "utf8"),
                sourceFile,
            );
            const state = await this.#configState(scope);
            const index = state.entries.findIndex(
                (entry) => entry.id === name,
            );
            if (index < 0) {
                throw new RuntimeError(
                    "package_not_found",
                    `package ${packageKey(scope, name)} does not exist`,
                );
            }
            const live = this.#entry(scope, name);
            if (!state.entries[index].disabled && live?.fiber?.uid) {
                return {
                    changed: false,
                    ...(await this.inspect({ scope, name })),
                };
            }
            const wasDisabled = Boolean(state.entries[index].disabled);
            if (wasDisabled) {
                await this.root.hmr.refreshFile(sourceFile);
                const entries = structuredClone(state.entries);
                delete entries[index].disabled;
                await this.#commitConfig(scope, state.content, entries);
            } else {
                await live.refresh();
                await this.root.loader.await();
            }
            return {
                changed: true,
                ...(await this.inspect({ scope, name })),
            };
        });
    }

    async reload({ scope, name }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(name);
        return this.#mutate(async () => {
            await this.root.hmr.awaitRefresh();
            const sourceFile = this.#sourceFile(scope, name);
            const source = prepareSource(await readFile(sourceFile, "utf8"));
            const update = await this.#replaceSource(scope, name, source, {
                force: true,
            });
            const snapshot = await this.inspect({ scope, name });
            return {
                changed: update.changed,
                ...(snapshot.enabled ? { activation: "pending" } : {}),
                ...snapshot,
            };
        });
    }

    async stop({ scope, name }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(name);
        return this.#mutate(async () => {
            await this.root.hmr.awaitRefresh();
            const state = await this.#configState(scope);
            const index = state.entries.findIndex(
                (entry) => entry.id === name,
            );
            if (index < 0) {
                throw new RuntimeError(
                    "package_not_found",
                    `package ${packageKey(scope, name)} does not exist`,
                );
            }
            if (state.entries[index].disabled) {
                return {
                    stopped: false,
                    ...(await this.inspect({ scope, name })),
                };
            }
            const entries = structuredClone(state.entries);
            entries[index].disabled = true;
            await this.#commitConfig(scope, state.content, entries);
            return { stopped: true, ...(await this.inspect({ scope, name })) };
        });
    }

    async remove({ scope, name }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(name);
        return this.#mutate(async () => {
            await this.root.hmr.awaitRefresh();
            const state = await this.#configState(scope);
            const entries = state.entries.filter((entry) => entry.id !== name);
            if (entries.length === state.entries.length) {
                await rm(this.#sourceFile(scope, name), { force: true });
                return {
                    scope,
                    name,
                    removed: false,
                    catalogVersion: this.#catalogVersion,
                };
            }
            await this.#commitConfig(scope, state.content, entries);
            try {
                await rm(this.#sourceFile(scope, name), { force: true });
            } catch (error) {
                try {
                    await this.#commitConfig(
                        scope,
                        dumpEntries(entries),
                        state.entries,
                    );
                } catch (rollbackError) {
                    throw new RuntimeError(
                        "remove_rollback_failed",
                        `failed to remove ${packageKey(scope, name)} and ` +
                            "restore its config entry",
                        {
                            operation: plainError(error),
                            rollback: plainError(rollbackError),
                        },
                    );
                }
                throw error;
            }
            return {
                scope,
                name,
                removed: true,
                catalogVersion: this.#catalogVersion,
            };
        });
    }

    async promote({ name, targetName = undefined, activate = false }) {
        this.#assertOpen();
        validateName(name);
        if (targetName !== undefined) validateName(targetName);
        const destination = targetName ?? name;
        return this.#mutate(async () => {
            await this.root.hmr.awaitRefresh();
            const source = await readFile(
                this.#sourceFile("scratch", name),
                "utf8",
            );
            const promoted = await this.#defineUnlocked({
                scope: "project",
                name: destination,
                source,
                activate,
            });
            return {
                ...promoted,
                promotedFrom: { scope: "scratch", name },
            };
        });
    }

    async listPackages({ scope = undefined } = {}) {
        this.#assertOpen();
        if (scope !== undefined) validateScope(scope);
        const scopes = scope ? [scope] : ["project", "scratch"];
        const packages = [];
        const errors = [];
        for (const currentScope of scopes) {
            if (!this.#includes.has(currentScope)) {
                if (scope !== undefined) this.#include(currentScope);
                errors.push({
                    scope: currentScope,
                    error: structuredClone(
                        this.#scopeErrors.get(currentScope) ??
                            plainError(
                                new RuntimeError(
                                    "scope_unavailable",
                                    `${currentScope} cordis.yaml did not load`,
                                ),
                            ),
                    ),
                });
                continue;
            }
            let state;
            try {
                state = await this.#configState(currentScope);
            } catch (error) {
                if (scope !== undefined) throw error;
                errors.push({
                    scope: currentScope,
                    error: plainError(error),
                });
                continue;
            }
            for (const entry of state.entries) {
                if (!NAME_PATTERN.test(entry.id ?? "")) continue;
                packages.push(
                    await this.#snapshot(currentScope, entry.id, entry, false),
                );
            }
        }
        packages.sort((left, right) => {
            return `${left.scope}/${left.name}`.localeCompare(
                `${right.scope}/${right.name}`,
            );
        });
        return { catalogVersion: this.#catalogVersion, packages, errors };
    }

    async inspect({ scope, name, includeSource = false }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(name);
        const state = await this.#configState(scope);
        const entry = state.entries.find((candidate) => candidate.id === name);
        if (!entry) {
            throw new RuntimeError(
                "package_not_found",
                `package ${packageKey(scope, name)} does not exist`,
            );
        }
        return {
            catalogVersion: this.#catalogVersion,
            ...(await this.#snapshot(scope, name, entry, includeSource)),
        };
    }

    async #snapshot(scope, name, options, includeSource) {
        const live = this.#entry(scope, name);
        const sourceFile = this.#sourceFile(scope, name);
        const source = await readMaybe(sourceFile);
        const callback = live?.fiber?.runtime?.callback;
        return {
            scope,
            name,
            description:
                typeof callback?.description === "string"
                    ? callback.description
                    : "",
            module: options.name,
            enabled: !Boolean(options.disabled),
            running: Boolean(live?.fiber?.uid),
            sourceFile: path.relative(this.workspaceRoot, sourceFile),
            sourceSha256: source === undefined ? null : hashSource(source),
            ...(includeSource ? { source: source ?? null } : {}),
        };
    }

    listTools({ scope = undefined, packageName = undefined } = {}) {
        this.#assertOpen();
        if (scope !== undefined) validateScope(scope);
        if (packageName !== undefined) validateName(packageName);
        const tools = [];
        for (const [key, packageTools] of this.#tools) {
            const separator = key.indexOf(":");
            const currentScope = key.slice(0, separator);
            const currentName = key.slice(separator + 1);
            if (scope !== undefined && scope !== currentScope) continue;
            if (packageName !== undefined && packageName !== currentName) {
                continue;
            }
            for (const record of packageTools.values()) {
                tools.push({
                    scope: currentScope,
                    package: currentName,
                    ...record.metadata,
                });
            }
        }
        tools.sort((left, right) => {
            return `${left.scope}/${left.package}/${left.name}`.localeCompare(
                `${right.scope}/${right.package}/${right.name}`,
            );
        });
        return { catalogVersion: this.#catalogVersion, tools };
    }

    async invoke({
        scope,
        packageName,
        tool,
        arguments: args = {},
        timeoutMs = undefined,
        catalogVersion = undefined,
    }) {
        this.#assertOpen();
        validateScope(scope);
        validateName(packageName);
        validateName(tool);
        args = normalizeArguments(args);
        if (
            catalogVersion !== undefined &&
            catalogVersion !== this.#catalogVersion
        ) {
            throw new RuntimeError(
                "stale_catalog",
                `catalog ${catalogVersion} is stale; current catalog is ` +
                    `${this.#catalogVersion}`,
                { catalogVersion: this.#catalogVersion },
            );
        }
        const record = this.#tools
            .get(packageKey(scope, packageName))
            ?.get(tool);
        if (!record) {
            throw new RuntimeError(
                "tool_not_found",
                `tool ${scope}:${packageName}/${tool} is not running`,
            );
        }
        if (!record.validate(args)) {
            throw new RuntimeError(
                "invalid_arguments",
                formatValidationErrors(record.validate.errors),
                { errors: record.validate.errors },
            );
        }
        const selectedTimeout = positiveInteger(
            timeoutMs,
            this.invokeTimeoutMs,
            "timeoutMs",
            300_000,
        );

        const done = deferred();
        let disposeLease;
        try {
            disposeLease = record.context.effect(
                () => () => done.promise,
                `mcp_cordis.invoke(${JSON.stringify(tool)})`,
            );
        } catch (error) {
            throw new RuntimeError(
                "package_not_running",
                `package ${packageKey(scope, packageName)} is disposing`,
                plainError(error),
            );
        }

        let completed = false;
        const invocationSupervisor = new ProcessSupervisor();
        const call = Promise.resolve()
            .then(() =>
                this.#invocationSupervisors.run(invocationSupervisor, () =>
                    record.handler(args),
                ),
            )
            .then((value) => jsonClone(value, `result from ${tool}`))
            .finally(() => {
                completed = true;
            });
        void call.catch(() => {});
        const timeoutError = new RuntimeError(
            "invoke_timeout",
            `tool invocation exceeded ${selectedTimeout} ms`,
        );
        try {
            const value = await withTimeout(
                call,
                selectedTimeout,
                timeoutError,
            );
            return {
                catalogVersion: this.#catalogVersion,
                scope,
                package: packageName,
                tool,
                value,
            };
        } finally {
            await invocationSupervisor.close(
                codedError(
                    "EXEC_DISPOSED",
                    completed
                        ? "ctx.exec invocation scope has completed"
                        : "ctx.exec was cancelled at the invocation deadline",
                ),
            );
            if (completed) {
                done.resolve();
                await disposeLease();
            } else {
                void call
                    .finally(async () => {
                        done.resolve();
                        await disposeLease();
                    })
                    .catch((error) => {
                        process.stderr.write(
                            `[mcp_cordis] invocation lease cleanup failed: ` +
                                `${error instanceof Error ? error.message : error}\n`,
                        );
                    });
            }
        }
    }

    async #replaceSource(scope, name, source, { force = false } = {}) {
        const filename = this.#sourceFile(scope, name);
        const previous = await readMaybe(filename);
        if (!force && previous === source) {
            return { changed: false };
        }
        const entry = this.#entry(scope, name);
        if (entry && entry.options.name !== managedSpecifier(name)) {
            throw new RuntimeError(
                "unmanaged_entry",
                `${packageKey(scope, name)} does not use its managed module path`,
            );
        }
        await atomicWrite(filename, source);
        return { changed: true };
    }

    async #configState(scope) {
        validateScope(scope);
        const content = await readFile(this.configFiles[scope], "utf8");
        return {
            content,
            entries: parseEntries(content, this.configFiles[scope]),
        };
    }

    async #validateConfigSources(scope) {
        const state = await this.#configState(scope);
        for (const entry of state.entries) {
            if (
                !NAME_PATTERN.test(entry.id ?? "") ||
                entry.name !== managedSpecifier(entry.id)
            ) {
                throw new RuntimeError(
                    "unsupported_cordis_entry",
                    `${scope} cordis.yaml contains an unsupported entry`,
                );
            }
            const sourceFile = this.#sourceFile(scope, entry.id);
            validateModuleSource(
                await readFile(sourceFile, "utf8"),
                sourceFile,
            );
        }
    }

    async #commitConfig(scope, previousContent, entries) {
        const include = this.#include(scope);
        const content = dumpEntries(entries);
        if (content === previousContent) return;
        await atomicWrite(this.configFiles[scope], content);
        try {
            await include.refresh();
            await this.root.loader.await();
        } catch (error) {
            let rollbackError;
            try {
                await atomicWrite(this.configFiles[scope], previousContent);
                await include.refresh();
                await this.root.loader.await();
            } catch (failure) {
                rollbackError = failure;
            }
            throw new RuntimeError(
                rollbackError ? "config_rollback_failed" : "activation_failed",
                `failed to apply ${scope} cordis.yaml` +
                    (rollbackError ? " and restore its prior entries" : ""),
                {
                    activation: plainError(error),
                    ...(rollbackError
                        ? { rollback: plainError(rollbackError) }
                        : {}),
                },
            );
        }
    }

    #include(scope) {
        const include = this.#includes.get(validateScope(scope));
        if (!include) {
            throw new RuntimeError(
                "scope_unavailable",
                `${scope} cordis.yaml did not load`,
            );
        }
        return include;
    }

    #entry(scope, name) {
        return this.#include(scope).store[validateName(name)];
    }

    #sourceFile(scope, name) {
        return path.join(
            this.pluginRoots[validateScope(scope)],
            `${validateName(name)}.mjs`,
        );
    }

    #loadedPackages() {
        const loaded = [];
        for (const [scope, include] of this.#includes) {
            for (const entry of Object.values(include.store)) {
                if (!entry.fiber?.uid) continue;
                loaded.push({
                    scope,
                    name: entry.options.id,
                    running: true,
                });
            }
        }
        loaded.sort((left, right) => {
            return `${left.scope}/${left.name}`.localeCompare(
                `${right.scope}/${right.name}`,
            );
        });
        return loaded;
    }

    #mutate(callback) {
        this.#assertOpen();
        const preceding = this.#mutation.catch(() => {});
        const current = preceding.then(() => {
            return callback();
        });
        this.#mutation = current;
        return this.#admit(() => current);
    }

    #admit(callback) {
        const admitted = Promise.resolve().then(callback);
        this.#admissions.add(admitted);
        return admitted.finally(() => this.#admissions.delete(admitted));
    }

    shutdown() {
        if (this.#shutdownPromise) return this.#shutdownPromise;
        this.#closed = true;
        const admitted = [...this.#admissions];
        this.#shutdownPromise = (async () => {
            await Promise.allSettled(admitted);
            // The root Fiber has uid 0, so truthiness would skip all Cordis
            // effects, including HMR watcher disposal.
            if (this.root?.fiber) {
                await this.root.fiber.dispose();
            }
        })();
        return this.#shutdownPromise;
    }

    #assertOpen() {
        if (this.#closed) {
            throw new RuntimeError("runtime_closed", "runtime is closed");
        }
    }
}

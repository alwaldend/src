/** URLs resolve relative to index.js. The callable Rust WASI backend loads lazily. */
export interface BrowserOptions {
    wasmURL?: string | URL;
    workerURL?: string | URL;
    serverWasmURL?: string | URL;
}

export interface ConfigurationApp {
    /**
     * Write a CollectionConfigSet protobuf JSON object as a TOML tree.
     * Use string values for int64 fields. Declared files are overwritten;
     * unrelated files remain and filesystem errors can leave partial output.
     * Loading checks supported versions and layout; resource relationships
     * remain outside configuration serialization.
     */
    configExport(options: {
        output: string;
        configs: Record<string, unknown>;
    }): Promise<void>;

    /** Load a config directory or marker and return encoded text. */
    configDump(options: {
        text: string;
        format?: "json" | "toml";
    }): Promise<string>;
}

/** Managed collection operations accept completed .colpkg packages. */
export interface CollectionOptions {
    collection: string;
    stateKey?: string;
}

/** Existing App workflows invoke Rust's JavaScript command API over the shared workspace. */
export interface CollectionApp extends ConfigurationApp {
    /** Export editable TOML and optional bindings without changing collection resources. */
    collectionExport(
        options: CollectionOptions & {
            output: string;
            stateOutput?: string;
        },
    ): Promise<void>;

    /** Merge bindings from an exported TOML state file into the collection. */
    collectionStateUpdate(
        options: CollectionOptions & {
            state: string;
        },
    ): Promise<void>;

    /** Return a CollectionChangePlan protobuf JSON object. */
    collectionPlan(
        options: CollectionOptions & {
            text: string;
        },
    ): Promise<Record<string, unknown>>;

    /** Reconcile a package in place, including managed publication and cleanup. */
    collectionApply(
        options: CollectionOptions & {
            text: string;
        },
    ): Promise<void>;

    /** Build an independent package; failed builds preserve an existing destination. */
    collectionBuild(options: {
        base: string;
        text: string;
        output: string;
        stateKey?: string;
    }): Promise<void>;
}

/** Paths are absolute virtual paths; /tmp is reserved for the collection runtime. */
export interface Workspace {
    /**
     * Copy file bytes, creating parents and replacing regular-file contents.
     * Directory targets and regular-file ancestors are rejected without changes.
     */
    writeFile(path: string, data: Uint8Array): Promise<void>;

    /** Return an independent copy of a file's bytes. */
    readFile(path: string): Promise<Uint8Array>;

    /** Return sorted absolute file paths beneath a path, recursively. */
    listFiles(path: string): Promise<string[]>;

    /** Remove a path and its descendants; absent paths are harmless. */
    remove(path: string): Promise<void>;

    /**
     * Atomically rename a regular file, replacing an existing regular file.
     * The destination parent must exist; directories cannot be renamed or replaced.
     */
    rename(source: string, destination: string): Promise<void>;
}

export interface SavedWorkspace {
    revision: number;
    /** ISO 8601 timestamp of the completed snapshot transaction. */
    savedAt: string;
}

export interface WorkspaceSession {
    /**
     * Atomically save all files and directories except /tmp after scoped backends
     * finish cleanup. Pass the revision returned by the last restore/save.
     * Zero (the default) creates a new key at revision 1; stale writers reject.
     * The caller owns its expected revision across workspace replacements.
     * Failed transactions preserve the previous snapshot and in-memory workspace.
     */
    save(key: string, expectedRevision?: number): Promise<SavedWorkspace>;

    /**
     * Validate and replace the workspace with the saved snapshot. Missing keys
     * return null without replacing files. Failed restores preserve current files.
     * Reads do not require configuration TOML to be valid.
     */
    restore(key: string): Promise<SavedWorkspace | null>;
}

export interface BrowserApplication {
    app: CollectionApp;
    workspace: Workspace;
    session: WorkspaceSession;

    /**
     * Terminate the worker and discard its memory workspace immediately.
     * Pending and future operations reject; repeated disposal is safe. No implicit
     * save occurs. Await an explicit save before disposal when persistence is needed.
     */
    dispose(): Promise<void>;
}

/**
 * Start an isolated worker and resolve when its Go application is ready.
 * Startup rejects after at most 30 seconds. Application, workspace and session
 * methods execute in request order. Collection calls asynchronously load the
 * callable Rust WASI backend in server.wasm; configuration-only usage skips it.
 * All workflows run without JSPI. Rust commands execute synchronously in the
 * worker and cannot be interrupted by cancellation before they return.
 */
export function createAnkiApp(
    options?: BrowserOptions,
): Promise<BrowserApplication>;

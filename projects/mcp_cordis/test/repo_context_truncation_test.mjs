import assert from "node:assert/strict";
import {
    mkdir,
    mkdtemp,
    readFile,
    symlink,
    writeFile,
} from "node:fs/promises";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import test from "node:test";
import plugin from "../plugins/repo_context.mjs";

const workspaceRoot = fileURLToPath(new URL("../plugins/", import.meta.url));
const temporaryRoot = process.env.TEST_TMPDIR;
const objectHash = "0123456789abcdef".repeat(2) + "01234567";
const normalStatus = "## main\n M fixture.txt\n";

function execution(stdout, overrides = {}) {
    return {
        code: 0,
        outputLimitExceeded: false,
        signal: null,
        stderr: "",
        stdout,
        truncated: false,
        ...overrides,
    };
}

function outputLimited(stdout = "", stderr = "") {
    return execution(stdout, {
        code: null,
        outputLimitExceeded: true,
        signal: "SIGTERM",
        stderr,
        truncated: true,
    });
}

function loadHandlers(overrides = {}, selectedRoot = workspaceRoot) {
    const handlers = new Map();
    const calls = [];
    const results = {
        root: execution(`${selectedRoot}\n`),
        head: execution(`${objectHash}\n`),
        status: execution(normalStatus),
        rg: execution(""),
        ...overrides,
    };
    plugin.apply({
        exec: async (file, args, options) => {
            calls.push({ file, args: [...args], options: { ...options } });
            if (file === "rg") {
                if (results.rg instanceof Error) throw results.rg;
                return results.rg;
            }
            assert.equal(file, "git");
            if (args.includes("--show-toplevel")) return results.root;
            if (args.includes("--verify")) return results.head;
            if (args.includes("status")) return results.status;
            throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
        },
        resolveWorkspace: (relativePath) => {
            return resolve(selectedRoot, relativePath);
        },
        readText:
            overrides.readText ??
            ((relativePath) => {
                return readFile(resolve(selectedRoot, relativePath), "utf8");
            }),
        tool: (definition, handler) => {
            handlers.set(definition.name, handler);
        },
    });
    return { calls, handlers };
}

async function getContext(
    overrides = {},
    maxBytes = 1024,
    selectedRoot = workspaceRoot,
) {
    const loaded = loadHandlers(overrides, selectedRoot);
    const handler = loaded.handlers.get("repo_context_get");
    const value = await handler({
        path: ".",
        include_instructions: false,
        max_bytes: maxBytes,
    });
    return { ...loaded, value };
}

function rgMatch(line, text = "needle") {
    return JSON.stringify({
        type: "match",
        data: {
            path: { text: "fixture.txt" },
            line_number: line,
            lines: { text: `${text}\n` },
            submatches: [
                {
                    start: 0,
                    end: text.length,
                    match: { text },
                },
            ],
        },
    });
}

function rgBytesMatch(field) {
    const event = JSON.parse(rgMatch(1));
    const encoded = Buffer.from("non-utf8\xff", "latin1").toString("base64");
    if (field === "path") event.data.path = { bytes: encoded };
    if (field === "lines") event.data.lines = { bytes: encoded };
    if (field === "submatch") {
        event.data.submatches[0].match = { bytes: encoded };
    }
    return JSON.stringify(event);
}

test("repo_context_get rejects truncated Git root discovery", async () => {
    const root = outputLimited(`${workspaceRoot}\n`);
    const { calls, value } = await getContext({ root });

    assert.equal(value.git.rootTruncated, true);
    assert.match(value.git.error, /root discovery output was truncated/i);
    assert.equal("root" in value.git, false);
    assert.equal(calls.length, 1);
});

test("repo_context_get never reports a truncated HEAD as complete", async () => {
    const head = outputLimited(objectHash.slice(0, 20));
    const { value } = await getContext({ head });

    assert.equal(value.git.root, ".");
    assert.equal(value.git.head, null);
    assert.equal(value.git.headTruncated, true);
    assert.equal(value.git.statusTruncated, false);
});

test("repo_context_get propagates host status truncation", async () => {
    const status = outputLimited(normalStatus);
    const { value } = await getContext({ status });

    assert.equal(value.git.status, normalStatus);
    assert.equal(value.git.statusTruncated, true);
});

test("repo_context_get reports local status clipping", async () => {
    const status = execution("x".repeat(1500));
    const { value } = await getContext({ status });

    assert.equal(Buffer.byteLength(value.git.status, "utf8"), 1024);
    assert.equal(value.git.statusTruncated, true);
});

test("repo_context_get marks complete Git results explicitly", async () => {
    const { value } = await getContext();

    assert.equal(value.git.root, ".");
    assert.equal(value.git.rootTruncated, false);
    assert.equal(value.git.head, objectHash);
    assert.equal(value.git.headTruncated, false);
    assert.equal(value.git.status, normalStatus);
    assert.equal(value.git.statusTruncated, false);
});

test("repo_context_search preserves ripgrep truncation", async () => {
    const rg = outputLimited("");
    const { handlers } = loadHandlers({ rg });
    const value = await handlers.get("repo_context_search")({
        query: "needle",
    });

    assert.equal(value.engine, "ripgrep");
    assert.equal(value.truncated, true);
});

test("repo_context_get distinguishes nonzero Git results", async (t) => {
    await t.test("root means no repository", async () => {
        const root = execution("", { code: 128, stderr: "not a repo" });
        const { value } = await getContext({ root });
        assert.equal(value.git, null);
    });

    await t.test("HEAD is explicitly unavailable", async () => {
        const head = execution("", { code: 128, stderr: "unborn" });
        const { value } = await getContext({ head });
        assert.equal(value.git.head, null);
        assert.equal(value.git.headAvailable, false);
        assert.equal(value.git.headTruncated, false);
        assert.equal(value.git.headError, "unborn");
    });

    await t.test("status is explicitly unavailable", async () => {
        const status = execution("", { code: 2, stderr: "status failed" });
        const { value } = await getContext({ status });
        assert.equal(value.git.status, null);
        assert.equal(value.git.statusAvailable, false);
        assert.equal(value.git.statusTruncated, false);
        assert.equal(value.git.statusError, "status failed");
    });
});

test("repo_context rejects unexpected process signals", async (t) => {
    const signaled = execution("partial", {
        code: null,
        signal: "SIGTERM",
    });
    for (const command of ["root", "head", "status", "rg"]) {
        await t.test(command, async () => {
            if (command === "rg") {
                const { handlers } = loadHandlers({ rg: signaled });
                await assert.rejects(
                    handlers.get("repo_context_search")({ query: "needle" }),
                    /terminated by SIGTERM/,
                );
                return;
            }
            await assert.rejects(
                getContext({ [command]: signaled }),
                /terminated by SIGTERM/,
            );
        });
    }
});

test("repo_context rejects unmarked incomplete output", async () => {
    const root = execution(`${workspaceRoot}\n`, { truncated: true });
    await assert.rejects(
        getContext({ root }),
        /without an output-limit marker/,
    );
});

test("repo_context_search distinguishes ripgrep exits", async (t) => {
    await t.test("exit one is a complete empty result", async () => {
        const { handlers } = loadHandlers({ rg: execution("", { code: 1 }) });
        const value = await handlers.get("repo_context_search")({
            query: "needle",
        });
        assert.equal(value.matchCount, 0);
        assert.equal(value.truncated, false);
    });

    await t.test("other nonzero exits reject", async () => {
        const { handlers } = loadHandlers({
            rg: execution("", { code: 2, stderr: "bad regex" }),
        });
        await assert.rejects(
            handlers.get("repo_context_search")({ query: "needle" }),
            /ripgrep exited with 2: bad regex/,
        );
    });
});

test("ripgrep max_matches needs an observed extra match", async (t) => {
    async function search(lines) {
        const { handlers } = loadHandlers({
            rg: execution(`${lines.join("\n")}\n`),
        });
        return handlers.get("repo_context_search")({
            query: "needle",
            max_matches: 1,
        });
    }

    await t.test("exactly one is complete", async () => {
        const value = await search([rgMatch(1)]);
        assert.equal(value.matchCount, 1);
        assert.equal(value.matches.length, 1);
        assert.equal(value.truncated, false);
    });

    await t.test("a second observed match truncates", async () => {
        const value = await search([rgMatch(1), rgMatch(2)]);
        assert.equal(value.matchCount, 1);
        assert.equal(value.matches.length, 1);
        assert.equal(value.truncated, true);
    });
});

test("ripgrep byte fields never masquerade as complete text", async (t) => {
    for (const field of ["path", "lines", "submatch"]) {
        await t.test(field, async () => {
            const { handlers } = loadHandlers({
                rg: execution(`${rgBytesMatch(field)}\n`),
            });
            const value = await handlers.get("repo_context_search")({
                query: "needle",
            });
            assert.equal(value.matchCount, 1);
            assert.equal(value.truncated, true);
            if (field === "path") {
                assert.equal(value.matches[0].path, "");
                assert.equal(value.matches[0].pathTruncated, true);
            } else if (field === "lines") {
                assert.equal(value.matches[0].text, "");
                assert.equal(value.matches[0].textTruncated, true);
            } else {
                assert.equal(value.matches[0].submatches[0].text, "");
                assert.equal(
                    value.matches[0].submatches[0].textTruncated,
                    true,
                );
            }
        });
    }
});

test("ripgrep drops an incomplete output-limited JSON event", async () => {
    const rg = outputLimited(`${rgMatch(1)}\n{"type":"match"`);
    const { handlers } = loadHandlers({ rg });
    const value = await handlers.get("repo_context_search")({
        query: "needle",
        max_matches: 5,
    });

    assert.equal(value.matchCount, 1);
    assert.equal(value.matches.length, 1);
    assert.equal(value.truncated, true);
});

test("JavaScript fallback uses exact max_matches semantics", async (t) => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, "fixture.txt"), "needle\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });

    async function search(content) {
        await writeFile(join(root, "fixture.txt"), content);
        const loaded = loadHandlers({ rg: missing }, root);
        return loaded.handlers.get("repo_context_search")({
            query: "needle",
            paths: ["fixture.txt"],
            max_matches: 1,
        });
    }

    await t.test("exactly one is complete", async () => {
        const value = await search("needle\n");
        assert.equal(value.engine, "javascript");
        assert.equal(value.matchCount, 1);
        assert.equal(value.truncated, false);
    });

    await t.test("an observed second match truncates", async () => {
        const value = await search("needle\nneedle\n");
        assert.equal(value.matchCount, 1);
        assert.equal(value.truncated, true);
    });
});

test("JavaScript fallback reports locally clipped details", async (t) => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });

    async function search(query, content) {
        await writeFile(join(root, "fixture.txt"), content);
        const loaded = loadHandlers({ rg: missing }, root);
        return loaded.handlers.get("repo_context_search")({
            query,
            paths: ["fixture.txt"],
            max_matches: 2,
        });
    }

    await t.test("a long submatch", async () => {
        const query = "n".repeat(600);
        const value = await search(query, `${query}\n`);
        assert.equal(value.matches[0].submatches[0].text.length, 512);
        assert.equal(value.truncated, true);
    });

    await t.test("a 33rd submatch", async () => {
        const value = await search("n", `${"n".repeat(33)}\n`);
        assert.equal(value.matches[0].submatches.length, 32);
        assert.equal(value.truncated, true);
    });

    await t.test("a long matching line", async () => {
        const value = await search("needle", `needle${"x".repeat(5000)}\n`);
        assert.equal(Buffer.byteLength(value.matches[0].text), 4096);
        assert.equal(value.matches[0].textTruncated, true);
        assert.equal(value.truncated, true);
    });
});

test("JavaScript fallback reports UTF-8 byte offsets", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, "fixture.txt"), "éneedle\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);

    const value = await handlers.get("repo_context_search")({
        query: "needle",
        paths: ["fixture.txt"],
    });
    assert.deepEqual(value.matches[0].submatches[0], {
        start: 2,
        end: 8,
        text: "needle",
    });
});

test("case-insensitive fallback keeps original-line offsets", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, "fixture.txt"), "İx\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);
    const value = await handlers.get("repo_context_search")({
        query: "x",
        paths: ["fixture.txt"],
        case_sensitive: false,
    });

    assert.deepEqual(value.matches[0].submatches[0], {
        start: 2,
        end: 3,
        text: "x",
    });
});

test("regular-expression fallback requires ripgrep", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, "fixture.txt"), "😀\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);
    await assert.rejects(
        handlers.get("repo_context_search")({
            query: "(a+)+$",
            paths: ["fixture.txt"],
            fixed: false,
        }),
        /ripgrep is required for regular-expression searches/,
    );
});

test("invalid UTF-8 fallback requires ripgrep", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(
        join(root, "fixture.txt"),
        Buffer.concat([Buffer.from([0xff]), Buffer.from("needle")]),
    );
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);
    await assert.rejects(
        handlers.get("repo_context_search")({
            query: "needle",
            paths: ["fixture.txt"],
        }),
        /ripgrep is required to search files containing invalid UTF-8/,
    );
});

test("fallback context excludes matching lines and stays ordered", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, "fixture.txt"), "needle\nneedle\nafter\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);
    const value = await handlers.get("repo_context_search")({
        query: "needle",
        paths: ["fixture.txt"],
        context_lines: 1,
    });

    assert.deepEqual(
        value.matches.map(({ kind, line }) => ({ kind, line })),
        [
            { kind: "match", line: 1 },
            { kind: "match", line: 2 },
            { kind: "context", line: 3 },
        ],
    );
});

test("fallback context does not invent a line after final newline", async (t) => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });

    await t.test("one newline ends the matching line", async () => {
        const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
        await writeFile(join(root, "fixture.txt"), "needle\n");
        const { handlers } = loadHandlers({ rg: missing }, root);
        const value = await handlers.get("repo_context_search")({
            query: "needle",
            paths: ["fixture.txt"],
            context_lines: 1,
        });

        assert.deepEqual(
            value.matches.map(({ kind, line }) => ({ kind, line })),
            [{ kind: "match", line: 1 }],
        );
    });

    await t.test("a second newline is a real empty line", async () => {
        const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
        await writeFile(join(root, "fixture.txt"), "needle\n\n");
        const { handlers } = loadHandlers({ rg: missing }, root);
        const value = await handlers.get("repo_context_search")({
            query: "needle",
            paths: ["fixture.txt"],
            context_lines: 1,
        });

        assert.deepEqual(
            value.matches.map(({ kind, line, text }) => ({
                kind,
                line,
                text,
            })),
            [
                { kind: "match", line: 1, text: "needle" },
                { kind: "context", line: 2, text: "" },
            ],
        );
    });
});

test("explicit fallback files override globs", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, "fixture.txt"), "needle\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);
    const value = await handlers.get("repo_context_search")({
        query: "needle",
        paths: ["fixture.txt"],
        globs: ["*.js"],
    });

    assert.equal(value.matchCount, 1);
    assert.equal(value.matches[0].path, "fixture.txt");
});

test("JavaScript fallback fails closed for directory recursion", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-fallback-"));
    await writeFile(join(root, ".env"), "needle=secret\n");
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);

    await assert.rejects(
        handlers.get("repo_context_search")({ query: "needle" }),
        /ripgrep is required to search directories/,
    );
    const explicit = await handlers.get("repo_context_search")({
        query: "needle",
        paths: [".env"],
    });
    assert.equal(explicit.matchCount, 1);
    assert.equal(explicit.matches[0].path, ".env");
});

test("instruction truncation reports exact-byte omissions", async (t) => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");

    await t.test(
        "exact content without a later document is complete",
        async () => {
            const root = await mkdtemp(
                join(temporaryRoot, "repo-instructions-"),
            );
            await writeFile(join(root, "AGENTS.md"), "a".repeat(1024));
            const { handlers } = loadHandlers({}, root);
            const value = await handlers.get("repo_context_get")({
                path: ".",
                include_git: false,
                max_bytes: 1024,
            });
            assert.equal(value.instructions.length, 1);
            assert.equal(value.instructions[0].truncated, false);
            assert.equal(value.instructionsTruncated, false);
        },
    );

    await t.test("a later applicable document is reported", async () => {
        const root = await mkdtemp(join(temporaryRoot, "repo-instructions-"));
        await mkdir(join(root, "nested"));
        await writeFile(join(root, "AGENTS.md"), "a".repeat(1024));
        await writeFile(join(root, "nested", "AGENTS.md"), "later\n");
        const { handlers } = loadHandlers({}, root);
        const value = await handlers.get("repo_context_get")({
            path: "nested",
            include_git: false,
            max_bytes: 1024,
        });
        assert.equal(value.instructions.length, 1);
        assert.equal(value.instructions[0].path, "AGENTS.md");
        assert.equal(value.instructionsTruncated, true);
    });
});

test("directory entry truncation observes the 129th entry", async (t) => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");

    async function contextWithEntries(count) {
        const root = await mkdtemp(join(temporaryRoot, "repo-entries-"));
        await Promise.all(
            Array.from({ length: count }, (_, index) => {
                const name = `entry-${String(index).padStart(3, "0")}`;
                return writeFile(join(root, name), "");
            }),
        );
        const { handlers } = loadHandlers({}, root);
        return handlers.get("repo_context_get")({
            path: ".",
            include_git: false,
            include_instructions: false,
            max_bytes: 1024,
        });
    }

    await t.test("128 entries are complete", async () => {
        const value = await contextWithEntries(128);
        assert.equal(value.entries.length, 128);
        assert.equal(value.entriesTruncated, false);
    });

    await t.test("129 entries are truncated", async () => {
        const value = await contextWithEntries(129);
        assert.equal(value.entries.length, 128);
        assert.equal(value.entriesTruncated, true);
    });
});

test("repo_context accepts legal '..config' path components", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-path-"));
    await mkdir(join(root, "..config"));
    const { handlers } = loadHandlers({}, root);
    const value = await handlers.get("repo_context_get")({
        path: "..config",
        include_git: false,
        include_instructions: false,
        max_bytes: 1024,
    });
    assert.equal(value.path, "..config");
});

test("repo_context rejects symlink escapes", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-symlink-root-"));
    const outside = await mkdtemp(
        join(temporaryRoot, "repo-symlink-outside-"),
    );
    await writeFile(join(outside, "outside.txt"), "outside\n");
    await symlink(join(outside, "outside.txt"), join(root, "escape.txt"));
    const { handlers } = loadHandlers({}, root);

    await assert.rejects(
        handlers.get("repo_context_get")({
            path: "escape.txt",
            include_git: false,
            include_instructions: false,
            max_bytes: 1024,
        }),
        /outside the workspace/,
    );
    await assert.rejects(
        handlers.get("repo_context_read")({
            files: [{ path: "escape.txt" }],
            max_total_bytes: 1024,
        }),
        /outside the workspace/,
    );
    await assert.rejects(
        handlers.get("repo_context_search")({
            query: "outside",
            paths: ["escape.txt"],
        }),
        /outside the workspace/,
    );
});

test("fallback follows an explicitly selected internal symlink", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-symlink-root-"));
    await writeFile(join(root, "target.txt"), "needle\n");
    await symlink("target.txt", join(root, "alias.txt"));
    const missing = Object.assign(new Error("missing rg"), { code: "ENOENT" });
    const { handlers } = loadHandlers({ rg: missing }, root);
    const value = await handlers.get("repo_context_search")({
        query: "needle",
        paths: ["alias.txt"],
        max_matches: 2,
    });

    assert.equal(value.engine, "javascript");
    assert.equal(value.matchCount, 1);
    assert.equal(value.matches[0].path, "alias.txt");
    assert.equal(value.truncated, false);
});

test("repo_context reads the verified symlink target", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-symlink-root-"));
    await writeFile(join(root, "target.txt"), "inside\n");
    await symlink("target.txt", join(root, "alias.txt"));
    const { handlers } = loadHandlers(
        {
            readText: async () => {
                throw new Error("the lexical alias was reopened");
            },
        },
        root,
    );

    const value = await handlers.get("repo_context_read")({
        files: [{ path: "alias.txt" }],
        max_total_bytes: 1024,
    });

    assert.equal(value.files[0].content, "inside");
});

test("repo_context opts into subprocess partial output", async () => {
    const { calls, handlers } = loadHandlers({
        rg: execution(`${rgMatch(1)}\n`),
    });
    await handlers.get("repo_context_get")({
        path: ".",
        include_instructions: false,
        max_bytes: 1024,
    });
    await handlers.get("repo_context_search")({ query: "needle" });

    for (const call of calls) {
        assert.equal(call.options.allowTruncatedOutput, true);
        if (call.file === "git") {
            assert.match(
                call.args[1],
                new RegExp(`^/proc/${process.pid}/fd/\\d+$`),
            );
            assert.equal(call.options.env.GIT_OPTIONAL_LOCKS, "0");
            assert.equal(call.options.env.GIT_NO_LAZY_FETCH, "1");
            for (const name of [
                "GIT_ALTERNATE_OBJECT_DIRECTORIES",
                "GIT_COMMON_DIR",
                "GIT_CONFIG",
                "GIT_CONFIG_COUNT",
                "GIT_CONFIG_PARAMETERS",
                "GIT_DIR",
                "GIT_GRAFT_FILE",
                "GIT_IMPLICIT_WORK_TREE",
                "GIT_INDEX_FILE",
                "GIT_INTERNAL_SUPER_PREFIX",
                "GIT_NAMESPACE",
                "GIT_NO_REPLACE_OBJECTS",
                "GIT_OBJECT_DIRECTORY",
                "GIT_PREFIX",
                "GIT_REPLACE_REF_BASE",
                "GIT_SHALLOW_FILE",
                "GIT_WORK_TREE",
            ]) {
                assert.equal(call.options.env[name], null);
            }
        } else {
            assert.equal(call.file, "rg");
            assert.match(
                call.args.at(-1),
                new RegExp(`^/proc/${process.pid}/fd/\\d+$`),
            );
        }
    }
});

test("bounded reads report only retained lines", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-read-"));
    await writeFile(join(root, "fixture.txt"), "abcdef\nsecond\nthird\n");
    const { handlers } = loadHandlers({}, root);
    const value = await handlers.get("repo_context_read")({
        files: [{ path: "fixture.txt", start_line: 1, end_line: 3 }],
        max_total_bytes: 3,
    });

    assert.equal(value.files[0].content, "abc");
    assert.equal(value.files[0].endLine, 1);
    assert.equal(value.files[0].truncated, true);

    const lineBoundary = await handlers.get("repo_context_read")({
        files: [{ path: "fixture.txt", start_line: 1, end_line: 3 }],
        max_total_bytes: 7,
    });
    assert.equal(lineBoundary.files[0].content, "abcdef\n");
    assert.equal(lineBoundary.files[0].endLine, 1);
    assert.equal(lineBoundary.files[0].truncated, true);
});

test("bounded reads distinguish empty retained lines from EOF", async () => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-read-"));
    await writeFile(join(root, "fixture.txt"), "first\n\nthird");
    const { handlers } = loadHandlers({}, root);
    const emptyLine = await handlers.get("repo_context_read")({
        files: [{ path: "fixture.txt", start_line: 2, end_line: 2 }],
    });
    const pastEnd = await handlers.get("repo_context_read")({
        files: [{ path: "fixture.txt", start_line: 4, end_line: 4 }],
    });

    assert.equal(emptyLine.files[0].content, "");
    assert.equal(emptyLine.files[0].endLine, 2);
    assert.equal(emptyLine.files[0].truncated, false);
    assert.equal(pastEnd.files[0].endLine, null);
});

test("bounded reads do not count terminal split sentinels", async (t) => {
    assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");
    const root = await mkdtemp(join(temporaryRoot, "repo-read-"));
    const { handlers } = loadHandlers({}, root);

    await t.test("final newline does not create another line", async () => {
        await writeFile(join(root, "fixture.txt"), "first\n");
        const value = await handlers.get("repo_context_read")({
            files: [{ path: "fixture.txt", start_line: 2 }],
        });

        assert.equal(value.files[0].endLine, null);
        assert.equal(value.files[0].totalLines, 1);
    });

    await t.test("an empty file has no lines", async () => {
        await writeFile(join(root, "fixture.txt"), "");
        const value = await handlers.get("repo_context_read")({
            files: [{ path: "fixture.txt", start_line: 1 }],
        });

        assert.equal(value.files[0].endLine, null);
        assert.equal(value.files[0].totalLines, 0);
    });
});

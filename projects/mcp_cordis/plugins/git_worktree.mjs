import { constants as fsConstants } from "node:fs";
import * as fs from "node:fs/promises";
import * as path from "node:path";
import process from "node:process";

export const description =
    "Read-only Git worktree snapshots and commit comparisons.";

const plugin = {
    name: "git_worktree",
    description,

    apply(ctx) {
        const workspaceRoot = path.resolve(ctx.resolveWorkspace("."));
        const canonicalWorkspaceRoot = fs.realpath(workspaceRoot);
        const MAX_DIFF_BYTES = 2 * 1024 * 1024;

        function integer(value, fallback, minimum, maximum, label) {
            const selected = value === undefined ? fallback : value;
            if (!Number.isInteger(selected)) {
                throw new TypeError(`${label} must be an integer`);
            }
            if (selected < minimum || selected > maximum) {
                throw new RangeError(
                    `${label} must be between ${minimum} and ${maximum}`,
                );
            }
            return selected;
        }

        function string(value, label, maximum, allowEmpty = false) {
            if (typeof value !== "string") {
                throw new TypeError(`${label} must be a string`);
            }
            if (
                value.includes("\0") ||
                value.includes("\n") ||
                value.includes("\r")
            ) {
                throw new TypeError(
                    `${label} contains an invalid control character`,
                );
            }
            if (
                (!allowEmpty && value.length === 0) ||
                value.length > maximum
            ) {
                throw new RangeError(`${label} has an invalid length`);
            }
            return value;
        }

        function clipUtf8(value, maximum) {
            const source = Buffer.from(String(value), "utf8");
            if (source.length <= maximum) {
                return {
                    text: source.toString("utf8"),
                    bytes: source.length,
                    truncated: false,
                };
            }
            let end = maximum;
            while (end > 0 && (source[end] & 0xc0) === 0x80) {
                end -= 1;
            }
            return {
                text: source.subarray(0, end).toString("utf8"),
                bytes: end,
                truncated: true,
            };
        }

        function resultText(result, field) {
            const value = result[field];
            if (Buffer.isBuffer(value)) {
                return value.toString("utf8");
            }
            return value === undefined ? "" : String(value);
        }

        function executionState(result, label) {
            if (!result || typeof result !== "object") {
                throw new Error(
                    `${label} returned an invalid execution result`,
                );
            }
            const outputLimitExceeded = result.outputLimitExceeded === true;
            if (result.truncated && !outputLimitExceeded) {
                throw new Error(
                    `${label} returned incomplete output without an output-limit marker`,
                );
            }
            if (outputLimitExceeded) {
                return { code: null, outputLimitExceeded: true };
            }
            if (result.signal !== null && result.signal !== undefined) {
                throw new Error(`${label} was terminated by ${result.signal}`);
            }
            const code = result.code ?? result.exitCode;
            if (!Number.isInteger(code)) {
                throw new Error(`${label} ended without an exit code`);
            }
            return { code, outputLimitExceeded: false };
        }

        function separatedFields(value, separator) {
            if (value === "") {
                return { fields: [], malformed: false };
            }
            const terminated = value.endsWith(separator);
            const fields = value.split(separator);
            fields.pop();
            return {
                fields,
                malformed: !terminated,
            };
        }

        function escapes(base, candidate) {
            const relative = path.relative(base, candidate);
            return (
                relative === ".." ||
                relative.startsWith(`..${path.sep}`) ||
                path.isAbsolute(relative)
            );
        }

        async function execute(repo, args, options = {}) {
            return ctx.exec("git", ["-C", repo, ...args], {
                cwd: workspaceRoot,
                env: {
                    GIT_ALTERNATE_OBJECT_DIRECTORIES: null,
                    GIT_COMMON_DIR: null,
                    GIT_CONFIG: null,
                    GIT_CONFIG_COUNT: null,
                    GIT_CONFIG_PARAMETERS: null,
                    GIT_DIR: null,
                    GIT_GRAFT_FILE: null,
                    GIT_IMPLICIT_WORK_TREE: null,
                    GIT_INDEX_FILE: null,
                    GIT_INTERNAL_SUPER_PREFIX: null,
                    GIT_NAMESPACE: null,
                    GIT_NO_LAZY_FETCH: "1",
                    GIT_NO_REPLACE_OBJECTS: null,
                    GIT_OBJECT_DIRECTORY: null,
                    GIT_OPTIONAL_LOCKS: "0",
                    GIT_PREFIX: null,
                    GIT_REPLACE_REF_BASE: null,
                    GIT_SHALLOW_FILE: null,
                    GIT_TERMINAL_PROMPT: "0",
                    GIT_WORK_TREE: null,
                },
                allowTruncatedOutput: true,
                timeoutMs: options.timeoutMs ?? 15_000,
                maxBytes: options.maxBytes ?? 512 * 1024,
            });
        }

        async function bindDirectory(absolute, label) {
            const [workspace, canonical] = await Promise.all([
                canonicalWorkspaceRoot,
                fs.realpath(absolute),
            ]);
            if (escapes(workspace, canonical)) {
                throw new RangeError(`${label} is outside the workspace`);
            }
            const handle = await fs.open(
                canonical,
                fsConstants.O_RDONLY |
                    fsConstants.O_DIRECTORY |
                    fsConstants.O_NOFOLLOW,
            );
            try {
                const opened = await fs.realpath(`/proc/self/fd/${handle.fd}`);
                if (escapes(workspace, opened)) {
                    throw new RangeError(`${label} is outside the workspace`);
                }
                return {
                    absolute: opened,
                    commandPath: `/proc/${process.pid}/fd/${handle.fd}`,
                    handle,
                };
            } catch (error) {
                await handle.close();
                throw error;
            }
        }

        async function repository(requested = ".") {
            const selected = string(requested, "repo", 4096);
            if (path.isAbsolute(selected)) {
                throw new TypeError("repo must be relative to the workspace");
            }
            const candidate = path.resolve(ctx.resolveWorkspace(selected));
            if (escapes(workspaceRoot, candidate)) {
                throw new RangeError(
                    "the Git repository is outside the workspace",
                );
            }
            const selectedDirectory = await bindDirectory(
                candidate,
                "the selected Git directory",
            );
            let discovery;
            try {
                discovery = await execute(
                    selectedDirectory.commandPath,
                    ["rev-parse", "--show-toplevel"],
                    { maxBytes: 16 * 1024 },
                );
            } finally {
                await selectedDirectory.handle.close();
            }
            const state = executionState(discovery, "git rev-parse");
            if (state.outputLimitExceeded) {
                throw new Error(
                    "Git repository discovery output was truncated",
                );
            }
            if (state.code !== 0) {
                throw new Error(
                    `${selected} is not inside a Git worktree: ` +
                        clipUtf8(resultText(discovery, "stderr"), 4096).text,
                );
            }
            const discovered = path.resolve(
                resultText(discovery, "stdout").trim(),
            );
            const directory = await bindDirectory(
                discovered,
                "the Git repository",
            );
            const relative = path.relative(workspaceRoot, directory.absolute);
            return {
                absolute: directory.absolute,
                close: () => directory.handle.close(),
                commandPath: directory.commandPath,
                relative: relative.split(path.sep).join("/") || ".",
            };
        }

        async function withRepository(requested, operation) {
            const repo = await repository(requested);
            try {
                return await operation(repo);
            } finally {
                await repo.close();
            }
        }

        function revision(value, label) {
            const selected = string(value, label, 256);
            if (selected.startsWith("-")) {
                throw new TypeError(`${label} must not begin with a hyphen`);
            }
            return selected;
        }

        async function resolveRevision(repo, value, label) {
            const selected = revision(value, label);
            const result = await execute(
                repo,
                [
                    "rev-parse",
                    "--verify",
                    "--quiet",
                    "--end-of-options",
                    `${selected}^{commit}`,
                ],
                { maxBytes: 16 * 1024 },
            );
            const state = executionState(result, `git rev-parse for ${label}`);
            if (state.outputLimitExceeded) {
                throw new Error(
                    `Git output was truncated while resolving ${label}`,
                );
            }
            if (state.code !== 0) {
                throw new Error(`${label} does not resolve to a commit`);
            }
            const hash = resultText(result, "stdout").trim();
            if (!/^[0-9a-f]{40,64}$/i.test(hash)) {
                throw new Error(`Git returned an invalid hash for ${label}`);
            }
            return hash;
        }

        function selectedPaths(repo, values) {
            if (values === undefined) {
                return [];
            }
            if (!Array.isArray(values) || values.length > 64) {
                throw new RangeError("paths must contain at most 64 entries");
            }
            return values.map((value, index) => {
                const selected = string(value, `paths[${index}]`, 4096);
                if (path.isAbsolute(selected)) {
                    throw new TypeError(`paths[${index}] must be relative`);
                }
                const absolute = path.resolve(repo, selected);
                const relative = path.relative(repo, absolute);
                if (escapes(repo, absolute)) {
                    throw new RangeError(
                        `paths[${index}] resolves outside the repo`,
                    );
                }
                const workspaceRelative = path.relative(
                    workspaceRoot,
                    absolute,
                );
                ctx.resolveWorkspace(workspaceRelative);
                return relative.split(path.sep).join("/") || ".";
            });
        }

        function parseStatus(value, maximum, outputLimitExceeded = false) {
            const framing = separatedFields(value, "\0");
            const records = framing.fields;
            const branch = {};
            const changes = [];
            let truncated = framing.malformed;
            for (let index = 0; index < records.length; index += 1) {
                const record = records[index];
                if (!record) {
                    continue;
                }
                if (record.startsWith("# branch.oid ")) {
                    branch.oid =
                        record.slice(13) === "(initial)"
                            ? null
                            : record.slice(13);
                    continue;
                }
                if (record.startsWith("# branch.head ")) {
                    branch.head =
                        record.slice(14) === "(detached)"
                            ? null
                            : record.slice(14);
                    continue;
                }
                if (record.startsWith("# branch.upstream ")) {
                    branch.upstream = record.slice(18);
                    continue;
                }
                if (record.startsWith("# branch.ab ")) {
                    const match = /^# branch\.ab \+(\d+) -(\d+)$/.exec(record);
                    if (match) {
                        branch.ahead = Number(match[1]);
                        branch.behind = Number(match[2]);
                    }
                    continue;
                }

                let consumed = 0;
                let change;
                let match = /^1 (\S+) \S+ \S+ \S+ \S+ \S+ \S+ ([\s\S]*)$/.exec(
                    record,
                );
                if (match) {
                    change = {
                        kind: "ordinary",
                        status: match[1],
                        path: match[2],
                    };
                } else {
                    match =
                        /^2 (\S+) \S+ \S+ \S+ \S+ \S+ \S+ \S+ ([\s\S]*)$/.exec(
                            record,
                        );
                }
                if (!change && match) {
                    const originalPath = records[index + 1] ?? null;
                    if (originalPath !== null) {
                        consumed = 1;
                        change = {
                            kind: "rename_or_copy",
                            status: match[1],
                            path: match[2],
                            originalPath,
                        };
                    }
                }
                if (!change) {
                    match =
                        /^u (\S+) \S+ \S+ \S+ \S+ \S+ \S+ \S+ \S+ ([\s\S]*)$/.exec(
                            record,
                        );
                    if (match) {
                        change = {
                            kind: "unmerged",
                            status: match[1],
                            path: match[2],
                        };
                    }
                }
                if (!change) {
                    match = /^([?!]) ([\s\S]*)$/.exec(record);
                    if (match) {
                        change = {
                            kind: match[1] === "?" ? "untracked" : "ignored",
                            status: match[1],
                            path: match[2],
                        };
                    }
                }
                if (!change) {
                    truncated = true;
                    continue;
                }
                if (changes.length >= maximum) {
                    truncated = true;
                    break;
                }
                changes.push(change);
                index += consumed;
            }
            return {
                branch,
                changes,
                truncated,
                framingTruncated: framing.malformed || outputLimitExceeded,
            };
        }

        function parseLog(value, maximum, outputLimitExceeded = false) {
            const framing = separatedFields(value, "\0");
            const fields = framing.fields;
            const commits = [];
            let truncated = framing.malformed || outputLimitExceeded;
            for (let index = 0; index + 4 < fields.length; index += 5) {
                const hash = fields[index];
                const parents = fields[index + 1];
                const author = fields[index + 2];
                const authoredAt = fields[index + 3];
                const subject = fields[index + 4];
                if (
                    !/^[0-9a-f]{40,64}$/i.test(hash) ||
                    (parents &&
                        !parents.split(" ").every((parent) => {
                            return /^[0-9a-f]{40,64}$/i.test(parent);
                        })) ||
                    !authoredAt
                ) {
                    truncated = true;
                    continue;
                }
                const clippedAuthor = clipUtf8(author, 512);
                const clippedSubject = clipUtf8(subject, 4096);
                truncated ||=
                    clippedAuthor.truncated || clippedSubject.truncated;
                commits.push({
                    hash,
                    parents: parents ? parents.split(" ") : [],
                    author: clippedAuthor.text,
                    authoredAt,
                    subject: clippedSubject.text,
                });
            }
            if (fields.length % 5 !== 0) {
                truncated = true;
            }
            if (commits.length > maximum) {
                truncated = true;
            }
            return {
                commits: commits.slice(0, maximum),
                truncated,
            };
        }

        function parseNameStatus(value, outputLimitExceeded = false) {
            const framing = separatedFields(value, "\0");
            const fields = framing.fields;
            const files = [];
            let truncated = framing.malformed;
            for (let index = 0; index < fields.length; index += 1) {
                const status = fields[index];
                if (!status) {
                    continue;
                }
                let file;
                let consumed;
                if (/^[RC]\d+$/.test(status)) {
                    const oldPath = fields[index + 1];
                    const selectedPath = fields[index + 2];
                    if (
                        oldPath === undefined ||
                        oldPath === "" ||
                        selectedPath === undefined ||
                        selectedPath === ""
                    ) {
                        truncated = true;
                        break;
                    }
                    file = { status, oldPath, path: selectedPath };
                    consumed = 2;
                } else {
                    const selectedPath = fields[index + 1];
                    if (selectedPath === undefined || selectedPath === "") {
                        truncated = true;
                        break;
                    }
                    file = { status, path: selectedPath };
                    consumed = 1;
                }
                if (files.length >= 1000) {
                    truncated = true;
                    break;
                }
                files.push(file);
                index += consumed;
            }
            return { files, truncated };
        }

        ctx.tool(
            {
                name: "git_snapshot",
                description:
                    "Return bounded branch, status, history, and diff-stat data " +
                    "without modifying Git.",
                inputSchema: {
                    type: "object",
                    additionalProperties: false,
                    properties: {
                        repo: {
                            type: "string",
                            default: ".",
                            maxLength: 4096,
                        },
                        log_limit: {
                            type: "integer",
                            minimum: 0,
                            maximum: 100,
                            default: 10,
                        },
                        max_changes: {
                            type: "integer",
                            minimum: 1,
                            maximum: 1000,
                            default: 200,
                        },
                        include_diff_stat: { type: "boolean", default: true },
                    },
                },
            },
            async (input = {}) => {
                return withRepository(input.repo ?? ".", async (repo) => {
                    const logLimit = integer(
                        input.log_limit,
                        10,
                        0,
                        100,
                        "log_limit",
                    );
                    const maxChanges = integer(
                        input.max_changes,
                        200,
                        1,
                        1000,
                        "max_changes",
                    );
                    const commands = [
                        execute(
                            repo.commandPath,
                            [
                                "status",
                                "--porcelain=v2",
                                "--branch",
                                "-z",
                                "--untracked-files=all",
                            ],
                            { maxBytes: 1024 * 1024 },
                        ),
                    ];
                    if (logLimit > 0) {
                        commands.push(
                            execute(
                                repo.commandPath,
                                [
                                    "log",
                                    "-z",
                                    "--no-show-signature",
                                    `--max-count=${logLimit + 1}`,
                                    "--format=%H%x00%P%x00%aN%x00%aI%x00%s",
                                ],
                                { maxBytes: 512 * 1024 },
                            ),
                        );
                    }
                    if (input.include_diff_stat !== false) {
                        commands.push(
                            execute(
                                repo.commandPath,
                                [
                                    "diff",
                                    "--no-ext-diff",
                                    "--no-textconv",
                                    "--stat=120,80",
                                ],
                                { maxBytes: 128 * 1024 },
                            ),
                        );
                        commands.push(
                            execute(
                                repo.commandPath,
                                [
                                    "diff",
                                    "--cached",
                                    "--no-ext-diff",
                                    "--no-textconv",
                                    "--stat=120,80",
                                ],
                                { maxBytes: 128 * 1024 },
                            ),
                        );
                    }

                    const results = await Promise.all(commands);
                    const statusResult = results.shift();
                    const statusState = executionState(
                        statusResult,
                        "git status",
                    );
                    if (
                        !statusState.outputLimitExceeded &&
                        statusState.code !== 0
                    ) {
                        throw new Error(
                            `git status failed: ` +
                                clipUtf8(
                                    resultText(statusResult, "stderr"),
                                    4096,
                                ).text,
                        );
                    }
                    const status = parseStatus(
                        resultText(statusResult, "stdout"),
                        maxChanges,
                        statusState.outputLimitExceeded,
                    );
                    const output = {
                        repo: repo.relative,
                        branch: status.branch,
                        changes: status.changes,
                        branchTruncated: status.framingTruncated,
                        changesTruncated:
                            status.truncated ||
                            statusState.outputLimitExceeded,
                        statusTruncated:
                            status.truncated ||
                            statusState.outputLimitExceeded,
                    };
                    if (logLimit > 0) {
                        const logResult = results.shift();
                        const logState = executionState(logResult, "git log");
                        if (
                            !logState.outputLimitExceeded &&
                            logState.code !== 0
                        ) {
                            const historyError = clipUtf8(
                                resultText(logResult, "stderr"),
                                4096,
                            );
                            output.commits = [];
                            output.historyAvailable = false;
                            output.historyTruncated = false;
                            output.historyError = historyError.text;
                            output.historyErrorTruncated =
                                historyError.truncated;
                        } else {
                            const history = parseLog(
                                resultText(logResult, "stdout"),
                                logLimit,
                                logState.outputLimitExceeded,
                            );
                            output.commits = history.commits;
                            output.historyAvailable = true;
                            output.historyTruncated = history.truncated;
                        }
                    }
                    if (input.include_diff_stat !== false) {
                        const workingResult = results.shift();
                        const stagedResult = results.shift();
                        const workingState = executionState(
                            workingResult,
                            "git diff --stat",
                        );
                        const stagedState = executionState(
                            stagedResult,
                            "git diff --cached --stat",
                        );
                        if (
                            !workingState.outputLimitExceeded &&
                            workingState.code !== 0
                        ) {
                            throw new Error(
                                "git diff --stat failed: " +
                                    clipUtf8(
                                        resultText(workingResult, "stderr"),
                                        4096,
                                    ).text,
                            );
                        }
                        if (
                            !stagedState.outputLimitExceeded &&
                            stagedState.code !== 0
                        ) {
                            throw new Error(
                                "git diff --cached --stat failed: " +
                                    clipUtf8(
                                        resultText(stagedResult, "stderr"),
                                        4096,
                                    ).text,
                            );
                        }
                        const working = clipUtf8(
                            resultText(workingResult, "stdout"),
                            128 * 1024,
                        );
                        const staged = clipUtf8(
                            resultText(stagedResult, "stdout"),
                            128 * 1024,
                        );
                        output.diffStat = {
                            working: working.text,
                            workingTruncated:
                                working.truncated ||
                                workingState.outputLimitExceeded,
                            staged: staged.text,
                            stagedTruncated:
                                staged.truncated ||
                                stagedState.outputLimitExceeded,
                        };
                    }
                    return output;
                });
            },
        );

        ctx.tool(
            {
                name: "git_compare",
                description:
                    "Compare two commits with bounded name-status and unified " +
                    "diff output.",
                inputSchema: {
                    type: "object",
                    additionalProperties: false,
                    required: ["base", "head"],
                    properties: {
                        repo: {
                            type: "string",
                            default: ".",
                            maxLength: 4096,
                        },
                        base: { type: "string", minLength: 1, maxLength: 256 },
                        head: { type: "string", minLength: 1, maxLength: 256 },
                        paths: {
                            type: "array",
                            maxItems: 64,
                            items: {
                                type: "string",
                                minLength: 1,
                                maxLength: 4096,
                            },
                        },
                        context_lines: {
                            type: "integer",
                            minimum: 0,
                            maximum: 100,
                            default: 3,
                        },
                        max_bytes: {
                            type: "integer",
                            minimum: 1024,
                            maximum: MAX_DIFF_BYTES,
                            default: 256 * 1024,
                        },
                    },
                },
            },
            async (input) => {
                return withRepository(input.repo ?? ".", async (repo) => {
                    const [base, head] = await Promise.all([
                        resolveRevision(repo.commandPath, input.base, "base"),
                        resolveRevision(repo.commandPath, input.head, "head"),
                    ]);
                    const paths = selectedPaths(repo.absolute, input.paths);
                    const context = integer(
                        input.context_lines,
                        3,
                        0,
                        100,
                        "context_lines",
                    );
                    const maximum = integer(
                        input.max_bytes,
                        256 * 1024,
                        1024,
                        MAX_DIFF_BYTES,
                        "max_bytes",
                    );
                    const separator = paths.length > 0 ? ["--", ...paths] : [];
                    const [namesResult, diffResult] = await Promise.all([
                        execute(
                            repo.commandPath,
                            [
                                "diff",
                                "--name-status",
                                "-z",
                                "--no-ext-diff",
                                "--no-textconv",
                                base,
                                head,
                                ...separator,
                            ],
                            { maxBytes: 256 * 1024 },
                        ),
                        execute(
                            repo.commandPath,
                            [
                                "diff",
                                "--no-color",
                                "--no-ext-diff",
                                "--no-textconv",
                                `--unified=${context}`,
                                base,
                                head,
                                ...separator,
                            ],
                            { maxBytes: maximum, timeoutMs: 30_000 },
                        ),
                    ]);
                    const namesState = executionState(
                        namesResult,
                        "git diff --name-status",
                    );
                    const diffState = executionState(diffResult, "git diff");
                    if (
                        (!namesState.outputLimitExceeded &&
                            namesState.code !== 0) ||
                        (!diffState.outputLimitExceeded &&
                            diffState.code !== 0)
                    ) {
                        throw new Error(
                            "git diff failed: " +
                                clipUtf8(
                                    resultText(diffResult, "stderr") ||
                                        resultText(namesResult, "stderr"),
                                    4096,
                                ).text,
                        );
                    }
                    const diff = clipUtf8(
                        resultText(diffResult, "stdout"),
                        maximum,
                    );
                    const names = parseNameStatus(
                        resultText(namesResult, "stdout"),
                        namesState.outputLimitExceeded,
                    );
                    const filesTruncated =
                        names.truncated || namesState.outputLimitExceeded;
                    const diffTruncated =
                        diff.truncated || diffState.outputLimitExceeded;
                    return {
                        repo: repo.relative,
                        base,
                        head,
                        paths,
                        files: names.files,
                        filesTruncated,
                        diff: diff.text,
                        diffBytes: diff.bytes,
                        diffTruncated,
                        truncated: filesTruncated || diffTruncated,
                    };
                });
            },
        );
    },
};

plugin.apply.description = description;

export default plugin;

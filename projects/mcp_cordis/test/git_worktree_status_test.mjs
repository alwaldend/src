import assert from "node:assert/strict";
import { mkdir, mkdtemp } from "node:fs/promises";
import { join, resolve } from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";
import test from "node:test";
import plugin from "../plugins/git_worktree.mjs";

const workspaceRoot = fileURLToPath(new URL(
    "../plugins/",
    import.meta.url,
));
const objectHash = "0".repeat(40);
const temporaryRoot = process.env.TEST_TMPDIR;
assert.ok(temporaryRoot, "Bazel must provide TEST_TMPDIR");

const fixtures = [
    {
        name: "ordinary",
        entries: [
            {
                records: [
                    `1 .M N... 100644 100644 100644 ${objectHash} ` +
                        `${objectHash} ordinary-one`,
                ],
                expected: {
                    kind: "ordinary",
                    status: ".M",
                    path: "ordinary-one",
                },
            },
            {
                records: [
                    `1 M. N... 100644 100644 100644 ${objectHash} ` +
                        `${objectHash} ordinary-two`,
                ],
                expected: {
                    kind: "ordinary",
                    status: "M.",
                    path: "ordinary-two",
                },
            },
        ],
    },
    {
        name: "rename and copy",
        entries: [
            {
                records: [
                    `2 R. N... 100644 100644 100644 ${objectHash} ` +
                        `${objectHash} R100 renamed-one`,
                    "original-one",
                ],
                expected: {
                    kind: "rename_or_copy",
                    status: "R.",
                    path: "renamed-one",
                    originalPath: "original-one",
                },
            },
            {
                records: [
                    `2 C. N... 100644 100644 100644 ${objectHash} ` +
                        `${objectHash} C100 copied-two`,
                    "original-two",
                ],
                expected: {
                    kind: "rename_or_copy",
                    status: "C.",
                    path: "copied-two",
                    originalPath: "original-two",
                },
            },
        ],
    },
    {
        name: "unmerged",
        entries: [
            {
                records: [
                    `u UU N... 100644 100644 100644 100644 ` +
                        `${objectHash} ${objectHash} ${objectHash} ` +
                        "unmerged-one",
                ],
                expected: {
                    kind: "unmerged",
                    status: "UU",
                    path: "unmerged-one",
                },
            },
            {
                records: [
                    `u AA N... 100644 100644 100644 100644 ` +
                        `${objectHash} ${objectHash} ${objectHash} ` +
                        "unmerged-two",
                ],
                expected: {
                    kind: "unmerged",
                    status: "AA",
                    path: "unmerged-two",
                },
            },
        ],
    },
    {
        name: "untracked",
        entries: [
            {
                records: ["? untracked-one"],
                expected: {
                    kind: "untracked",
                    status: "?",
                    path: "untracked-one",
                },
            },
            {
                records: ["? untracked-two"],
                expected: {
                    kind: "untracked",
                    status: "?",
                    path: "untracked-two",
                },
            },
        ],
    },
    {
        name: "ignored",
        entries: [
            {
                records: ["! ignored-one"],
                expected: {
                    kind: "ignored",
                    status: "!",
                    path: "ignored-one",
                },
            },
            {
                records: ["! ignored-two"],
                expected: {
                    kind: "ignored",
                    status: "!",
                    path: "ignored-two",
                },
            },
        ],
    },
];

function commandResult(stdout = "", options = {}) {
    return {
        code: 0,
        outputLimitExceeded: false,
        signal: null,
        stderr: "",
        stdout,
        truncated: false,
        ...options,
    };
}

function outputLimited(stdout = "", stderr = "") {
    return commandResult(stdout, {
        code: null,
        outputLimitExceeded: true,
        signal: "SIGTERM",
        stderr,
        truncated: true,
    });
}

function plain(value) {
    return JSON.parse(JSON.stringify(value));
}

function statusOutput(entries) {
    const records = [
        `# branch.oid ${objectHash}`,
        "# branch.head main",
        ...entries.flatMap((entry) => entry.records),
    ];
    return `${records.join("\0")}\0`;
}

function logRecord({
    hash = objectHash,
    parents = "",
    author = "Example Author",
    authoredAt = "2026-08-30T00:00:00Z",
    subject = "subject",
} = {}) {
    return [hash, parents, author, authoredAt, subject].join("\0") + "\0";
}

function loadHandlers(execute, calls = [], selectedRoot = workspaceRoot) {
    const handlers = new Map();
    plugin.apply({
        exec: async (file, args, options) => {
            assert.equal(file, "git");
            calls.push({ args: [...args], options: { ...options } });
            return execute(args, options);
        },
        resolveWorkspace: (relativePath) => {
            return resolve(selectedRoot, relativePath);
        },
        tool: (definition, handler) => {
            handlers.set(definition.name, handler);
        },
    });
    handlers.calls = calls;
    return handlers;
}

function loadSnapshotHandler(output, statusOptions = {}) {
    return loadHandlers((args) => {
        if (args.includes("--show-toplevel")) {
            return commandResult(`${workspaceRoot}\n`);
        }
        if (args.includes("status")) {
            return commandResult(output, statusOptions);
        }
        throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
    }).get("git_snapshot");
}

function snapshotRequest(maxChanges) {
    return {
        repo: ".",
        log_limit: 0,
        max_changes: maxChanges,
        include_diff_stat: false,
    };
}

test("git_snapshot bounds every porcelain-v2 record kind", async (t) => {
    for (const fixture of fixtures) {
        await t.test(fixture.name, async () => {
            const expected = fixture.entries.map((entry) => entry.expected);
            const handler = loadSnapshotHandler(statusOutput(fixture.entries));

            const limited = await handler(snapshotRequest(1));
            assert.deepEqual(plain(limited.changes), expected.slice(0, 1));
            assert.equal(limited.changesTruncated, true);
            assert.equal(limited.statusTruncated, true);
            assert.equal(limited.branchTruncated, false);

            const complete = await handler(snapshotRequest(2));
            assert.deepEqual(plain(complete.changes), expected);
            assert.equal(complete.changesTruncated, false);
            assert.equal(complete.statusTruncated, false);
            assert.equal(complete.branchTruncated, false);

            const exact = await loadSnapshotHandler(statusOutput(
                fixture.entries.slice(0, 1),
            ))(snapshotRequest(1));
            assert.deepEqual(plain(exact.changes), expected.slice(0, 1));
            assert.equal(exact.changesTruncated, false);
            assert.equal(exact.statusTruncated, false);
            assert.equal(exact.branchTruncated, false);
        });
    }
});

test("git_snapshot preserves newlines in NUL-delimited paths", async () => {
    const entries = [
        {
            records: [
                `1 .M N... 100644 100644 100644 ${objectHash} ` +
                    `${objectHash} ordinary\npath`,
            ],
            expected: {
                kind: "ordinary",
                status: ".M",
                path: "ordinary\npath",
            },
        },
        {
            records: [
                `2 R. N... 100644 100644 100644 ${objectHash} ` +
                    `${objectHash} R100 renamed\npath`,
                "original\npath",
            ],
            expected: {
                kind: "rename_or_copy",
                status: "R.",
                path: "renamed\npath",
                originalPath: "original\npath",
            },
        },
        {
            records: ["? untracked\npath"],
            expected: {
                kind: "untracked",
                status: "?",
                path: "untracked\npath",
            },
        },
    ];
    const snapshot = await loadSnapshotHandler(statusOutput(entries))(
        snapshotRequest(entries.length),
    );

    assert.deepEqual(
        plain(snapshot.changes),
        entries.map((entry) => entry.expected),
    );
    assert.equal(snapshot.changesTruncated, false);
});

test(
    "git_snapshot reports host truncation and drops partial records",
    async () => {
        const completeEntry = fixtures[0].entries[0];
        const incompleteRename =
            `2 R. N... 100644 100644 100644 ${objectHash} ` +
            `${objectHash} R100 missing-original`;
        const status = statusOutput([completeEntry]) +
            `${incompleteRename}\0`;
        const completeCommit = [
            objectHash,
            "",
            "Example Author",
            "2026-08-30T00:00:00Z",
            "complete subject",
        ].join("\0") + "\0";
        const handlers = loadHandlers((args) => {
            if (args.includes("--show-toplevel")) {
                return commandResult(`${workspaceRoot}\n`);
            }
            if (args.includes("status")) {
                return outputLimited(status);
            }
            if (args.includes("log")) {
                return outputLimited(`${completeCommit}partial commit`);
            }
            if (args.includes("--stat=120,80")) {
                if (args.includes("--cached")) {
                    return commandResult("staged complete\n");
                }
                    return outputLimited("working complete\npartial");
            }
            throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
        });
        const snapshot = await handlers.get("git_snapshot")({
            repo: ".",
            log_limit: 2,
            max_changes: 10,
            include_diff_stat: true,
        });

        assert.deepEqual(plain(snapshot.changes), [completeEntry.expected]);
        assert.equal(snapshot.changesTruncated, true);
        assert.equal(snapshot.statusTruncated, true);
        assert.equal(snapshot.branchTruncated, true);
        assert.equal(snapshot.commits.length, 1);
        assert.equal(snapshot.commits[0].subject, "complete subject");
        assert.equal(snapshot.historyAvailable, true);
        assert.equal(snapshot.historyTruncated, true);
        assert.equal(snapshot.diffStat.working, "working complete\npartial");
        assert.equal(snapshot.diffStat.workingTruncated, true);
        assert.equal(snapshot.diffStat.staged, "staged complete\n");
        assert.equal(snapshot.diffStat.stagedTruncated, false);
    },
);

test("git_snapshot rejects truncated repository discovery", async () => {
    const handler = loadHandlers((args) => {
        assert.ok(args.includes("--show-toplevel"));
        return outputLimited(workspaceRoot);
    }).get("git_snapshot");

    await assert.rejects(
        handler(snapshotRequest(1)),
        /repository discovery output was truncated/,
    );
});

function loadCompareHandler(namesResult, diffResult, revisionOptions = {}) {
    return loadHandlers((args) => {
        if (args.includes("--show-toplevel")) {
            return commandResult(`${workspaceRoot}\n`);
        }
        if (args.includes("--verify")) {
            const label = args.at(-1).startsWith("base") ? "base" : "head";
            return commandResult(`${objectHash}\n`, revisionOptions[label]);
        }
        if (args.includes("--name-status")) {
            return namesResult;
        }
        if (args.includes("--no-color")) {
            return diffResult;
        }
        throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
    }).get("git_compare");
}

function compareRequest() {
    return {
        repo: ".",
        base: "base",
        head: "head",
        max_bytes: 1024,
    };
}

test("git_compare rejects truncated revision discovery", async () => {
    const handler = loadCompareHandler(
        commandResult(""),
        commandResult(""),
        {
            base: {
                code: null,
                outputLimitExceeded: true,
                signal: "SIGTERM",
                truncated: true,
            },
        },
    );

    await assert.rejects(
        handler(compareRequest()),
        /output was truncated while resolving base/,
    );
});

test("git_compare separates files and diff truncation", async (t) => {
    await t.test("host-truncated name status", async () => {
        const names = [
            "M",
            "complete.txt",
            "R100",
            "renamed.txt",
        ].join("\0") + "\0partial-original";
        const comparison = await loadCompareHandler(
            outputLimited(names),
            commandResult("complete diff\n"),
        )(compareRequest());

        assert.deepEqual(plain(comparison.files), [
            { status: "M", path: "complete.txt" },
        ]);
        assert.equal(comparison.filesTruncated, true);
        assert.equal(comparison.diffTruncated, false);
        assert.equal(comparison.truncated, true);
    });

    await t.test("host-truncated diff", async () => {
        const comparison = await loadCompareHandler(
            commandResult("M\0complete.txt\0"),
            outputLimited("partial diff"),
        )(compareRequest());

        assert.equal(comparison.filesTruncated, false);
        assert.equal(comparison.diffTruncated, true);
        assert.equal(comparison.truncated, true);
    });

    await t.test("complete outputs", async () => {
        const comparison = await loadCompareHandler(
            commandResult("M\0complete.txt\0"),
            commandResult("complete diff\n"),
        )(compareRequest());

        assert.equal(comparison.filesTruncated, false);
        assert.equal(comparison.diffTruncated, false);
        assert.equal(comparison.truncated, false);
    });
});

test("git_snapshot frames history independently of RS and US", async () => {
    const subject = `legal ${"\x1e"} record and ${"\x1f"} unit separators`;
    const handlers = loadHandlers((args) => {
        if (args.includes("--show-toplevel")) {
            return commandResult(`${workspaceRoot}\n`);
        }
        if (args.includes("status")) return commandResult(statusOutput([]));
        if (args.includes("log")) {
            assert.ok(args.includes("--max-count=2"));
            assert.ok(args.includes("--no-show-signature"));
            return commandResult(logRecord({ subject }));
        }
        throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
    });
    const snapshot = await handlers.get("git_snapshot")({
        ...snapshotRequest(1),
        log_limit: 1,
    });

    assert.equal(snapshot.commits.length, 1);
    assert.equal(snapshot.commits[0].subject, subject);
    assert.equal(snapshot.historyAvailable, true);
    assert.equal(snapshot.historyTruncated, false);
});

test("git_snapshot reports every history truncation cause", async (t) => {
    async function snapshotFor(logOutput, limit = 1) {
        const handlers = loadHandlers((args) => {
            if (args.includes("--show-toplevel")) {
                return commandResult(`${workspaceRoot}\n`);
            }
            if (args.includes("status")) {
                return commandResult(statusOutput([]));
            }
            if (args.includes("log")) return commandResult(logOutput);
            throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
        });
        return handlers.get("git_snapshot")({
            ...snapshotRequest(1),
            log_limit: limit,
        });
    }

    await t.test("an observed extra commit", async () => {
        const snapshot = await snapshotFor(
            logRecord() + logRecord({ hash: "1".repeat(40) }),
        );
        assert.equal(snapshot.commits.length, 1);
        assert.equal(snapshot.historyTruncated, true);
    });

    await t.test("locally clipped author and subject", async () => {
        const snapshot = await snapshotFor(logRecord({
            author: "a".repeat(513),
            subject: "s".repeat(4097),
        }));
        assert.equal(Buffer.byteLength(snapshot.commits[0].author), 512);
        assert.equal(Buffer.byteLength(snapshot.commits[0].subject), 4096);
        assert.equal(snapshot.historyTruncated, true);
    });

    await t.test("a malformed field group", async () => {
        const snapshot = await snapshotFor(logRecord() + "bad\0group\0");
        assert.equal(snapshot.commits.length, 1);
        assert.equal(snapshot.historyTruncated, true);
    });
});

test(
    "git_snapshot distinguishes nonzero history from complete history",
    async () => {
        const handlers = loadHandlers((args) => {
            if (args.includes("--show-toplevel")) {
                return commandResult(`${workspaceRoot}\n`);
            }
            if (args.includes("status")) {
                return commandResult(statusOutput([]));
            }
            if (args.includes("log")) {
                return commandResult("", { code: 128, stderr: "no history" });
            }
            throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
        });
        const snapshot = await handlers.get("git_snapshot")({
            ...snapshotRequest(1),
            log_limit: 1,
        });

        assert.deepEqual(plain(snapshot.commits), []);
        assert.equal(snapshot.historyAvailable, false);
        assert.equal(snapshot.historyTruncated, false);
        assert.equal(snapshot.historyError, "no history");
        assert.equal(snapshot.historyErrorTruncated, false);
    },
);

test("git_snapshot marks a clipped history error", async () => {
    const handlers = loadHandlers((args) => {
        if (args.includes("--show-toplevel")) {
            return commandResult(`${workspaceRoot}\n`);
        }
        if (args.includes("status")) return commandResult(statusOutput([]));
        if (args.includes("log")) {
            return commandResult("", {
                code: 128,
                stderr: "e".repeat(5000),
            });
        }
        throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
    });
    const snapshot = await handlers.get("git_snapshot")({
        ...snapshotRequest(1),
        log_limit: 1,
    });

    assert.equal(Buffer.byteLength(snapshot.historyError), 4096);
    assert.equal(snapshot.historyErrorTruncated, true);
});

test("git_snapshot rejects both diff-stat command failures", async (t) => {
    for (const failed of ["working", "staged"]) {
        await t.test(failed, async () => {
            const handlers = loadHandlers((args) => {
                if (args.includes("--show-toplevel")) {
                    return commandResult(`${workspaceRoot}\n`);
                }
                if (args.includes("status")) {
                    return commandResult(statusOutput([]));
                }
                if (args.includes("--stat=120,80")) {
                    const staged = args.includes("--cached");
                    if ((failed === "staged") === staged) {
                        return commandResult("", {
                            code: 2,
                            stderr: `${failed} failed`,
                        });
                    }
                    return commandResult("");
                }
                throw new Error(
                    `unexpected Git invocation: ${args.join(" ")}`,
                );
            });
            await assert.rejects(
                handlers.get("git_snapshot")({
                    ...snapshotRequest(1),
                    include_diff_stat: true,
                }),
                new RegExp(`${failed} failed`),
            );
        });
    }
});

test("git tools reject unexpected signal terminations", async (t) => {
    const signalResult = commandResult("partial", {
        code: null,
        signal: "SIGTERM",
    });
    const snapshotCases = ["status", "log", "working stat", "staged stat"];
    for (const failed of snapshotCases) {
        await t.test(failed, async () => {
            const handlers = loadHandlers((args) => {
                if (args.includes("--show-toplevel")) {
                    return commandResult(`${workspaceRoot}\n`);
                }
                if (args.includes("status")) {
                    return failed === "status"
                        ? signalResult
                        : commandResult(statusOutput([]));
                }
                if (args.includes("log")) {
                    return failed === "log"
                        ? signalResult
                        : commandResult(logRecord());
                }
                if (args.includes("--stat=120,80")) {
                    const kind = args.includes("--cached")
                        ? "staged stat"
                        : "working stat";
                    return failed === kind ? signalResult : commandResult("");
                }
                throw new Error(
                    `unexpected Git invocation: ${args.join(" ")}`,
                );
            });
            await assert.rejects(
                handlers.get("git_snapshot")({
                    ...snapshotRequest(1),
                    log_limit: failed === "status" ? 0 : 1,
                    include_diff_stat: failed.includes("stat"),
                }),
                /terminated by SIGTERM/,
            );
        });
    }

    for (const failed of ["names", "diff"]) {
        await t.test(failed, async () => {
            const handler = loadCompareHandler(
                failed === "names"
                    ? signalResult
                    : commandResult("M\0file\0"),
                failed === "diff"
                    ? signalResult
                    : commandResult("diff\n"),
            );
            await assert.rejects(handler(compareRequest()), /SIGTERM/);
        });
    }
});

test("git tools reject unmarked incomplete output", async () => {
    const handler = loadSnapshotHandler(statusOutput([]), {
        truncated: true,
    });
    await assert.rejects(
        handler(snapshotRequest(1)),
        /without an output-limit marker/,
    );
});

test("git_snapshot discards an unframed status record", async () => {
    const complete = statusOutput([]);
    const partial = `? partial-path`;
    const snapshot = await loadSnapshotHandler(complete + partial)(
        snapshotRequest(2),
    );

    assert.deepEqual(plain(snapshot.changes), []);
    assert.equal(snapshot.changesTruncated, true);
    assert.equal(snapshot.statusTruncated, true);
    assert.equal(snapshot.branchTruncated, true);
});

test("git tools reject ordinary command failures", async (t) => {
    await t.test("status", async () => {
        const handler = loadSnapshotHandler("", {
            code: 2,
            stderr: "status failed",
        });
        await assert.rejects(handler(snapshotRequest(1)), /status failed/);
    });

    for (const failed of ["names", "diff"]) {
        await t.test(failed, async () => {
            const failure = commandResult("", {
                code: 2,
                stderr: `${failed} failed`,
            });
            const handler = loadCompareHandler(
                failed === "names" ? failure : commandResult(""),
                failed === "diff" ? failure : commandResult(""),
            );
            await assert.rejects(
                handler(compareRequest()),
                new RegExp(`${failed} failed`),
            );
        });
    }
});

test(
    "Git commands use read-only options and verified directory handles",
    async () => {
        const calls = [];
        const handlers = loadHandlers((args) => {
            if (args.includes("--show-toplevel")) {
                return commandResult(`${workspaceRoot}\n`);
            }
            if (args.includes("--verify")) {
                return commandResult(`${objectHash}\n`);
            }
            if (args.includes("status")) return commandResult(statusOutput([]));
            if (args.includes("log")) return commandResult(logRecord());
            if (args.includes("--name-status")) return commandResult("");
            if (args.includes("diff")) return commandResult("");
            throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
        }, calls);

        await handlers.get("git_snapshot")({
            ...snapshotRequest(1),
            log_limit: 1,
            include_diff_stat: true,
        });
        await handlers.get("git_compare")(compareRequest());

        assert.ok(calls.length >= 9);
        for (const call of calls) {
            assert.match(
                call.args[1],
                new RegExp(`^/proc/${process.pid}/fd/\\d+$`),
            );
            assert.equal(call.options.allowTruncatedOutput, true);
            assert.equal(call.options.env.GIT_OPTIONAL_LOCKS, "0");
            assert.equal(call.options.env.GIT_NO_LAZY_FETCH, "1");
            assert.equal(call.options.env.GIT_TERMINAL_PROMPT, "0");
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
        }
    },
);

test("containment accepts '..config' but rejects parent escapes", async () => {
    const root = await mkdtemp(join(temporaryRoot, "git-containment-"));
    const legalRepo = resolve(root, "..config");
    await mkdir(legalRepo);
    const snapshotHandlers = loadHandlers((args) => {
        if (args.includes("--show-toplevel")) {
            return commandResult(`${legalRepo}\n`);
        }
        if (args.includes("status")) return commandResult(statusOutput([]));
        throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
    }, [], root);

    const legalSnapshot = await snapshotHandlers.get("git_snapshot")({
        ...snapshotRequest(1),
        repo: "..config",
    });
    assert.equal(legalSnapshot.repo, "..config");

    const compareHandlers = loadHandlers((args) => {
        if (args.includes("--show-toplevel")) {
            return commandResult(`${root}\n`);
        }
        if (args.includes("--verify")) {
            return commandResult(`${objectHash}\n`);
        }
        if (args.includes("--name-status")) return commandResult("");
        if (args.includes("--no-color")) return commandResult("");
        throw new Error(`unexpected Git invocation: ${args.join(" ")}`);
    }, [], root);

    const legalCompare = await compareHandlers.get("git_compare")({
        ...compareRequest(),
        paths: ["..config/file"],
    });
    assert.deepEqual(plain(legalCompare.paths), ["..config/file"]);

    await assert.rejects(
        compareHandlers.get("git_compare")({
            ...compareRequest(),
            paths: ["../outside"],
        }),
        /outside the repo/,
    );

    const outside = await mkdtemp(join(temporaryRoot, "git-outside-"));
    await mkdir(resolve(root, "inside"));
    const outsideHandlers = loadHandlers((args) => {
        assert.ok(args.includes("--show-toplevel"));
        return commandResult(`${outside}\n`);
    }, [], root);
    await assert.rejects(
        outsideHandlers.get("git_snapshot")({
            ...snapshotRequest(1),
            repo: "inside",
        }),
        /outside the workspace/,
    );
});

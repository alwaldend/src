const assert = require("node:assert/strict");
const { run } = require("./main.cjs");

// Treat the workspace as an opaque path, including shell metacharacters.
const workspace = "/checkout/repository with spaces; $(ignored)";
for (const status of [0, 7]) {
    const exitCode = run(workspace, (command, args, options) => {
        assert.equal(command, "bazel");
        assert.deepEqual(args, ["run", "--config=ci", "//tools/ci"]);
        assert.deepEqual(options, {
            cwd: workspace,
            stdio: "inherit",
            shell: false,
        });
        return { status };
    });
    assert.equal(exitCode, status);
}
assert.equal(
    run(workspace, () => ({ status: null, signal: "SIGTERM" })),
    1,
);
assert.equal(
    run(workspace, () => ({ error: new Error("missing executable") })),
    1,
);
assert.equal(
    run(undefined, () => assert.fail("must not start Bazel")),
    1,
);

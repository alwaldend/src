const assert = require("node:assert/strict");
const { run } = require("./main.cjs");

// Treat the workspace as an opaque path, including shell metacharacters.
const workspace = "/checkout/repository with spaces; $(ignored)";
for (const status of [0, 7]) {
    let launcher;
    const exitCode = run(workspace, (command, args, options) => {
        assert.deepEqual(options, {
            cwd: workspace,
            stdio: "inherit",
            shell: false,
        });
        if (!launcher) {
            assert.equal(command, "bazel");
            const [runCommand, config, scriptFlag, target] = args;
            assert.deepEqual(
                [runCommand, config, target],
                ["run", "--config=ci", "//tools/ci"],
            );
            assert.match(scriptFlag, /^--script_path=.+\/ci$/);
            launcher = scriptFlag.slice("--script_path=".length);
            return { status: 0 };
        }
        assert.equal(command, launcher);
        assert.deepEqual(args, []);
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

const { spawnSync } = require("node:child_process");

function run(workspace, spawn = spawnSync) {
    if (!workspace) {
        console.error("Repository CI requires a checked-out GITHUB_WORKSPACE");
        return 1;
    }
    const result = spawn("bazel", ["run", "--config=ci", "//tools/ci"], {
        cwd: workspace,
        stdio: "inherit",
        shell: false,
    });
    if (result.error) {
        console.error("Repository CI could not start Bazel");
        return 1;
    }
    return result.status ?? 1;
}

if (require.main === module) {
    process.exitCode = run(process.env.GITHUB_WORKSPACE);
}

module.exports = { run };

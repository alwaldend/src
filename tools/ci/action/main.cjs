const { execFileSync, spawnSync } = require("node:child_process");
const {
    existsSync,
    mkdirSync,
    mkdtempSync,
    realpathSync,
    rmSync,
} = require("node:fs");
const { tmpdir, userInfo } = require("node:os");
const { isAbsolute, join, relative, sep } = require("node:path");

function contains(parent, child) {
    const path = relative(parent, child);
    return (
        path === "" ||
        (!isAbsolute(path) && path !== ".." && !path.startsWith(`..${sep}`))
    );
}

function git(cwd, ...args) {
    return execFileSync("git", args, {
        cwd,
        encoding: "utf8",
        stdio: ["ignore", "pipe", "inherit"],
    }).trim();
}

function prepareStableCheckout(source, destination, revision) {
    source = realpathSync(source);
    if (
        git(source, "rev-parse", "--show-toplevel") !== source ||
        git(source, "rev-parse", "HEAD") !== revision
    ) {
        throw new Error("The job checkout must contain the event commit");
    }
    const lfsStorage = join(
        git(source, "rev-parse", "--absolute-git-dir"),
        "lfs",
    );
    mkdirSync(join(destination, ".."), { recursive: true, mode: 0o700 });
    const parent = realpathSync(join(destination, ".."));
    if (contains(source, parent) || contains(destination, source)) {
        throw new Error(
            "The stable checkout must be outside the job workspace",
        );
    }
    if (!existsSync(destination)) {
        mkdirSync(destination, { mode: 0o700 });
        git(destination, "init", "--template=", "--initial-branch=ci");
        git(destination, "lfs", "install", "--local");
        git(destination, "config", "--local", "ci.stable-checkout", "true");
    }
    if (
        git(destination, "rev-parse", "--show-toplevel") !== destination ||
        git(
            destination,
            "config",
            "--local",
            "--get",
            "ci.stable-checkout",
        ) !== "true"
    ) {
        throw new Error(
            "Refusing to clean a directory not managed by repository CI",
        );
    }

    git(
        destination,
        "fetch",
        "--atomic",
        "--prune",
        "--no-tags",
        "--",
        source,
        revision,
        "+refs/*:refs/*",
    );
    git(
        destination,
        "-c",
        `lfs.storage=${lfsStorage}`,
        "checkout",
        "--detach",
        "--force",
        revision,
    );
    git(destination, "clean", "-qffdx");
    return destination;
}

function withLogGroup(title, operation) {
    if (process.env.FORGEJO_ACTIONS !== "true") {
        return operation();
    }
    console.log(`##[group]${title}`);
    try {
        return operation();
    } finally {
        console.log("\n##[endgroup]");
    }
}

function run(workspace, spawn = spawnSync) {
    if (!workspace) {
        console.error("Repository CI requires a checked-out GITHUB_WORKSPACE");
        return 1;
    }
    const options = {
        cwd: workspace,
        stdio: "inherit",
        shell: false,
    };
    const directory = mkdtempSync(join(tmpdir(), "repository-ci-"));
    const launcher = join(directory, "ci");
    try {
        const buildStatus = withLogGroup("Build repository CI command", () => {
            const result = withLogGroup("Bazel launcher diagnostics", () =>
                spawn(
                    "bazel",
                    [
                        "run",
                        "--config=ci",
                        `--script_path=${launcher}`,
                        "//tools/ci",
                    ],
                    options,
                ),
            );
            if (result.error) {
                console.error("Repository CI could not start Bazel");
                return 1;
            }
            const status = result.status ?? 1;
            if (status !== 0) {
                console.log(
                    `Repository CI launcher build exited with status ${status}`,
                );
            }
            return status;
        });
        if (buildStatus !== 0) {
            return buildStatus;
        }
        const result = spawn(launcher, [], options);
        if (result.error) {
            console.error("Repository CI could not start its launcher");
            return 1;
        }
        return result.status ?? 1;
    } finally {
        rmSync(directory, { recursive: true, force: true });
    }
}

if (require.main === module) {
    try {
        let workspace = process.env.GITHUB_WORKSPACE;
        if (process.env["INPUT_STABLE-CHECKOUT"] === "true") {
            const repository = process.env.GITHUB_REPOSITORY ?? "";
            if (!/^[\w-][\w.-]*\/[\w-][\w.-]*$/.test(repository)) {
                throw new Error(
                    "Stable checkout requires an owner/repository name",
                );
            }
            const home = realpathSync(userInfo().homedir);
            workspace = withLogGroup("Prepare repository checkout", () =>
                prepareStableCheckout(
                    workspace,
                    join(home, ".cache", "repository-ci", repository),
                    process.env.GITHUB_SHA,
                ),
            );
            process.env.GITHUB_WORKSPACE = workspace;
            process.env.FORGEJO_WORKSPACE = workspace;
        }
        process.exitCode = run(workspace);
    } catch (error) {
        console.error(
            `Repository CI could not prepare the checkout: ${error.message}`,
        );
        process.exitCode = 1;
    }
}

module.exports = { prepareStableCheckout, run };

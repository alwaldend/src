"""Analysis tests for Hugo build actions."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")

def _has_input_suffix(inputs, suffix):
    return any([file.path.endswith(suffix) for file in inputs])

def _hugo_site_action_test_impl(ctx):
    env = analysistest.begin(ctx)
    actions = analysistest.target_actions(env)

    run_actions = [
        action
        for action in actions
        if action.mnemonic == "HugoSite" or "hugo" in " ".join(action.argv).lower()
    ]
    asserts.equals(env, 1, len(run_actions))
    action = run_actions[0]
    inputs = action.inputs.to_list()
    site_archives = [
        file
        for file in inputs
        if file.basename == "site_archive.tar"
    ]
    postcss_executables = [
        file
        for file in inputs
        if file.basename == "mock_postcss.sh" or
           "mock_postcss" in file.short_path
    ]

    asserts.equals(env, 1, len(site_archives))
    asserts.true(
        env,
        site_archives[0].path.startswith(
            analysistest.target_bin_dir_path(env) + "/",
        ),
        "site archive must stay in the target configuration",
    )
    asserts.true(
        env,
        len(postcss_executables) >= 1,
        "the action must declare the PostCSS executable",
    )
    command = " ".join(action.argv)
    asserts.true(
        env,
        "node_modules/.bin/postcss" in command,
        "the action must create the project-local PostCSS shim Hugo prefers",
    )
    asserts.true(
        env,
        "NODE_PATH" in command,
        "the PostCSS shim must expose its declared Node modules",
    )
    asserts.true(
        env,
        "hugo" in command,
        "Hugo must be executed by the action",
    )
    asserts.true(
        env,
        'DART_SASS_BINARY="sass"' in command or "DART_SASS_BINARY=sass" in command,
        "the site environment must define the bare DART_SASS_BINARY Hugo whitelists",
    )
    return analysistest.end(env)

hugo_site_action_test = analysistest.make(
    _hugo_site_action_test_impl,
)

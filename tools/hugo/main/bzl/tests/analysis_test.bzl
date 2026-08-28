"""Analysis tests for Hugo build actions."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")

def _has_input_suffix(inputs, suffix):
    return any([file.path.endswith(suffix) for file in inputs])

def _hugo_site_action_test_impl(ctx):
    env = analysistest.begin(ctx)
    actions = analysistest.target_actions(env)

    asserts.equals(env, 1, len(actions))
    inputs = actions[0].inputs.to_list()
    site_archives = [
        file
        for file in inputs
        if file.basename == "site_archive.tar"
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
        _has_input_suffix(inputs, "/tools/postcss/postcss_/postcss"),
        "PostCSS runfiles must be declared as action inputs",
    )
    asserts.true(
        env,
        _has_input_suffix(inputs, "/bin.exe"),
        "Hugo must be declared as an action tool",
    )
    asserts.true(
        env,
        _has_input_suffix(inputs, "/toolchain_impl.env.txt"),
        "generated Hugo environment must be declared as an action input",
    )
    return analysistest.end(env)

hugo_site_action_test = analysistest.make(
    _hugo_site_action_test_impl,
)

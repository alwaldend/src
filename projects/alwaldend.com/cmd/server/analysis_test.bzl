"""Analysis tests for the server image action."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")

def _ops_action_test_impl(ctx):
    env = analysistest.begin(ctx)
    actions = analysistest.target_actions(env)

    asserts.equals(env, 1, len(actions))
    server_inputs = [
        file
        for file in actions[0].inputs.to_list()
        if file.path.endswith(
            "/projects/alwaldend.com/cmd/server/server_/server",
        )
    ]
    asserts.equals(env, 1, len(server_inputs))
    asserts.true(
        env,
        server_inputs[0].path.startswith(
            analysistest.target_bin_dir_path(env) + "/",
        ),
        "the embedded server must stay in the target configuration",
    )
    return analysistest.end(env)

ops_action_test = analysistest.make(
    _ops_action_test_impl,
)

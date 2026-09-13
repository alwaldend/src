"""Analysis tests for the persistent-worker Hugo action."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")

def _hugo_worker_action_test_impl(ctx):
    env = analysistest.begin(ctx)
    actions = analysistest.target_actions(env)
    worker_actions = [
        action
        for action in actions
        if action.mnemonic == "HugoSite"
    ]
    asserts.equals(env, 1, len(worker_actions))
    action = worker_actions[0]
    inputs = action.inputs.to_list()
    flagfile = None
    for input in inputs:
        if input.basename.endswith(".flagfile.json"):
            flagfile = input
    asserts.true(
        env,
        flagfile != None,
        "the worker action must declare a flagfile input",
    )
    command = " ".join(action.argv)
    asserts.true(
        env,
        "--flagfile=" in command,
        "the worker must be invoked with the flagfile argument",
    )
    return analysistest.end(env)

hugo_worker_action_test = analysistest.make(
    _hugo_worker_action_test_impl,
)

"""Analysis tests for the Promptfoo skill workspace."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")
load(
    "//promptfoo/private:workspace.bzl",
    "PromptfooWorkspaceInfo",
)

def _workspace_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    workspace = target[PromptfooWorkspaceInfo]

    asserts.equals(env, ["fixture_skill"], workspace.skill_names)
    asserts.equals(
        env,
        [
            ".agents/skills/fixture_skill/SKILL.md",
            ".agents/skills/fixture_skill/reference.txt",
        ],
        sorted(workspace.files_by_path.keys()),
    )
    asserts.equals(env, 2, len(workspace.files.to_list()))
    asserts.equals(env, 2, len(target[DefaultInfo].files.to_list()))
    return analysistest.end(env)

workspace_test = analysistest.make(_workspace_test_impl)

def _workspace_failure_test_impl(ctx):
    env = analysistest.begin(ctx)
    asserts.expect_failure(
        env,
        "promptfoo workspace has duplicate skill name \"fixture_skill\"",
    )
    return analysistest.end(env)

workspace_failure_test = analysistest.make(
    _workspace_failure_test_impl,
    expect_failure = True,
)

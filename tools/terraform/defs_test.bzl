"""Analysis tests for Terraform Bazel macros."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")

def _terraform_test_map_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    executable = target[DefaultInfo].files_to_run.executable
    asserts.true(env, executable != None)
    asserts.true(env, str(target.label).endswith(":fixture.fmt_test"))
    return analysistest.end(env)

terraform_test_map_test = analysistest.make(_terraform_test_map_impl)

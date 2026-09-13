"""Analysis checks for public operation maps and caller-owned wrappers."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")
load("//terraform:defs.bzl", "TerraformRunInfo")

_InvocationInfo = provider("Observed Terraform invocation attributes.", fields = ["arguments", "kind", "testonly"])
_WrapperInfo = provider("Observed caller wrapper attributes and dependencies.", fields = ["arguments", "data", "mode", "size"])

def _invocation_aspect_impl(_target, ctx):
    return [_InvocationInfo(
        arguments = ctx.rule.attr.arguments,
        kind = ctx.rule.kind,
        testonly = ctx.rule.attr.testonly,
    )]

_invocation_aspect = aspect(implementation = _invocation_aspect_impl)

def _record_wrapper_impl(ctx):
    return [
        DefaultInfo(runfiles = ctx.runfiles(files = ctx.files.data).merge_all([
            target[DefaultInfo].default_runfiles
            for target in ctx.attr.data
        ])),
        _WrapperInfo(
            arguments = ctx.attr.args,
            data = ctx.attr.data,
            mode = ctx.attr.mode,
            size = ctx.attr.requested_size,
        ),
    ]

_record_wrapper = rule(
    implementation = _record_wrapper_impl,
    attrs = {
        "args": attr.string_list(),
        "data": attr.label_list(allow_files = True),
        "mode": attr.string(),
        "requested_size": attr.string(),
    },
)

def recording_wrapper(name, args, data, size = "", **kwargs):
    """Record the generic wrapper contract without executing an operation."""
    _record_wrapper(
        name = name,
        args = args,
        data = data,
        requested_size = size,
        **kwargs
    )

def _invocation_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    invocation = target[_InvocationInfo]
    asserts.equals(env, ctx.attr.expected_arguments, invocation.arguments)
    asserts.equals(env, ctx.attr.expected_kind, invocation.kind)
    asserts.true(env, invocation.testonly, "operation fixtures must remain test-only")
    asserts.true(env, TerraformRunInfo in target, "maps must use the public Terraform execution contract")
    return analysistest.end(env)

invocation_test = analysistest.make(
    _invocation_test_impl,
    attrs = {
        "expected_arguments": attr.string_list(),
        "expected_kind": attr.string(default = "terraform_binary"),
    },
    extra_target_under_test_aspects = [_invocation_aspect],
)

def _wrapper_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    wrapper = target[_WrapperInfo]
    invocation_name = target.label.name + ".terraform"
    asserts.equals(env, ["$(rootpath :{})".format(invocation_name)], wrapper.arguments)
    asserts.equals(env, "caller-authentication", wrapper.mode)
    asserts.equals(env, ctx.attr.expected_size, wrapper.size)
    asserts.equals(env, ["context.json", invocation_name], [dependency.label.name for dependency in wrapper.data])
    invocation = wrapper.data[-1]
    asserts.true(env, TerraformRunInfo in invocation, "wrapper must launch the declared Terraform target")
    inputs = [file.basename for file in invocation[DefaultInfo].default_runfiles.files.to_list()]
    asserts.true(env, "main.tf" in inputs, "configuration must be packaged by the inner invocation")
    asserts.true(env, "context.json" in inputs, "declared data must also reach the inner invocation")
    wrapped_files = target[DefaultInfo].default_runfiles.files.to_list()
    asserts.true(env, invocation[TerraformRunInfo].config in wrapped_files, "wrapper runfiles must retain the invocation configuration")
    return analysistest.end(env)

wrapper_test = analysistest.make(
    _wrapper_test_impl,
    attrs = {"expected_size": attr.string()},
)

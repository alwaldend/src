"""A build-input fixture with target inputs, tools, and visibility metadata."""

def _inputs_fixture_impl(ctx):
    return [DefaultInfo(files = depset(ctx.files.srcs + ctx.files._implicit_input))]

inputs_fixture = rule(
    implementation = _inputs_fixture_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True),
        "tools": attr.label_list(cfg = "exec"),
        "_implicit_input": attr.label(default = Label(":implicit_input")),
        "_exec_tool": attr.label(
            default = Label("//projects/bazel_deps_analyzer:bazel_deps_analyzer"),
            executable = True,
            cfg = "exec",
        ),
    },
)

filter_inputs_fixture = rule(
    implementation = _inputs_fixture_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True),
        "deps": attr.label_list(default = [Label(":implicit_input")]),
        "data": attr.label_list(allow_files = True),
        "tools": attr.label_list(cfg = "exec"),
        "_implicit_input": attr.label(default = Label(":implicit_input")),
    },
)

def _string_inputs_fixture_impl(_ctx):
    return [DefaultInfo(files = depset())]

string_inputs_fixture = rule(
    implementation = _string_inputs_fixture_impl,
    attrs = {"srcs": attr.string()},
)

load("//tools/release:al_release_files_info.bzl", "AlReleaseFilesInfo")

def _impl(ctx):
    manifest = ctx.actions.declare_file("{}.release.json".format(ctx.label.name))
    inputs = []
    outputs = [manifest]
    args = ctx.actions.args()
    args.add("generate")
    args.add_all(["--output", manifest.path])
    for src in ctx.files.srcs:
        args.add_all(["--item", "{}:{}".format("file", src.path)])
        inputs.append(src)
    ctx.actions.run(
        executable = ctx.executable.release_tool,
        arguments = [args],
        inputs = inputs,
        outputs = outputs,
    )
    return [
        DefaultInfo(
            files = depset(
                direct = ctx.files.srcs,
                transitive = [dep[DefaultInfo].files for dep in ctx.attr.deps],
            ),
        ),
        AlReleaseFilesInfo(
            manifest = manifest,
            srcs = depset(direct = ctx.files.srcs),
            deps = depset(transitive = [dep[AlReleaseFilesInfo] for dep in ctx.attr.deps]),
        ),
    ]

al_release_files = rule(
    implementation = _impl,
    provides = [AlReleaseFilesInfo],
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Sources",
        ),
        "deps": attr.label_list(
            providers = [AlReleaseFilesInfo],
            doc = "Deps",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = "//tools/release",
            doc = "Release tool",
        ),
    },
)

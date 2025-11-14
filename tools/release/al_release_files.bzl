load("//tools/release:al_release_files_info.bzl", "AlReleaseFilesInfo")

def _impl(ctx):
    manifest = ctx.actions.declare_file("{}.release.json".format(ctx.label.name))
    outputs = [manifest]
    inputs = []
    args = ctx.actions.args()
    args.add("generate")
    args.add_all(["--output", manifest.path])
    srcs = []
    for i, src in enumerate(ctx.files.srcs):
        ignore = False
        for suffix in ctx.attr.ignore_suffixes:
            if src.path.endswith(suffix):
                ignore = True
                break
        if ignore:
            continue

        item = ctx.actions.declare_file("{}.release_item.{}.json".format(ctx.label.name, i))
        inputs.extend([item, src])
        srcs.append(src)
        args.add_all(["--item", item.path])
        ctx.actions.write(
            output = item,
            content = json.encode(
                {
                    "file": {
                        "local_path": src.path,
                    },
                },
            ),
        )
    ctx.actions.run(
        executable = ctx.executable.release_tool,
        arguments = [args],
        inputs = inputs,
        outputs = outputs,
    )
    return [
        DefaultInfo(
            files = depset(
                direct = srcs,
                transitive = [dep[DefaultInfo].files for dep in ctx.attr.deps],
            ),
        ),
        AlReleaseFilesInfo(
            manifest = manifest,
            srcs = depset(srcs),
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
        "ignore_suffixes": attr.string_list(
            doc = "Ignore src files ending with these suffixes",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = "//tools/release",
            doc = "Release tool",
        ),
    },
)

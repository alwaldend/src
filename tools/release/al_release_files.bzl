load("//tools/release:al_release_files_info.bzl", "AlReleaseFilesInfo")

def _impl(ctx):
    manifest = ctx.actions.declare_file("{}.manifest.json".format(ctx.label.name))
    ctx.actions.write(
        output = manifest,
        content = "{}",
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
    },
)

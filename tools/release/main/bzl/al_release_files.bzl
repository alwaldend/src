load("//tools/release/main/bzl:al_release_deployment_info.bzl", "AlReleaseDeploymentInfo")
load("//tools/release/main/bzl:al_release_files_info.bzl", "AlReleaseFilesInfo")

def _contains_suffixes(val, suffixes):
    """
    Args:
        val (string): value to check
        suffixes (list[string]): suffixes to check against
    """
    for suffix in suffixes:
        if val.endswith(suffix):
            return True
    return False

def _impl(ctx):
    manifest = ctx.actions.declare_file("{}.release.json".format(ctx.label.name))
    outputs = [manifest]
    inputs = []
    args = ctx.actions.args()
    files = {}

    args.add("generate")
    args.add_all(["--output_manifest", manifest.path])

    for src in ctx.files.srcs:
        files[src.basename] = src
        if _contains_suffixes(src.basename, ctx.attr.ignore_suffixes):
            continue
        inputs.append(src)
        args.add_all(["--add_file", src.path])

    for deployment in ctx.attr.deployments:
        info = deployment[AlReleaseDeploymentInfo].info_file
        args.add_all(["--deployment", info.path])
        inputs.append(info)

    for dep in ctx.attr.deps:
        info = dep[AlReleaseFilesInfo]
        files.update(info.files)
        inputs.extend(info.values())
        args.add_all(["--merge_manifest", info.manifest.path])

    ctx.actions.run(
        executable = ctx.executable.release_tool,
        arguments = [args],
        inputs = inputs,
        outputs = outputs,
    )
    return [
        DefaultInfo(
            files = depset([manifest]),
        ),
        AlReleaseFilesInfo(
            manifest = manifest,
            files = files,
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
        "deployments": attr.label_list(
            providers = [AlReleaseDeploymentInfo],
            doc = "Deployment info",
        ),
        "ignore_suffixes": attr.string_list(
            doc = "Ignore src files ending with these suffixes",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = "//tools/release/main/go",
            doc = "Release tool",
        ),
    },
)

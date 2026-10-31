load("//rules:release_deployment_info.bzl", "ReleaseDeploymentInfo")
load("//rules:release_files_info.bzl", "ReleaseFilesInfo")

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

    args.add("release")
    args.add("generate")
    args.add_all(["--output_manifest", manifest.path])

    for src in ctx.files.srcs:
        if _contains_suffixes(src.basename, ctx.attr.ignore_suffixes):
            continue
        if src.basename in files:
            fail("duplicate release filename: " + src.basename)
        files[src.basename] = src
        inputs.append(src)
        args.add_all(["--add_file", src.path])

    for deployment in ctx.attr.deployments:
        info = deployment[ReleaseDeploymentInfo].info_file
        args.add_all(["--deployment", info.path])
        inputs.append(info)

    for dep in ctx.attr.deps:
        info = dep[ReleaseFilesInfo]
        for name, file in info.files.items():
            if name in files:
                fail("duplicate release filename: " + name)
            files[name] = file
        inputs.extend(info.files.values())
        inputs.append(info.manifest)
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
        ReleaseFilesInfo(
            manifest = manifest,
            files = files,
        ),
    ]

release_files = rule(
    implementation = _impl,
    provides = [ReleaseFilesInfo],
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Sources",
        ),
        "deps": attr.label_list(
            providers = [ReleaseFilesInfo],
            doc = "Deps",
        ),
        "deployments": attr.label_list(
            providers = [ReleaseDeploymentInfo],
            doc = "Deployment info",
        ),
        "ignore_suffixes": attr.string_list(
            doc = "Ignore src files ending with these suffixes",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = Label("//cmd/al:al"),
            doc = "Release tool",
        ),
    },
)

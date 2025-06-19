load("//bzl/providers:al_hugo_module_info.bzl", "AlHugoModuleInfo")

def _impl(ctx):
    return [
        DefaultInfo(
            files = depset(ctx.files.src),
        ),
        AlHugoModuleInfo(
            path = ctx.attr.path,
            archive = ctx.file.src,
        ),
    ]

al_hugo_module = rule(
    implementation = _impl,
    doc = "Describe hugo module",
    attrs = {
        "path": attr.string(doc = "Module path", mandatory = True),
        "src": attr.label(doc = "Module archive", mandatory = True, allow_single_file = [".tar"]),
    },
)

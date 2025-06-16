load("//bzl/providers:al_toml_info.bzl", "AlTomlInfo")

def _impl(ctx):
    return [
        DefaultInfo(files = depset(direct = ctx.files.srcs, transitive = ctx.files.deps)),
        AlTomlInfo(srcs = ctx.files.srcs, deps = ctx.files.deps),
    ]

al_toml_data = rule(
    implementation = _impl,
    attrs = {
        "srcs": attr.label_list(doc = "Toml files", allow_files = [".toml"]),
        "deps": attr.label_list(doc = "Toml data targets", providers = [AlTomlInfo]),
    },
)

load("//bzl/providers:al_toml_info.bzl", "AlTomlInfo")

def _impl(ctx):
    return [
        DefaultInfo(files = depset(direct = ctx.files.srcs)),
        AlTomlInfo(srcs = ctx.files.srcs, deps = ctx.files.deps),
    ]

al_toml_data = rule(
    implementation = _impl,
    provides = [AlTomlInfo, DefaultInfo],
    attrs = {
        "srcs": attr.label_list(doc = "Toml files", allow_files = [".toml"]),
        "deps": attr.label_list(doc = "Toml data targets", providers = [AlTomlInfo]),
        "tomlv": attr.label(
            executable = True,
            doc = "Tomlv label for the validation aspect",
            cfg = "exec",
            default = "//tools:tomlv",
        ),
    },
)

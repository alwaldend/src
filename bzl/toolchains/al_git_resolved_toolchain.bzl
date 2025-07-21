load(":al_resolved_toolchain.bzl", "al_resolved_toolchain")

al_git_resolved_toolchain = al_resolved_toolchain(
    doc = "Resolved git toolchain",
    toolchain_label = "//bzl/toolchain-types:git",
)

load("//tools/resolved_toolchain/main/bzl:al_resolved_toolchain.bzl", "al_resolved_toolchain")

al_git_resolved_toolchain = al_resolved_toolchain(
    doc = "Resolved git toolchain",
    toolchain_label = "//tools/git/main/bzl:toolchain_type",
)

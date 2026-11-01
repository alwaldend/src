load("//tools/resolved_toolchain/rules:al_resolved_toolchain.bzl", "al_resolved_toolchain")

al_git_resolved_toolchain = al_resolved_toolchain(
    doc = "Resolved git toolchain",
    toolchain_label = "//tools/git/rules:toolchain_type",
)

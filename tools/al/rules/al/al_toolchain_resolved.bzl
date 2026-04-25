load("//tools/resolved_toolchain/main/bzl:al_resolved_toolchain.bzl", "al_resolved_toolchain")

al_toolchain_resolved = al_resolved_toolchain(
    doc = "Resolved al toolchain (for genrules)",
    toolchain_label = "//tools/al/rules/al:toolchain_type",
)

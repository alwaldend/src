load("//bzl/rules/resolved_toolchain:al_resolved_toolchain.bzl", "al_resolved_toolchain")

al_qt_resolved_toolchain = al_resolved_toolchain(
    doc = "Resolved qt toolchain",
    toolchain_label = "//bzl/toolchain-types:qt",
)

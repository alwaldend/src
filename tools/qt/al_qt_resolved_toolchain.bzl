load("//tools/resolved_toolchain:al_resolved_toolchain.bzl", "al_resolved_toolchain")

al_qt_resolved_toolchain = al_resolved_toolchain(
    doc = "Resolved qt toolchain",
    toolchain_label = "//tools/qt:toolchain_type",
)

## Registry consumer verification

Bazel 9.3.1 built the ruleset archive from the monorepo root. A fresh workspace
with a local Bzlmod registry resolved `bazel_dep(name = "rules_iso", version =
"0.0.0")` using the release module declaration and a checksum of the generated
archive. The fixture used a small file as the ISO input.

The consumer built `@rules_iso//pkg/bzl:iso_extension`, `@fixture_iso//file:file`,
and `@fixture_iso_flash//:iso_flash`. The root workspace also resolved both
existing generated flash targets for XCP-ng and TrueNAS SCALE.

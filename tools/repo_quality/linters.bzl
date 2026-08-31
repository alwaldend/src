"""Repository linter aspects backed by upstream aspect_rules_lint."""

load("@aspect_rules_lint//lint:buildifier.bzl", "lint_buildifier_aspect")
load("@aspect_rules_lint//lint:ruff.bzl", "lint_ruff_aspect")
load("@aspect_rules_lint//lint:shellcheck.bzl", "lint_shellcheck_aspect")

buildifier = lint_buildifier_aspect(
    binary = Label("//tools/buildifier"),
    warnings = "-bzl-visibility,-function-docstring-args,-function-docstring-header,-module-docstring,-name-conventions,-no-effect,-print,-unnamed-macro,-unsorted-dict-items,-unused-variable",
)

ruff = lint_ruff_aspect(
    binary = Label("@aspect_rules_lint//lint:ruff_bin"),
    configs = [Label("//:pyproject.toml")],
)

shellcheck = lint_shellcheck_aspect(
    binary = Label("@aspect_rules_lint//lint:shellcheck_bin"),
    config = Label("//:.shellcheckrc"),
)

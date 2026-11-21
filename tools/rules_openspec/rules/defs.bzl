"""Public declarations for OpenSpec workspaces and their validation."""

load(":openspec.bzl", _openspec_validate_test = "openspec_validate_test")
load(":workspace.bzl", _openspec_workspace = "openspec_workspace")

openspec_validate_test = _openspec_validate_test
openspec_workspace = _openspec_workspace

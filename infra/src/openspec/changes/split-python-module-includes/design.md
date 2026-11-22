## Context

The root module generator discovers `include.MODULE.bazel` files in the
monorepo. InfiniTime already uses a component-local pip declaration.

## Goals / Non-Goals

Keep each component hub beside its lockfile and preserve the existing hub
names, lockfiles, and Python runtime.

## Decisions

Use the existing rules_python pip extension from each component fragment.
Set the shared pip default from `tools/py/python-version.toml` and let the
moved hubs and repository-wide hub inherit it. Keep the shared toolchains
and the repository-wide lockfile's hub in `tools/py`.

## Risks / Trade-offs

Update the shared pip version and the explicit default toolchain registration
together when changing the repository Python version. A component using a
different toolchain can set its own pip version.

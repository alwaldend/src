---
title: Terraform launcher
description: Execute Terraform with providers declared by Bazel
---

This implementation launcher reads the generated target sidecar and resolves
Terraform, configuration inputs, and provider archives through Bazel runfiles.
Consumer targets use the public Terraform rules instead of invoking this binary
directly.

The launcher selects the target's working directory inside the runfiles tree.
Source files can remain symlinks, so formatting retains its existing source
behavior. Initialization writes lockfiles and Terraform metadata in that runfiles
directory. Saved plans and local state retain their normal paths relative to
the working directory; the launcher never removes them.

Each invocation uses a private CLI configuration containing only the declared
filesystem provider mirror. It overrides user CLI configuration, removes inherited
provider-cache settings, disables checkpoint checks, and rejects injected CLI
arguments, provider reattachment, and provider-download commands. Missing provider
versions and platforms fail without a registry installation fallback. Temporary
CLI configuration and any materialized provider mirror are removed on exit.

Before commands that can load providers, the launcher reads Terraform's lock
selections through `version -json` and checks every selected installed package
against its declared archive. Source addresses, versions, file paths, and file
contents must match; undeclared selections and extra, missing, modified, or linked
package files fail before the requested command starts. The check honors
`TF_DATA_DIR`, including paths relative to Terraform's working directory. It also
runs after explicit initialization. This adds local archive decompression and
hashing work; it does not initialize the backend or contact a registry. The
working directory must not be modified concurrently during execution.

Initialization runs before the requested command unless `--direct` precedes it.
Its standard output goes to standard error, preserving structured command output.
`AL_TF_BACKEND_CONFIG*` environment values remain literal initialization
arguments; other provider and backend environment, including `TF_HTTP_*`, is
preserved. Terraform process exit codes are returned unchanged.

Configuration supports declared local modules. Provider installation is isolated;
the launcher does not inspect module source expressions or sandbox Terraform's
service connections.

`--require-saved-plan` accepts only `apply <saved-plan-file>`. The file must exist
and be regular; relative paths resolve against the selected working directory.
Argument and environment validation happens before initialization.

Manifest-only tests materialize declared configuration inputs beneath
`TEST_TMPDIR`, retaining their workspace until Bazel removes that test directory.
Outside tests, a directory runfiles tree is required to preserve persistent
working-directory semantics. An explicit `--chdir` before the Terraform command
selects a caller-owned directory when that behavior is intended.

Invocation scratch uses `TEST_TMPDIR`, then an explicitly configured `TMPDIR`,
then `BUILD_WORKSPACE_DIRECTORY/out/rules_terraform/runtime`. It never falls back
to the system temporary directory.

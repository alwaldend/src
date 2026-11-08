---
title: Command and plugin lifecycle
description: Startup, shutdown, and secret cleanup guarantees
---

Build merges AL configs and propagates plugin runfiles into the command's
runfiles. The Runfiles subgraph contains the merged AL configuration, plugin
runfiles and command executable; source AL configs remain build inputs. AL reads the merged configuration,
starts plugin subprocesses and communicates with them using gRPC over stdio.
Plugins return environment values, file paths and data through that connection,
as shown by the dashed Resources edge. AL runs the command with its
prepared environment. Plugins connect to external services. Build arrows show configuration
dependencies; Run arrows show AL's interactions with the built outputs and plugins.

![Build and Run phases of al run](/docs/projects/al/docs/al/assets/run-lifecycle.svg)

The [Mermaid source](/docs/projects/al/docs/al/run.mmd) is authoritative. This repository renders and
packages its light/dark assets through `//tools/al:run_lifecycle`.

`al run` starts selected plugins, passes their environment to a command, and
keeps plugin resources available until that command exits. On cancellation it
asks the command to terminate and waits up to ten seconds before killing it.
Plugin shutdown then drains requests, releases resources in reverse registration
order, and waits for plugin exit. Startup failure cancels sibling starts and
rolls back partially initialized resources. Cleanup errors fail the invocation.

A standalone installation can run ordinary commands and explicitly selected
absolute plugin executables without Bazel runfiles. Selected relative plugin
identifiers still require runfiles. Available declared or executable-adjacent
runfiles are passed to children; an invalid declared runtime context fails
instead of silently falling back to standalone execution.

Selected plugins may export the same environment name only when their values
agree. Different values fail before the child starts and trigger plugin cleanup.
The diagnostic names the key and plugins in a deterministic order without
printing the values. Explicit injection exports override plugin startup defaults,
including otherwise conflicting startup values. Agreed plugin exports can
override inherited environment values intentionally.

`--stdout` and `--stderr` create output files when absent and truncate existing
files before execution. New files use owner-only permissions; existing file
permissions and ordinary symlink-following behavior are retained. `--stdin`
opens an existing readable file. AL closes its own file descriptors after the
child and plugins stop, while ordinary inherited streams remain borrowed.
Output preparation is not transactional: later startup failure does not restore
previous output contents.

The injector creates secret files with owner-only permissions and removes them
after consumers stop. Vault stores revoke AppRole tokens they issue; failed
AppRole login attempts destroy their unused single-use SecretID. The role's
bootstrap policy must permit `secret-id-accessor/destroy`, as declared by
`projects/tf_modules/vault_approle`. The user's existing token-helper credential
is never erased or revoked. Secret inputs and remote response bodies are omitted
from diagnostic messages. OIDC authorization uses the configured Vault HTTPS
origin and TLS transport, rejecting alternate origins and redirects.
Explicit config dumps still contain the requested configuration; file outputs
use owner-only permissions and truncate previous contents.

`no_auth` prevents AL's Vault client from loading credentials. For environment
injection it clears inherited `VAULT_TOKEN`; it does not sandbox the command or
prevent a command from independently reading the user's token helper. Plugins
and invoked commands remain trusted programs with the user's filesystem access.

Cleanup is observable best effort, not secure erasure. Memory copies are not
zeroized. Forced termination, host failure, unresponsive cleanup code, or an
unavailable service can prevent deletion or revocation. Subprocess signaling
targets direct children, not arbitrary descendants. Credential expiry remains
a fallback; service-specific logout requirements are documented with each plugin.

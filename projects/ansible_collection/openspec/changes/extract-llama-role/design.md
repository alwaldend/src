## Decision

Preserve the separate inference account, loopback binding, CUDA checks, model
digest verification, and readiness after notified restarts. Parameterize
archive paths, runtime/model paths, account, context, port, and model alias.
The caller owns driver provisioning and image constraints.

Bazel packages large archives only for the deployment that consumes them;
adding the role to the shared collection must not download CUDA everywhere.
Ansible unarchive compares contents on each run to repair incomplete extraction.
The role's handler is registered before installation tasks notify it.

## Failure cases and validation

Check missing archives, wrong root stripping, incomplete library extraction,
wrong model digest, lost runtime environment, and health before pending restart.
Verify packaged files and rendered units, and exercise archive repair and
unchanged reapplication with disposable controller-side inputs. GPU loading
and live systemd activation require separately authorized deployment.

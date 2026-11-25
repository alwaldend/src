---
title: MLflow
description: MLflow tracking service and Vast.ai job coordinator
---

`alwaldend.main.mlflow` installs a caller-supplied runtime archive and enables
the MLflow tracking service and Vast.ai coordinator on a Linux systemd host.
Run it with privilege escalation. The services listen on loopback and keep
their data in `/var/lib/mlflow-tracking` and `/var/lib/mlflow-vastai`.

Supply `mlflow_runtime_archive` as a controller-side path to a tar archive with
`mlflow`, its runfiles, and `mlflow-vastai` at the archive root. The role
extracts each archive revision into its own directory and activates it at
`/opt/mlflow/current`.

Set `mlflow_tracking_uri` to the authenticated URL used by clients and workers.
The runtime archive must include the MLflow OIDC plugin and a coordinator
that supports a bearer-token file. Supply the OIDC provider registry,
permitted groups, client secret file, session signing key file, and coordinator token
file. The coordinator also needs its Vast.ai API key, allowed rentals, SSH
identity, and verified host keys. The caller configures the HTTPS reverse
proxy and prepares the workers.

When the caller provides dedicated volumes, set `mlflow_require_data_mounts`
to `true` and mount them before invoking the role. Tracking then requires
mounts at `/var/lib/mlflow-tracking` and `/var/lib/mlflow-tracking/artifacts`.
The coordinator requires a mount at `/var/lib/mlflow-vastai`.

Set `mlflow_enable_coordinator` to `false` to install tracking independently.
Provide only the browser OIDC registry entries in this mode. The role skips
coordinator setup and leaves any previously installed coordinator untouched.

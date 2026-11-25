---
title: MLflow
date: 2026-11-25
description: MLflow tracking and Vast.ai project execution
---

MLflow stores experiment metadata in SQLite and serves run artifacts over HTTPS
at `https://mlflow.alwaldend.com`. MLflow authenticates browser users through Vault OIDC.
The tracking API accepts bearer tokens for workload accounts.
External clients also need a client certificate trusted by shared ingress.
The Xen Orchestra guest keeps service data on four dedicated persistent disks.

| Service data             | Xen device  |   Size | Mount point                          |
| ------------------------ | ----------- | -----: | ------------------------------------ |
| MLflow tracking database | `/dev/xvdb` | 10 GiB | `/var/lib/mlflow-tracking`           |
| MLflow artifacts         | `/dev/xvdc` | 40 GiB | `/var/lib/mlflow-tracking/artifacts` |
| Vast.ai coordinator      | `/dev/xvdd` |  5 GiB | `/var/lib/mlflow-vastai`             |
| Traefik state            | `/dev/xvde` |  2 GiB | `/opt/traefik`                       |

[The storage declaration](storage.json) supplies the disk order, sizes, and
mount points. Back up each volume. Ansible checks the attached disks before
formatting and refuses to mount a new volume over existing data. The services
require their data mounts at startup.

Bazel packages the pinned MLflow CLI, Python runtime, Vast.ai Projects plugin,
OIDC plugin, and coordinator. Ansible installs the package at `/opt/mlflow/current`. The
tracking process and the coordinator listen on loopback. Use the deployed CLI
from a trusted shell on the guest to submit projects. The coordinator assigns
jobs to existing Vast.ai rentals and contacts their workers over SSH.

## Credentials

Vault's MLflow OIDC provider admits members of `global_admins`. The
`tf_setup` root generates the session-signing key and stores it in Vault. The
tracking volume stores plugin accounts and sessions in `auth.db`.

Automation uses a token from an OIDC issuer that signs workload JWTs. Initialize
`alwaldend.com/vault1/approles/src_infra_mlflow/workload` using the component
AppRole, as in the coordinator example below. The JSON fields are `issuer`,
`discovery_url`, `audience`, and `token`. The discovery endpoint must advertise
a reachable JWKS URL. Tokens must carry the pinned issuer and audience and
an `email` or `preferred_username` claim.

In the MLflow UI, create external service accounts bound to the `workloads`
provider and each token's subject. Grant each account access to its intended
experiments. The coordinator rereads
`/var/lib/mlflow-vastai/tracking-token` on every tracking request. Replace the
file atomically with a fresh token before expiration, or update the Vault
entry and rerun Ansible. Credential renewal belongs to the workload issuer
and operator.

Set `MLFLOW_TRACKING_TOKEN` for the CLI and rental worker processes. Refresh
the credential before it expires, including inside long-running MLflow
clients. External clients also need `MLFLOW_TRACKING_CLIENT_CERT_PATH` for
shared ingress. Direct guest clients can set
`MLFLOW_TRACKING_SERVER_CERT_PATH` to the trusted Vault CA bundle.

Initialize the coordinator KV entry from a task-private JSON file containing
these fields.

| Field             | Value                                                                     |
| ----------------- | ------------------------------------------------------------------------- |
| `api_key`         | Vast.ai key with permission to inspect the configured existing rentals    |
| `workers`         | Comma-separated exact rental labels or IDs in preference order            |
| `ssh_private_key` | Dedicated unencrypted OpenSSH Ed25519 identity, with its trailing newline |
| `ssh_known_hosts` | Verified OpenSSH host keys for the endpoint and port of every rental      |

Authorize that key on the worker account, which defaults to `root` in
[the inventory](ansible/group_vars/mlflow.yaml). Use bracketed
`[hostname]:port` known-hosts entries for non-default SSH ports. Write the
initial JSON object using the component AppRole and prevent replacement of
an existing entry with `-cas=0`.

```sh
bazel run //infra/mlflow:vault.kv_put -- -cas=0 \
  alwaldend.com/vault1/approles/src_infra_mlflow/coordinator \
  - < "$MLFLOW_COORDINATOR_JSON"
```

## Deployment

The [AppRole](approle/README.md), shared Vault memberships, and [DNS
records](dns/dnsconfig.json) are deployment prerequisites.

1. Apply `//infra/mlflow/approle:tf.apply` to register the Vault OIDC client.
2. Synchronize the AppRole in Xen Orchestra using the
   [resource-set workflow](../xcp_ng/tf/README.md). Apply the shared resource
   set before running the guest setup.
3. Apply `//infra/mlflow/tf_setup:tf_setup.apply` to create the VM and the
   OIDC session key.
4. Initialize the coordinator and workload KV entries.
5. Run `bazel run //infra/mlflow/ansible`, then deploy the HTTPS route with
   `bazel run //infra/ingress/ansible`.
6. Sign in with Vault and register the workload accounts before submitting runs.

## Run a project

Prepare each existing rental with the worker binary, MLflow CLI, tracking
credentials, and project runtime. Ensure the worker service and the SSH user
use the same state directory and both can find the packaged binaries. Follow
the [worker setup](../../projects/mlflow_vastai_backend/README.md#rental-worker).
The packaged Linux runtime is available from
`bazel build //infra/mlflow/runtime`.

Create a backend file on the MLflow guest using the local coordinator.

```json
{
  "server": "127.0.0.1:50051",
  "insecure": true
}
```

After loading a workload bearer token into the submitting shell, run:

```sh
export MLFLOW_TRACKING_URI=https://mlflow.alwaldend.com
/opt/mlflow/current/mlflow run /path/to/project \
  --backend vastai --backend-config /path/to/backend.json --env-manager local
```

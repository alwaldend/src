## Offline verification

- The packaged Ansible playbook passes syntax checking. The MLflow and shared
  ingress Ansible packages and the MLflow documentation target build.
- The MLflow Terraform format and validation tests, OpenSpec validation,
  repository DNS inventory test, and plugin, worker, and coordinator integration
  tests pass. The coordinator fixture requires a rotating bearer token for tracking
  updates and checks that tracking and discovery credentials stay out of worker SSH
  subprocesses.
- The unpacked Linux runtime starts its own MLflow tracking server. A CLI
  invocation resolves the Vast.ai project entry point, and a run created
  through the tracking API and its uploaded artifact survive a server restart.

Live Vault, Xen Orchestra, DNS, ingress, and Vast.ai operation require the
deployment and worker setup in the project README.

The unpacked OIDC runtime was exercised against a loopback discovery fixture.
The server started, rejected anonymous and invalid bearer API requests, and
redirected browser sign-in to the configured issuer with the expected callback
and scopes. Live Vault sign-in and workload account provisioning are deployment
acceptance steps.

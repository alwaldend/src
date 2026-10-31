---
name: repo-secrets
description: Safely add, route, rotate, or review secrets and authentication in this repository using Vault, AppRole, the al.lua DSL, injectors, and Bazel wrappers. Use whenever work involves credentials, tokens, private keys, sensitive Terraform values, Ansible secrets, Vault policy, or secret-bearing configuration.
---

# Handle Secrets

## Follow the repository trust boundary

- Store secret values in Vault, never in Git, Terraform source, Ansible vars,
  command history, test fixtures, ordinary logs, screenshots, or PR text. An
  unavoidable secret-bearing raw artifact follows the task-private temporary
  handling below and is not durable evidence.
- Commit only secret **references**: Vault mount/path, environment variable name,
  plugin label, policy, role, or non-sensitive metadata.
- Treat checked-in `infra/` source as public repository content while
  preserving its repository-internal target visibility, artifact-publication,
  and production-build dependency boundaries. Decrypted, generated, or live
  artifacts are not a separate confidentiality class, but inspect them before
  disclosure because they can contain secrets or personal information.
- Treat an unexpected credential already in the tree as compromised. Do not
  repeat it in output. Stop exposing it, preserve evidence minimally, and tell
  the user it needs revocation/rotation.

## Reuse the injection flow

1. Find the nearest component `al.lua` and its parent `al_config` dependency.
2. Reuse `lib.vault_auth` for AppRole authentication.
3. Use `lib.plugin_call` with the existing component plugin (`injector`,
   `tf_backend`, provider login, or another established plugin).
4. Select calls through labels such as `tf=main`, `tf=setup`, or `ansible=1`.
5. Ensure the corresponding plugin binary is present in the Bazel target's
   `data` and its label is present in `run_args`.
6. Inject the minimum values and grant the minimum Vault policy needed.

Follow sibling `al.lua` files exactly for the DSL shape. Keep AppRole names and
paths aligned with the component label; do not silently share a broader role.

## Prevent disclosure while working

- Do not print secrets or personal information. Do not print whole environment
  dumps, rendered templates, Terraform state, Vault responses, token files, or
  generated credentials without review; select or structurally summarize only
  fields proven not to contain prohibited content.
- Use placeholders in examples. Prefer stdin or `@file` input over a literal
  value on a command line. Keep temporary secret files outside tracked or
  committable source. When a file is unavoidable, use ignored task-private
  `out/<task>/` storage with restrictive permissions, minimum lifetime, and
  explicit cleanup; never stage it or promote it as evidence.
- Do not run Vault writes, token generation, unseal, apply, or deployment merely
  to test a code change.
- Review diffs with `git diff --check` and a targeted search for newly introduced
  private keys, literal passwords/tokens, and suspicious high-entropy values.
  Do not use a scan that echoes suspected values into the final response.

## Interactive Bash credential input

When a user requests a pasteable command that prompts for a token, read it
silently at runtime instead of putting the value in the command. Use
`IFS= read -r -s -p` to preserve spaces and backslashes and suppress terminal
echo. Print a newline after the hidden prompt. Stop on failed input or an empty
value before invoking the consumer.

Use a subshell to scope the temporary variable and preserve the consumer's exit
status. Disable shell tracing before reading or expanding the secret, and unset
the variable first so an inherited export attribute cannot expose it to child
processes. Prefer the consumer's standard-input credential option; quote the
expansion and use a fixed `printf` format. For example, replace the placeholder
consumer with the actual tool and its documented stdin option:

```bash
(
  set +x
  unset secret_input
  IFS= read -r -s -p 'API token: ' secret_input || exit 1
  printf '\n'
  [ -n "$secret_input" ] || exit 1
  printf '%s\n' "$secret_input" | consumer --token-stdin
)
```

The literal token stays out of shell history and process arguments. Do not
substitute an argument such as `--token "$secret_input"`, export the variable,
or echo the value for confirmation. Subshell exit removes the shell variable;
it does not guarantee memory erasure or control how the consumer stores login
credentials. If stdin authentication is unavailable, use the tool's documented
secure input mechanism and describe any necessary storage. Check prompt,
input-failure, consumer-failure and cleanup behavior with dummy values; never
use a real credential as a fixture or run a login merely to test the example.

## Validate structure

Build or test the narrow Bazel targets that package the changed `al.lua`, Vault
policy, Terraform, or Ansible configuration. Explain when live authentication is
required and leave the live operation to an authorized operator unless the user
explicitly requests it and the environment is clearly non-production.

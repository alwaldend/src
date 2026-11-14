## Decision

Proceed with a shared role receiving model entries, immutable Git node pins,
controller-side first-party modules/workflows/requirements and runtime paths.
The playbook resolves these from its JSON registry. Large assets
remain deployment dependencies rather than dependencies of every role consumer.
Terraform retains constrained rental creation and access bootstrap; Ansible
owns ComfyUI's approved-node arguments as well as installation. On reboot the
bootstrap preserves the Ansible-managed argument line in `/etc/environment`.

Keeping the updater as a deployment-specific task list duplicates installation
behavior. Keeping large application files in Terraform also couples ordinary
asset changes to replacement of a paid rental. Neither meets the requested
ownership boundary. The existing Vast image owns ComfyUI and its Python venv;
the role configures that installation rather than building a second one.

## Failure cases and validation

Before implementation, exercise installation under a non-default root, correct
checkpoint/LoRA/ControlNet destinations, pinned Git checkout, malformed paths
and missing pins, checksum mismatch, first-party module and browser files,
workflow JSON preservation, dependency integrity flags, changed-node restart,
unchanged reapplication, workflow-copy opt-out, and busy queue refusal before
mutation. Read only the aggregate `/prompt` count and readiness endpoints.
Emit sanitized local installation evidence; use no GPU or remote host.

## Verification

The packaged role fixture installs checkpoints, LoRAs, ControlNet, an exact
local Git revision, first-party Python/browser files and a native App below a
non-default root. It verifies no additional restart on reapplication, a changed
node restart with workflow copies skipped, and refusal of busy, wrong-checksum
and path-escape updates. `installation.json` is the repeatable sanitized artifact.
The deployment fixture additionally installs all 75 generated workflows across
11 profiles through the real playbook and runs the full selected NoobAI profile
with local model/Git substitutes and real packaged first-party source.
Its `profile-workflows.json` receipt records installed counts. Real dependency
installation and Vulkan/GPU execution are not covered by these fixtures.

The generic fixture initially exposed YAML quoting and whitelist sorting
issues; both were corrected. The deployment fixture required the pinned
controller interpreter for local modules, rather than the host's incompatible
interpreter. Failure diagnostics now omit embedded module archives. This is a
fixture correction; no host Python or configuration changed.

## Plain role variables

The playbook translates its JSON registry into ordinary role variables once.
The role accepts models, node sources, workflow files, model directories, runtime
arguments and token inputs; it does not know the deployment's registry or
profile schema. Remove the conditional dual interface and profile facts.
The caller owns selection and path mapping; installation, validation, idle
protection and restart handling remain in the role. AL owns Vault access.

Use the already-pinned community.general dict_kv filter to pair registry model
entries with their filenames without mutable accumulation. Other mappings use
standard Ansible filters. No new dependency or live operation is required.

Failure cases: selecting a non-default profile must install its own workflow
catalog; selected model filenames and source paths must survive translation;
explicit role overrides must keep precedence; workflow-only runs require no
model token; generic callers must work without registry variables. Existing
local deployment and role fixtures cover those boundaries.

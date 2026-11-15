## Why

The Forgejo runner and Host Bot have separate infrastructure owners and deployment wiring. Consolidate their desired state under infra/dev, naming the secure runner secure1.dev.alwaldend.com and Host Bot dev1.dev.alwaldend.com.

## What Changes

- **BREAKING**: Move both owners, maintained OpenSpec work, and consumers to infra/dev.
- Provide a common packaged Ansible entry point with host-specific configuration.
- Rename host and application DNS names, ingress routes, and certificate scope.
- Merge DNS into one root using the shared src_infra_dev identity and one backend; document the separately authorized transfer from both legacy states.
- Remove the obsolete Proxmox runner declaration, unused runner credential injection and old inventory names while preserving host-specific behavior.

## Capabilities

### New Capabilities

- `development-deployment`: Unified owner, deployment, host identities, and safe state continuity.

### Modified Capabilities

Existing component contracts move with their owner.

## Impact

infra/forgejo_runner, users/simeonwarren/host_bot, infra/vault, infra/ingress, security inventory, build consumers, and documentation. No new dependencies or live deployment. Acceptance requires offline packaging, source checks, inventory and route rendering, and explicit documentation of pending operational cutover.

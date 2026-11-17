# Validation

The VM root validates with only Vault/XO providers. Its configured dependency
graph has no TrueNAS provider. The development AL configuration contains no NAS
authentication or API injector; its VM backend, disk and metadata are unchanged.

The namespace/cache project at infra/nas/approles/src_infra_dev validates as an
independent root with its own state and targets. Generated configuration selects
its distinct backend using NAS operator authentication. Its package contains the
namespace and cache, while shared NAS launchers exclude all AppRole Terraform
files. Both NAS roots' plan/apply launchers were regenerated and inspected.

The stage order is Vault groups/grants, shared NAS accounts/settings, independent
development NAS storage, XCP-ng, VM provisioning and Ansible. Each apply covers
its own root. Existing packaged Ansible syntax and seven metadata cases remain
valid for unchanged playbook inputs. The generic dataset implementation and its
synthetic ACL rendering are unchanged. No live state move, plan, apply, credential
operation or guest mount ran. Ownership reconciliation is required before apply
if an older shared-state storage revision was deployed elsewhere.

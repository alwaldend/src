## 1. Implementation

- [x] 1.1 Grant shared DNS access and remove the unused runner-config read; verify formatted Terraform and packaged configuration without live changes.
- [x] 1.2 Under the user's deployment request, apply the scoped shared identity and verify its bootstrap access; preserve legacy identities until state reconciliation establishes whether they remain needed.
- [x] 1.3 Retire both legacy identities and their remaining grants after confirming both old backends are empty; inspect the scoped interactive cleanup apply.
- [x] 1.4 Remove the retired OpenHands read policy and development AppRole attachment through the owning interactive apply.

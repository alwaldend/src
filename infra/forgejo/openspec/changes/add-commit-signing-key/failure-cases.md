# Failure cases before implementation

- Vault permission, network, TLS, redirect, or malformed responses must not be
  interpreted as a missing key. Soft deletion and destroyed history must not
  permit regeneration. CAS conflicts must not overwrite another pair.
- Missing fields, mismatched public/private keys, wrong algorithms, incorrect
  fingerprints must stop cache reuse.
- Terraform planning and ordinary repeat applies must not generate keys;
  explicit replacement must reuse Vault. Terraform state/output must omit keys.
- Missing injected values and a different installed private key must stop before
  service activation. Missing public content must be repairable, and private
  permissions must be repaired on repeats without rotating the identity.
- Forgejo must expose the expected public signing key and generate a genuinely
  verifiable SSH signature through its file-edit API. A repeat role run must
  preserve both identity and service lifecycle. Retain sanitized test receipts.

## Import-based reuse

The revised Terraform CLI cases first failed against the counted-resource
implementation: the singleton import target did not exist and the legacy
upgrade changed its address. The module now retains its singleton address and
has no existing-team selector or moved block. The full loopback suite passes
for managed bootstrap and upgrade, imported-team key issuance through an aliased
provider, invalid imports, scoped-permission rejection, failed activation with
tracked-key recovery, and key-only cleanup with refused team deletion. Fixture
receipts are repeatable Bazel undeclared outputs. Temporary key-move fixture
blocks are removed after apply and followed by a no-change plan.

The owning live bootstrap plan also returned no changes after the completed
root migration block was removed. No team was imported again or recreated.
Live alwaldend key binding remains unverified after the earlier HTTP 400.

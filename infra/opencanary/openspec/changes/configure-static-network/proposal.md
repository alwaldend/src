## Why

Give OpenCanary a stable address and a site-local name on the regular Xen
Orchestra network.

## What Changes

Declare the static IPv4 address in the component DNS document. Use it for
guest configuration and deploy the record from a dedicated DNS root. Install
a Vault-signed certificate on the decoy HTTPS listener for ingress traffic.

## Capabilities

### New Capabilities

- `static-guest-network`: Static addressing, site-local DNS, and private HTTPS access.

### Modified Capabilities

None.

## Impact

The Vault AppRole needs the site-local DNS policy and Xen Orchestra needs a
resource-set network assignment.

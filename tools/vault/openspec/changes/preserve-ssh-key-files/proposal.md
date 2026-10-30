## Why

Publishing the requested schema exposed a signed SSH certificate overwriting the generated public-key file and triggering an OpenSSH key-permission warning.

## What Changes

- Preserve the generated public key and write the signed certificate to its standard -cert.pub sibling.

## Capabilities

No specification delta; this corrects temporary-file ownership within the existing signing contract.

## Impact

The shared SSH injector; both authorized download-host deployment and publication exercise the corrected paths.

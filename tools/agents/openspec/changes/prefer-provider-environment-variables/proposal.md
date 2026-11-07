## Why

Provider-supported environment variables avoid redundant Terraform credential
variables. The user requested this preference in the infrastructure skill after
reviewing the printers configuration.

## What Changes

Prefer documented provider-native environment variables for credentials and
endpoints, reusing the repository injection flow. Keep explicit provider inputs
when aliases need distinct identities or the provider lacks environment support.
Add evaluation scenarios for the preference and its alias exception.

## Impact

Only canonical repo-infra guidance and its offline evaluation fixtures change.
The related implementation is maintained by the active printers and MikroTik
changes. This is a procedural refinement without a product specification delta.

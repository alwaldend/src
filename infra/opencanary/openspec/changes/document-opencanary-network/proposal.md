## Why

Show the network paths to the OpenCanary VM in the project documentation.

## What Changes

Add a Mermaid diagram showing public mTLS ingress and one local edge that
lists each enabled decoy on its own line. Place the local actor and the
OpenCanary VM inside the Local infra container, and label the VM with its
fully qualified name. Publish the diagram through the existing documentation
build.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None.

## Impact

The OpenCanary README displays the network diagram in the site's light and
dark themes.

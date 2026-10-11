## Why

Keep the 3x-ui architecture beside its owner and make node placement and internal
service relationships readable without repeating the same detail in every host.

## What Changes

- Move the Threexui Drawio page to a canonical component-owned Mermaid source.
- Publish site SVG and WebP renders at the top of the component README.
- Show njal.la and Yandex Cloud as separate datacenters with simple node cells.
- Show Traefik, 3x-ui, Xray, and Mullvad relationships in one shared detail subgraph.
- Remove the superseded Drawio page, render mapping, and maintained SVG.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `infra-threexui`: Own the architecture source and describe node placement and internals.

## Impact

Documentation and rendering only. Node configuration, DNS, credentials, and live
deployment are outside this change. No new dependencies are required.

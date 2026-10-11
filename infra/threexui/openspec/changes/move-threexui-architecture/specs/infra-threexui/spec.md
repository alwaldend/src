## ADDED Requirements

### Requirement: Component-owned architecture diagram

The component README SHALL begin with a Mermaid architecture diagram showing
the control plane, data-plane node placement, ACME, and Mullvad traffic paths.
The diagram SHALL show separate njal.la and Yandex Cloud datacenters with simple
node cells for njalla1 and work1, and yc1 respectively. A shared node-detail
subgraph SHALL show Traefik routing to 3x-ui and Xray, 3x-ui managing Xray, and
Xray connecting to Mullvad. Site SVG and WebP renders SHALL share one source.
The superseded Threexui Drawio page SHALL no longer be maintained in infra/arch.

#### Scenario: Consult the 3x-ui architecture

- **WHEN** a reader opens the component README
- **THEN** its first section shows node placement in both datacenters
- **AND** a separate subgraph explains the services shared by every node.

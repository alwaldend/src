## MODIFIED Requirements

### Requirement: Declare explicit node and outbound routing relationships

Terraform SHALL configure the `njalla1` and `yc1` nodes through the shared node
module. The node routing declaration SHALL block the advertising category and
private IP destinations before routing the freedom inbound directly and each
`mullvad_min` or `http_proxy` inbound through its named WireGuard outbound.
The `mullvad_min_lb` inbound SHALL route through the node's Mullvad balancer.
The observatory SHALL probe the balancer's Mullvad WireGuard outbounds, and
the balancer SHALL use the `blocked` fallback when it has no healthy candidate.

#### Scenario: Inspect a configured Mullvad inbound route

- **WHEN** a node's selected relay creates a Mullvad inbound
- **THEN** the inbound listens on loopback
- **AND** its routing rule selects the matching `out-mullvad-min-` outbound tag
- **AND** the node declaration supplies its own relay selection to the module

#### Scenario: Route an eligible load-balanced connection

- **WHEN** a connection reaches the Mullvad load-balancing inbound and passes
  the blocking rules
- **THEN** Xray uses the node's Mullvad balancer for that connection
- **AND** the balancer selects from healthy observed Mullvad WireGuard outbounds
- **AND** the balancer selects `blocked` when no healthy candidate is available

## MODIFIED Requirements

### Requirement: Declare explicit node and outbound routing relationships

Terraform SHALL configure the `njalla1`, `yc1`, and `work1` nodes through the shared node
module. The node routing declaration SHALL block the advertising category and
private IP destinations before routing the freedom inbound directly and each
`mullvad_min` or `http_proxy` inbound through its named WireGuard outbound.
The enabled `mullvad_min_lb` inbound has no active dedicated balancing rule;
commented balancing rules SHALL NOT be treated as enabled routing behavior.

#### Scenario: Inspect a configured Mullvad inbound route

- **WHEN** a node's selected relay creates a Mullvad inbound
- **THEN** the inbound listens on loopback
- **AND** its routing rule selects the matching `out-mullvad-min-` outbound tag
- **AND** the node declaration supplies its own relay selection to the module

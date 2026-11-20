## ADDED Requirements

### Requirement: VictoriaMetrics service role

The VictoriaMetrics role SHALL install the caller-supplied executable,
render and validate the scrape YAML template, and manage the single-node
systemd service.

#### Scenario: Configure a guest

- **WHEN** a play invokes the role with a compatible executable
- **THEN** the role installs the executable and validates the scrape configuration before enabling the loopback service.

#### Scenario: Configure a predefined VMUI dashboard

- **WHEN** a caller supplies a VMUI dashboard JSON file
- **THEN** the role validates and installs it and configures the service to serve it through VMUI.

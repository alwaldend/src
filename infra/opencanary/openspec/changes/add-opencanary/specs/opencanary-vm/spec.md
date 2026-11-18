## ADDED Requirements

### Requirement: Dedicated native OpenCanary guest

The OpenCanary infrastructure SHALL declare one XCP-ng guest on an operator-selected DHCP network and deploy OpenCanary through an Ansible role as a native systemd service.

#### Scenario: Provision and configure the dedicated guest

- **WHEN** an operator provisions and configures the guest with the owning Terraform and Ansible commands
- **THEN** the guest runs the configured OpenCanary decoys without a container runtime
- **AND** decoy interaction events are available in the local system journal
- **AND** administrative SSH and the SSH decoy listen on separate ports
- **AND** Xen Orchestra reports the DHCP-assigned IPv4 address for Ansible

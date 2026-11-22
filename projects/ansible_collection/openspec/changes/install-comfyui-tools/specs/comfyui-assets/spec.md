## ADDED Requirements

### Requirement: Selected asset installation

The role SHALL install caller-selected model kinds, pinned node packs,
first-party sources and packaged node bundles, locked requirements, and native workflow JSON beneath
caller-configured application paths without selecting a GPU or provisioning a host.

#### Scenario: Install and reapply a bundle

- **WHEN** the same valid bundle is applied twice to an idle installation
- **THEN** files SHALL retain their intended contents and destinations
- **AND** unchanged node inputs SHALL NOT restart the service again.

### Requirement: Preserve active generation and edited workflows

The role SHALL wait for an aggregate idle count before mutation and again
before changing loaded application code, without inspecting job contents.
Workflow copies SHALL support the `workflows` tag.

#### Scenario: Update nodes while retaining workflow edits

- **WHEN** the caller skips the workflows tag
- **THEN** existing workflows SHALL remain untouched while changed nodes are installed.

### Requirement: Ordinary role variables

The role SHALL accept explicit model entries, node sources, workflow paths,
native model directories, runtime arguments and token inputs. Its defaults
SHALL NOT interpret a deployment registry or create temporary profile facts.
The deployment playbook SHALL translate its JSON into those variables once.
The role SHALL consume the supplied model token without contacting Vault.

#### Scenario: Translate a selected deployment profile

- **WHEN** the playbook selects a profile from its JSON registry
- **THEN** its role variables SHALL identify that profile's models, sources and workflows
- **AND** workflow-only runs SHALL require no model token or Vault access

#### Scenario: Supply role variables directly

- **WHEN** another caller supplies the ordinary role variables without a registry
- **THEN** the role SHALL validate and install those assets normally

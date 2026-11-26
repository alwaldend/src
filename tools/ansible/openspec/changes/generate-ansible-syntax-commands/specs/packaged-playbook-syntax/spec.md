## Purpose

Make it possible to check the syntax of packaged Ansible playbooks using
the same runtime inputs and arguments as their playbook commands.

## ADDED Requirements

### Requirement: Packaged playbook syntax tests

Each packaged Ansible playbook command SHALL have a sibling Bazel test whose
name appends `.syntax`. The test SHALL run the packaged executable
with `--syntax-check` and the command's playbook arguments, including any
packaged extra-vars files, independently of deployment plugins.

#### Scenario: Check a default command

- **WHEN** an operator selects `<name>.syntax` with `bazel test`
- **THEN** Ansible checks the default playbook using its configured inventory
  and arguments without running playbook tasks or loading deployment plugins

#### Scenario: Check named commands

- **WHEN** a package declares named playbook commands
- **THEN** each `<name>.<suffix>.syntax` target checks the playbook and arguments
  associated with that command

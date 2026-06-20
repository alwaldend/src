---
title: Agents
---

# AGENTS.md

This file contains high-signal, repository-specific guidance for AI agents working in this repository. It focuses on executable truths and non-obvious workflows that would otherwise be difficult to infer.

## Architecture Overview

The repository follows a monorepo structure with the following major directories:

- `infra/`: Infrastructure as Code (IaC) for data centers and services. Contains Terraform configurations and Ansible playbooks.
- `projects/`: Project-specific code and modules, including Terraform modules for various services.
- `tools/`: Bazel toolchain definitions and utilities used across the repository.
- `data/`: Static data files used across the repository, including SSL certificates, PGP keys, and documentation assets.
- `third_party/`: External dependencies and vendored code that are managed locally.

## Key Configuration Files

- `pyproject.toml`: Python dependency and linting configuration (black, mypy, isort, flake8)
- `.editorconfig`: Editor configuration with specific indentation rules
- `.bazelrc`: Bazel configuration with project-specific and CI settings
- `.gitignore`: Comprehensive ignore patterns including Bazel, Terraform, and IDE files

## Bazel Build System

The repository uses Bazel as the primary build system with custom toolchains and configurations.

### Bazel Configuration

The Bazel configuration is split across multiple files:
- `.bazelrc`: Main configuration file
- `tools/bazelrc/root.bazelrc`: Imports presets and project configurations
- `tools/bazelrc/preset.bazelrc`: Standard Bazel best practices
- `tools/bazelrc/project.bazelrc`: Project-specific overrides

## Infrastructure Workflow

### Data Center Structure

The infrastructure is organized by data center (`dc1`) with multiple services:
- `vault/`: HashiCorp Vault setup and configuration
- `pve1/`: Proxmox Virtual Environment cluster
- `consul1/`: Consul service mesh
- `forgejo1/`: Forgejo (Gitea fork) code hosting
- `cl1/`: Kubernetes cluster (k3s)

### Cloud Infrastructure

The repository also manages cloud infrastructure through:
- `yandex_cloud/`: Configuration for Yandex Cloud resources and services

### DNS Configuration

DNS is managed through:
- `dns/`: DNS configuration and zone files using dnscontrol
- `infra/dc1/pve1/ansible/roles/dns/`: Ansible role for DNS server configuration

## Authentication and Secrets

The repository uses HashiCorp Vault for secrets management with AppRole authentication. Configuration is managed through Lua files (`al.lua`) in each infrastructure directory.

## Lua Configuration

Infrastructure configuration is defined in Lua files named `al.lua` that use a custom DSL to configure services and authentication. These files are located in:
- `/infra/dc1/vault/al.lua`
- `/infra/dc1/pve1/al.lua`
- `/infra/dc1/consul1/al.lua`
- `/infra/dc1/forgejo1/al.lua`
- `/infra/dns/al.lua`
- `/infra/yandex_cloud/org1/al.lua`
- `/infra/users/simeonwarren/al.lua`

## Terraform Usage

Terraform configurations are managed through Bazel using the `terraform_binary_map` rule. Each infrastructure component has its own Terraform configuration in a `tf/` subdirectory.

### Terraform Directory Structure

- `infra/dc1/vault/tf/`: Vault Terraform configuration
- `infra/dc1/pve1/tf/`: Proxmox Terraform configuration
- `infra/dc1/consul1/tf/`: Consul Terraform configuration
- `infra/dc1/consul1/tf_setup/`: Consul setup Terraform configuration
- `infra/yandex_cloud/org1/tf/`: Yandex Cloud Terraform configuration
- `infra/users/simeonwarren/tf/`: User-specific Terraform configuration

## Development Environment

### Python Environment

Python dependencies are managed through `pyproject.toml` with the following tools:
- black: Code formatting
- mypy: Type checking
- isort: Import sorting
- flake8: Linting

### Editor Configuration

The repository uses `.editorconfig` with the following settings:
- Python, Go, and Makefile: tabs for indentation
- YAML, JSON, Markdown, HTML, proto, TOML: 2 spaces for indentation
- All other files: 4 spaces for indentation
- Maximum line length: 79 characters

## Testing and Verification

### Key Verification Commands

- `bazel test //...` - Run all tests in the repository
- `bazel build //...` - Build all targets
- `black --check .` - Verify Python code formatting
- `mypy .` - Run Python type checking

## Special Considerations

### Git Hooks

The repository configures git hooks through Bazel:
- Pre-commit hook is configured in `BUILD.bazel` using `write_source_file`

### Remote Execution

The Bazel configuration is optimized for remote execution and caching, with settings in `preset.bazelrc` that:
- Enable remote caching and execution
- Set appropriate timeouts for CI
- Configure output downloading strategies

### Platform-Specific Configuration

Bazel automatically enables platform-specific configuration (`--enable_platform_specific_config`), so platform-specific flags can be used in `.bazelrc` files with `:linux`, `:macos`, `:windows` suffixes.

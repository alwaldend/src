# Rules Hugo

## Purpose

Build and run Hugo sites with registered Bazel toolchains, site source archives,
and declared PostCSS tooling. This source baseline was observed on 2026-09-08
at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`.

Sources: [project description](../../../README.md),
[site provider rule](../../../pkg/bzl/al_hugo_site.bzl),
[build rule](../../../pkg/bzl/al_hugo_run_binary.bzl),
[toolchain extension](../../../pkg/bzl/al_hugo_extension.bzl),
and [worker rule](../../../pkg/bzl/al_hugo_worker.bzl).

## Requirements

### Requirement: Separate site inputs from execution tools

`al_hugo_site` SHALL retain the `.tar` site archive in the target configuration
and resolve its executable PostCSS dependency in the execution configuration.
Build rules SHALL obtain Hugo from the registered Hugo toolchain.

#### Scenario: The target and execution platforms differ

- **WHEN** Bazel analyzes a Hugo site build for distinct target and execution platforms
- **THEN** the site archive remains a target input and Hugo and PostCSS are selected as execution tools.

### Requirement: Build into a declared destination directory

`al_hugo_run_binary` SHALL unpack the site archive, make its PostCSS executable
available to Hugo, and invoke Hugo with `--destination` pointing to the declared
output directory.

#### Scenario: A site build specifies an output directory

- **WHEN** a caller supplies `out_dir` and additional Hugo arguments
- **THEN** the action appends the declared destination to those arguments and exposes the output directory through `DefaultInfo`.

### Requirement: Generate platform-specific toolchain repositories

The Hugo module extension SHALL create toolchain repositories from the archive
set associated with the requested version, carrying each archive's declared
integrity and execution-platform constraints.

#### Scenario: A supported version has a Linux x86-64 archive

- **WHEN** a toolchain tag requests that version under a repository name prefix
- **THEN** the extension creates the corresponding `<name>_os_linux_cpu_x86_64` repository with a Hugo toolchain constrained to that platform.

### Requirement: Support an explicit persistent-worker build path

`al_hugo_worker` SHALL serialize site inputs, arguments, tools, environment, and
output location into a flag file and invoke its worker with Bazel's protobuf
worker protocol requirements.

#### Scenario: A caller chooses the worker rule

- **WHEN** a site is built with `al_hugo_worker`
- **THEN** Bazel receives a worker-capable `HugoSite` action whose declared output directory is `<target>.dest`.

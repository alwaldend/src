## Purpose

Provide a reproducible local CPU image-generation trial with inspectable output.

## ADDED Requirements

### Requirement: CPU thread configuration

Both Python entry points SHALL accept an optional positive `--threads` value
without a parser default. When omitted, they SHALL preserve the inherited or
library-default thread settings and SHALL NOT synthesize a child thread flag.
Public Bazel targets SHALL supply `--threads 12` through shared
BUILD arguments and SHALL accept a caller's explicit override. When provided, the benchmark
SHALL forward the selected value to its owned server, which SHALL configure
OpenMP, MKL, and OpenBLAS with that value.

#### Scenario: Run a public target

- **WHEN** the caller starts a public server or benchmark without a thread override
- **THEN** BUILD supplies 12 threads to the runtime

#### Scenario: Override CPU threads

- **WHEN** the caller appends `--threads 4` to the public target command
- **THEN** its server receives four computation threads

#### Scenario: Omitted thread setting

- **WHEN** either Python entry point receives no thread setting
- **THEN** it leaves runtime thread selection unchanged

#### Scenario: Invalid thread setting

- **WHEN** either Python entry point receives a non-positive thread value
- **THEN** argument parsing fails before starting inference

### Requirement: Declared CPU inference inputs

The tool SHALL obtain ComfyUI, Python, dependencies, and model weights through
versioned and integrity-checked Bazel inputs. It SHALL run CPU inference with
custom nodes and external API nodes disabled and no runtime package installation.

#### Scenario: Local generation

- **WHEN** the caller runs the benchmark with a prompt and output directory
- **THEN** it produces a valid 512 by 512 PNG and records the workflow, seed,
  model identity, and elapsed generation time separately from startup.

### Requirement: Bounded benchmark lifecycle

The benchmark SHALL bind to localhost, bound startup and generation waits,
report failed workflows, and stop its own server on success or failure.

#### Scenario: Generation fails

- **WHEN** the workflow fails or the deadline expires
- **THEN** the command fails, preserves diagnostic evidence, and stops its server.

### Requirement: Declared Z-Image-Turbo trial

The Z-Image-Turbo benchmark SHALL obtain its diffusion weights, text encoder,
and VAE through checksum-verified Bazel inputs and use the existing CPU
ComfyUI runtime without runtime downloads or custom node installation.

#### Scenario: Z-Image-Turbo generation

- **WHEN** the caller runs the Z-Image-Turbo benchmark
- **THEN** it produces a valid PNG, preserves the workflow and all model hashes, and records startup and generation times separately.

### Requirement: Bounded Z-Image process memory

The default Z-Image launch and benchmark targets SHALL refuse inference unless
the effective cgroup RAM ceiling is at most 44 GiB and cgroup swap is disabled.
The documented run SHALL use a transient scope and leave persistent host
configuration unchanged.

#### Scenario: Unbounded invocation

- **WHEN** the caller starts a default Z-Image target without those kernel limits
- **THEN** the launcher refuses before starting the inference server.

#### Scenario: Bounded invocation

- **WHEN** the caller starts the generated Bazel command in the documented memory scope
- **THEN** the kernel enforces the RAM ceiling and prevents inference swap use.

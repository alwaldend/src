## ADDED Requirements

### Requirement: Run and check declared Bend sources

The rules SHALL use a registered Bend 2 toolchain to run or check the entry
point with its declared sources.

#### Scenario: Resolve tools for distinct platforms

- **WHEN** a runnable target or check target resolves its toolchain
- **THEN** Bazel selects the runnable compiler for the target platform or the
  checker for the execution platform, respectively

#### Scenario: Run a program with a local import

- **WHEN** a runnable target declares its entry point and imported sources
- **THEN** `bazel run` executes the Bend program with its relative imports
  intact and forwards command-line arguments and the program's exit status

#### Scenario: Check a Bend source

- **WHEN** a check target builds successfully
- **THEN** its action has invoked the compiler in check-only mode and produced
  a success marker

#### Scenario: A source imports an undeclared sibling

- **WHEN** Bend resolves a relative import to a file available only in the checkout
- **THEN** the invocation fails because that file is absent from its staged sources

#### Scenario: Bend rejects a source

- **WHEN** Bend returns a failure
- **THEN** the rule reports the error and the build or execution fails

### Requirement: Provide the default Linux distribution

The rules SHALL register the checksum-pinned Bend 2.0.40 executable and its
standard library for Linux x86-64.

#### Scenario: Use the default toolchain

- **WHEN** a Linux x86-64 build uses a Bend rule with the default registrations
- **THEN** Bazel selects the upstream compiler and makes its standard library
  and runtime files available to it

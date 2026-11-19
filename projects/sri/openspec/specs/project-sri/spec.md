# Subresource Integrity calculator

## Purpose

Calculate file digest output for Subresource Integrity using OpenSSL. This
baseline was observed at repository revision `550d7e79` on 2026-09-08. Sources
are the [project README](../../../README.md),
[command implementation](../../../main/c/cmd.c),
[digest implementation](../../../main/c/sri.c),
[existing test cases](../../../main/c/sri_test.c), and
[build targets](../../../main/c/BUILD.bazel).
The checked-in fixture describes the intended digest format. The current
`sri_write` passes an uninitialized output buffer to `sri_calculate`, which
appends with `strncat`; this undefined behavior prevents a deterministic CLI
output guarantee until the buffer handling is repaired.

## Requirements

### Requirement: Preserve the intended named-digest interface and fixture

The command SHALL accept `--digest` (`-d`) and `--file` (`-f`). Its checked-in
fixture SHALL retain the intended output format: the digest name, a hyphen,
and the Base64-encoded OpenSSL digest, followed by a newline.

#### Scenario: Inspect the SHA-256 fixture for a known input

- **WHEN** the checked-in test case specifies input `Hello World` and digest `sha256`
- **THEN** its expected output is `sha256-pZGm1Av0IEBKARczz7exkNYsZb8LzaMrV7J32a2fFG4=` followed by a newline; the expectation alone does not establish reliable current CLI execution

### Requirement: Select an output destination

The command SHALL write to standard output by default and SHALL support
`--output` (`-o`) to write to the specified file instead.

#### Scenario: Output file is requested

- **WHEN** the caller supplies an output path that the command opens successfully
- **THEN** the command passes that file stream to its writer instead of standard output; the buffer limitation still applies to the written bytes

### Requirement: Report invalid inputs as command failures

The command SHALL return a nonzero exit status and report a diagnostic when a
required digest or input file option is missing, the input file cannot be
opened, or OpenSSL cannot resolve the requested digest.

#### Scenario: Digest name is unknown

- **WHEN** the input file opens but OpenSSL does not recognize the supplied digest name
- **THEN** the command reports the digest failure to standard error and returns a nonzero status

#### Scenario: Input path is unavailable

- **WHEN** the input file cannot be opened
- **THEN** the command identifies the input-open failure and returns a nonzero status

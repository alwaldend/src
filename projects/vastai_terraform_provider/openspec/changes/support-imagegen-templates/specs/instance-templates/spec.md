## Purpose

Expose validated template lookup and explicit template selection to reusable
Terraform modules without introducing live deployment side effects.

## ADDED Requirements

### Requirement: Read an exact template

The provider SHALL resolve exactly one matching template and expose its image
and whether Docker options consist only of environment assignments.

#### Scenario: Invalid template response

- **WHEN** the API returns no template, ambiguous results, or incomplete metadata
- **THEN** lookup fails without rental creation or response-body disclosure

### Requirement: Select a template for instance creation

The instance resource SHALL accept an optional creation-only template hash.

#### Scenario: Template hash changes

- **WHEN** a caller changes a known template hash on an existing instance
- **THEN** Terraform plans replacement

### Requirement: Bounded optional startup compression

The instance resource SHALL enforce a 16,384-byte startup command limit before
creation. Explicit gzip compression SHALL preserve script contents and refresh
stability without exposing script values in diagnostics. Compression SHALL
remain opt-in for images with the documented shell utilities.

#### Scenario: Large compressible script

- **WHEN** the caller enables compression for a compressible startup script
- **THEN** the API receives a bounded command that runs the original script and refresh SHALL preserve the configured script

#### Scenario: Oversized encoded command

- **WHEN** the encoded command still exceeds the API limit
- **THEN** validation or creation SHALL fail before a rental request and without disclosing the command

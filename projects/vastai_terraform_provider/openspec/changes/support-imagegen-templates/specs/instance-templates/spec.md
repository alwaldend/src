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

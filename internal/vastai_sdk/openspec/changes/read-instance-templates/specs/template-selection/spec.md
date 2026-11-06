## Purpose

Supply validated template identity and Docker option classification to callers
using the same bounded transport as other Vast.ai operations.

## ADDED Requirements

### Requirement: Validate template identity and options

Template lookup SHALL require one exact hash match with complete image and
environment metadata. It SHALL classify non-environment Docker options and
malformed quoting as unsafe for environment-only consumers.

#### Scenario: Template contains a published port

- **WHEN** template options include a Docker port or unknown option
- **THEN** the returned environment-only classification is false

#### Scenario: Malformed identity response

- **WHEN** template results are missing, duplicated, incomplete, or mismatched
- **THEN** lookup returns an error without including environment values

### Requirement: Forward explicit template selection

Creation SHALL send the caller's nonempty template hash through the generated
instance operation and omit it for existing image-only callers.

#### Scenario: Image-only rental

- **WHEN** the caller supplies no template hash
- **THEN** the API request remains compatible with existing rental behavior

# configuration-filesystem Specification

## Purpose

Let application callers supply storage for all application-owned filesystem operations, enabling configuration workflows without a host filesystem.

## Requirements

### Requirement: Application storage is supplied at construction

The application SHALL accept a filesystem at construction and use it for every application-owned file read, write, directory operation, path inspection and archive staging/publication operation. Omitting it SHALL retain existing OS filesystem behavior, including the usable zero-value application. Server and backend-library storage implementations SHALL retain their existing contracts.

#### Scenario: Independent application storage

- **WHEN** two applications receive different filesystems containing the same paths
- **THEN** configuration reads and writes remain independent and no operation falls back to host files

#### Scenario: Native backend boundary

- **WHEN** a collection workflow delegates file paths to its existing backend
- **THEN** application-owned checks and staging use the supplied filesystem and the backend remains responsible for accessing those paths under its existing contract

### Requirement: Configuration workflows operate entirely in supplied storage

Configuration export, loading, discovery, validation and JSON/TOML dumps SHALL operate entirely within the supplied filesystem, including configuration access from collection planning, apply, state import and export. Custom filesystem paths SHALL use that filesystem's namespace without consulting the host working directory or home directory.

#### Scenario: Memory-backed hierarchy round trip

- **WHEN** a caller exports nested decks, note types and notes to memory-backed storage
- **THEN** directory and root-file dumps return the complete declarations, derived titles and default deck paths, preserve content and leave disk untouched

#### Scenario: Storage failures

- **WHEN** supplied storage is read-only, missing an input, contains malformed declarations or fails a file operation
- **THEN** the application returns the relevant error with its cause and never retries against host storage

### Requirement: Filesystem protections survive injection

Configuration loading SHALL reject symbolic links when exposed by the supplied filesystem, including links in ancestor paths. Application export and build SHALL resolve aliases in the supplied filesystem and retain source/output conflict checks, staged publication and cleanup behavior.

Internal application owners SHALL use the shared filesystem package for file operations. Rooted filesystems SHALL reject paths and resolved links that leave their base, treat equivalent virtual path spellings as the same source, and retain temporary names in their own namespace.

#### Scenario: Rooted configuration boundary

- **WHEN** configuration output or a declared resource traverses outside an injected BasePathFs, directly or through a symlink
- **THEN** the operation fails with a permission error and leaves outside storage unchanged

#### Scenario: Nested rooted state publication

- **WHEN** collection export writes configuration and state through directly nested rooted filesystems
- **THEN** state is published under its virtual path and no temporary state file remains

#### Scenario: Source aliases retain filesystem semantics

- **WHEN** rooted source and destination names differ only by a leading separator, or a source link target contains an intermediate symlink followed by a parent component
- **THEN** source conflicts are rejected before staging and intermediate links resolve before the parent component is applied

#### Scenario: Build in memory

- **WHEN** a staged build uses memory-backed storage and its apply callback succeeds or fails
- **THEN** successful output is published there, failed output remains unchanged and staging is removed in both cases

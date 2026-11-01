## REMOVED Requirements

### Requirement: Resolve declared tools through AL configuration

**Reason**: AL no longer manages repository tool executables.
**Migration**: Run the ordinary Bazel tool target or its explicitly installed command.

### Requirement: Reuse validated executable cache entries

**Reason**: Private executable caching has been removed.
**Migration**: Use explicit installation to update standalone host tools.

### Requirement: Install misses through a configured Bazel method

**Reason**: There are no AL cache misses or implicit installation.
**Migration**: Run the tool's install target explicitly.

### Requirement: Preserve tool execution identity

**Reason**: AL no longer starts cached tools.
**Migration**: Invoke the installed command directly.

## ADDED Requirements

### Requirement: Explicit standalone tool installation

Repository tools formerly cached by AL SHALL expose ordinary local build and
explicit installation targets. Installation SHALL copy a standalone executable
atomically to the user's local binary directory, with an explicit directory
override for isolated use. A failed source read SHALL preserve the previous
installation. No private executable cache, download target or selection flag
SHALL be required.

#### Scenario: Install and replace a tool

- **WHEN** the install target runs successfully
- **THEN** the destination contains the current executable with mode 0755 and
  can run independently of Bazel outputs.

#### Scenario: Failed update

- **WHEN** a source executable cannot be read
- **THEN** installation fails and the previous installed executable remains.

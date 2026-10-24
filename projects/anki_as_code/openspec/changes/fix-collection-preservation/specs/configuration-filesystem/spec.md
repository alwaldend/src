## ADDED Requirements

### Requirement: Export destinations belong to the caller

Export SHALL write declared configuration files and leave unrelated entries
alone. Callers SHALL own selection of output locations and existing contents.
The writer SHALL NOT receive protected source identities or reject outputs based
on source containment, equality, symbolic links or shared file identity.
Ordinary filesystem errors and rooted namespace boundaries SHALL still apply.

#### Scenario: Export beside a collection

- **WHEN** the caller chooses an output directory containing its collection and
  an existing declared configuration file
- **THEN** export overwrites the declaration and preserves unrelated entries
  without rejecting the directory merely because it contains the source

### Requirement: Public paths preserve filesystem traversal order

Public collection source and output paths SHALL resolve intermediate symbolic
links before applying subsequent parent components. Configuration loading SHALL
retain its existing rejection of symbolic links in input paths.

#### Scenario: Collection selected through a symlink followed by a parent

- **WHEN** a supplied collection path includes a symbolic link followed by `..`
- **THEN** collection access selects the filesystem-resolved collection rather
  than the collection selected by lexical cleanup

#### Scenario: Configuration input traverses a symbolic link

- **WHEN** a configuration input contains a symbolic link followed by `..`
- **THEN** loading rejects the link rather than silently loading a different tree

#### Scenario: Build output is absent after a linked parent

- **WHEN** the build output includes a symbolic link followed by `..` and the
  resolved output does not exist
- **THEN** build publishes in the resolved existing parent and preserves any
  existing file at the destination selected by lexical cleanup

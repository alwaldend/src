## ADDED Requirements

### Requirement: Allow preserving nested source paths

`docs_filegroup` SHALL provide an opt-in mode that keeps each source's path
relative to its package, so a package whose sources sit in subdirectories is
not forced to flatten identically named files onto one destination. The
default SHALL remain flattened so existing published paths do not change.

#### Scenario: A package holds identically named sources in subdirectories

- **WHEN** a package declares documentation whose sources include two files
  with the same basename in different subdirectories and opts into preserving
  paths
- **THEN** the macro packages both files at distinct package-relative
  destinations instead of failing analysis

#### Scenario: A package keeps the default behavior

- **WHEN** a package declares documentation without opting into preserving
  paths
- **THEN** each source keeps its flattened basename under the archive prefix,
  as before

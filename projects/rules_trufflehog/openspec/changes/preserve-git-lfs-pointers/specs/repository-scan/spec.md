## ADDED Requirements

### Requirement: Scan Git history containing LFS pointers

The scanner's temporary Git checkout SHALL retain Git LFS pointer files and
scan the Git history without downloading the associated LFS payloads.

#### Scenario: A Git LFS payload is unavailable

- **WHEN** a scanned repository contains a Git LFS pointer whose payload is
  unavailable to the local clone
- **THEN** the scan covers the selected Git history and applies its normal
  findings policy

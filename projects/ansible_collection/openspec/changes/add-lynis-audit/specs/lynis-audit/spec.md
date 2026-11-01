## Purpose

Provide an explicitly invoked, bounded host security audit with a pinned scanner
and private evidence retained on the controller.

## ADDED Requirements

### Requirement: Explicit bounded audit

The role SHALL require a packaged archive and an absolute private controller
report directory. It SHALL execute Lynis as root with a finite timeout and SHALL
NOT install a service, schedule, or automatic remediation.

#### Scenario: Required destination omitted

- **WHEN** the caller omits the controller report directory
- **THEN** the role fails before copying or executing the scanner

#### Scenario: Audit exceeds its time limit

- **WHEN** the scanner exceeds the configured runtime
- **THEN** the scanner is terminated and its exit status and available evidence
  are retained before the failure is reported

### Requirement: Private retained evidence

The role SHALL retain the audit log, report, console output, and exit status in a
unique controller result directory. Directories SHALL have mode 0700 and reports
SHALL have mode 0600. Warnings SHALL NOT prevent collection. Remote cleanup SHALL
remove only the role's temporary directory; upstream private temporary files
remain under the scanner's own cleanup control.

An incomplete private receipt SHALL precede remote work. The role SHALL record
the allocated remote recovery path before preparing the scanner and SHALL retain
meaningful command transport failures without requiring an exit code. Native async polling SHALL supply scanner status only from a finished valid result. It SHALL
attempt collection of available evidence after command failure or unreachability.
Remote scratch SHALL be removed only after successful inspection of every
evidence path, successful fetch of every existing evidence file, persistence
of the collection receipt, and a known terminal scanner status. Collection failures SHALL retain scratch and fail.

#### Scenario: Scanner reports warnings

- **WHEN** Lynis reports warnings
- **THEN** the reports and exit status are collected before remote cleanup

#### Scenario: Evidence inspection or transfer fails

- **WHEN** any report stat or fetch fails or is unreachable
- **THEN** collection of other available files is attempted, the private recovery
  path and diagnosis are retained, remote scratch is preserved, and the role fails

#### Scenario: Scanner command is unreachable

- **WHEN** command execution returns an unreachable result without an exit code
- **THEN** the private receipt retains the transport diagnosis and recovery path,
  collection is attempted, and the audit is incomplete

### Requirement: Evidence-based completion

Completion SHALL require exit 0 or 78, successful retention of all existing
evidence, and a fully fetched report with `finish=true` and a positive integer
`lynis_tests_done`. Failed and timed-out scans SHALL retain available partial
artifacts before failure is reported.

#### Scenario: Accepted exit with incomplete report

- **WHEN** exit 0 or 78 accompanies an absent, unfinished, or zero-test report
- **THEN** the receipt records incomplete status and the role fails after retention

### Requirement: Pinned local execution

The scanner SHALL come from the approved upstream release with SHA256 integrity
enforced by Bazel. The role SHALL disable report uploads, plugins, and scanner
update checks and SHALL use the packaged configuration instead of host profiles.
The target SHALL provide Python supported by the packaged Ansible runtime.
The role SHALL retain `audit system` with all groups and noninteractive
`--cronjob`, and SHALL skip the DNF/YUM checks that can mutate package caches:
`PKGS-7350`, `PKGS-7352`, `PKGS-7354`, `PKGS-7383`, `PKGS-7384`, `PKGS-7386`, and
`PKGS-7387`. The role documentation SHALL describe the lost coverage and SHALL
NOT present a hardening audit as proof of absence of compromise.
It SHALL explain that the pinned malware tests mainly detect tools and processes,
without executing a comprehensive antivirus or rootkit scan.

#### Scenario: Host has an existing scanner profile

- **WHEN** the host has a separately installed Lynis profile
- **THEN** this audit uses the packaged profile and role overrides

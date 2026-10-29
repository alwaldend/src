## Purpose

Run official native ClamAV streaming scans with private transport, bounded
execution, and retained evidence without claiming complete host coverage.

## ADDED Requirements

### Requirement: Native private streaming

The role SHALL unpack the official payload into private endpoint scratch and
invoke `clamdscan --stream` through a private Unix socket. It SHALL NOT install
host packages, use a custom scanner helper, expose TCP, reboot, quarantine, or
delete scanned files. Unsupported runtimes SHALL fail without a fallback.

#### Scenario: Compatible endpoint

- **WHEN** the native runtime preflight succeeds
- **THEN** configured roots are scanned with the official native client
- **AND** symlink following is disabled and exclusions are recorded

#### Scenario: Unsupported endpoint

- **WHEN** the payload cannot execute on the endpoint architecture or libc
- **THEN** the receipt records failure and no compatibility helper is installed

#### Scenario: Payload extraction fails

- **WHEN** extraction fails or its connection becomes unreachable
- **THEN** preparation stops before runtime probing and preserves the extraction
  diagnostic in the private receipt
- **AND** a missing executable is not misreported as proven architecture or
  library incompatibility

#### Scenario: Kernel device nodes are present

- **WHEN** the host exposes block or character device nodes beneath `/dev`
- **THEN** native Ansible metadata discovery identifies and excludes those nodes
  without opening their contents or following symlinks
- **AND** their paths are retained and regular files beneath `/dev/shm` remain eligible
- **AND** failed or incomplete device discovery prevents scanning

### Requirement: Oversized streams cannot succeed as prefixes

The client and caller's controller stream limits SHALL both be 104857601 bytes.
The role's `max_file_bytes` SHALL remain 104857600 and its client limit SHALL
retain the `max_file_bytes + 1` expression. The caller's controller SHALL use
`MaxFileSize 104857600`, `HeuristicAlerts yes`, and `AlertExceedsMax yes`.
Oversized input SHALL reach the engine and produce a native file-limit
detection, allowing traversal to continue without a successful verdict for a
truncated prefix.

#### Scenario: Input exceeds the client limit

- **WHEN** a stable readable file is larger than 104857601 bytes
- **THEN** the engine reports `Heuristics.Limits.Exceeded.MaxFileSize FOUND`
- **AND** native traversal continues without a successful prefix verdict
- **AND** raw diagnostics and native status are retained

#### Scenario: Oversized input precedes EICAR

- **WHEN** an oversized file is encountered before a later EICAR fixture in
  the same native traversal
- **THEN** retained native logs show the limit detection followed by the
  later EICAR detection
- **AND** the size limit does not abort traversal with a protocol error

### Requirement: Isolated transport lifecycle

The role SHALL manage only its dedicated SSH master and forward, preserve caller
options/authentication, and leave unrelated masters intact. It SHALL allocate
the dedicated ControlPath entirely on the controller and establish its
persistent master before any endpoint command, including boot identity,
tempfile creation, and Ansible module bootstrap. It SHALL capture caller
arguments once before overrides, restore them during finalization, use explicit
own-master mux operations, and avoid `meta: reset_connection`. Evidence
collection and endpoint cleanup SHALL use that same authenticated master before
it is closed. Local connections SHALL skip SSH setup.

#### Scenario: Master startup fails

- **WHEN** the dedicated master cannot be established for a remote endpoint
- **THEN** endpoint commands are not attempted through another connection
- **AND** the master-start result is retained in the private incomplete receipt

#### Scenario: Forward release

- **WHEN** forward setup returned success with exit code 0
- **THEN** recovery explicitly cancels only that forward on its own master
- **AND** a forward that was never attempted needs no release operation

#### Scenario: Private native IPC paths

- **WHEN** a remote scan starts
- **THEN** its dedicated SSH control directory uses private native user runtime storage
- **AND** it is independent of daemon-owned runtime storage so recovery survives daemon exit
- **AND** one final native assertion derives completion from registered results

#### Scenario: Unsafe configuration path or excessive control path

- **WHEN** the report root, scan root, or exclusion contains `#`, or ControlPath
  exceeds 90 bytes
- **THEN** the role rejects the input before scanning
- **AND** the accepted control path reserves 17 bytes for OpenSSH's temporary suffix

#### Scenario: Scan outlasts certificate validity

- **WHEN** the dedicated master stays connected through scan completion
- **THEN** evidence is collected through it without fresh authentication

#### Scenario: Master or collection fails

- **WHEN** retention or termination cannot be confirmed
- **THEN** recovery paths remain in a private incomplete receipt
- **AND** unverified scan evidence is preserved

### Requirement: Bounded execution and honest native evidence

The role SHALL bound execution, retain native logs, validate native async status, and
parse final summary counts. Exit 1 SHALL NOT hide reported errors. Execution
completion SHALL be separate from host cleanliness and coverage.
Receipts SHALL retain `coverage_complete=false` and `clean=null` rather than
infer complete coverage or host cleanliness from native exit status. Native
infected counts SHALL include limits and encrypted-content heuristics without
being presented as malware-only counts. The private receipt SHALL retain the
master-start result and both boot-command results, including failed,
unreachable, or skipped outcomes. The initial boot result SHALL be saved before
its success assertion.

#### Scenario: Async execution does not provide a finished native result

- **WHEN** the async result is unfinished, unreachable, or lacks a valid return code
- **THEN** no scanner exit code is inferred and workflow completion remains false
- **AND** available evidence and recovery paths remain retained

#### Scenario: Explicit higher-throughput scan

- **WHEN** the caller disables background scheduling and enables multiscan
- **THEN** the native client uses normal CPU/I/O scheduling and `--multiscan`
  while retaining `--stream`, scan roots, exclusions, deadlines, and evidence
- **AND** receipts record both controls; defaults retain background scheduling
  and sequential submission
- **AND** multiscan exit 1 retains findings but marks native completion uncertain
  instead of treating a detection status as proof of completed traversal

#### Scenario: Boot identity cannot be verified

- **WHEN** a boot-command result is failed or unreachable
- **THEN** its diagnostics remain in the private receipt
- **AND** unchanged boot identity and workflow completion are not established

#### Scenario: No connection for final boot check

- **WHEN** neither local transport nor a successful dedicated master is available
- **THEN** the boot-after command is skipped and its result is retained privately

#### Scenario: Detection and error coexist

- **WHEN** native exit is 1 and the summary reports errors
- **THEN** detections remain visible and workflow completion is false

#### Scenario: Successful execution with limited coverage

- **WHEN** execution and evidence collection succeed
- **THEN** summary counts and raw logs remain available
- **AND** no host-clean or complete-coverage assertion is made

#### Scenario: Deadline expires

- **WHEN** GNU timeout terminates the scanner
- **THEN** its actual status and available logs are retained
- **AND** the workflow remains incomplete

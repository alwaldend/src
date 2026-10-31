---
title: Lynis audit
description: Bounded on-demand security audit with private controller reports
---

`alwaldend.main.lynis_audit` extracts the packaged CISOfy source archive into a
root-owned 0700 temporary directory and runs an audit as root.
The Linux target must have a POSIX shell, tar, GNU `timeout`, and a remote Python
version supported by the packaged Ansible runtime. The role requires privilege
escalation to root and does not support Ansible check mode. It does not install
Python or provide a legacy transport fallback.

Supply `lynis_audit_archive` as the controller path to the archive packaged from
`//third_party/com_github_cisofy_lynis:archive`. Supply
`lynis_audit_report_directory` as an existing absolute controller directory,
owned by the controller user with mode 0700. In repository operations, place it
under the task's ignored `out/` directory. `lynis_audit_timeout_seconds` defaults
to 900, with a further 15-second termination grace period and a maximum of 3600.

Each invocation creates a unique 0700 result directory and exposes its path as
`lynis_audit_result_directory`. It retains 0600 `console.log`, `lynis.log`,
`lynis-report.dat`, and `result.json` files. The role writes an incomplete receipt
first, then records the remote scratch recovery path before preparing or running
the scanner. The receipt and `lynis_audit_exit_code` record the scanner status;
the exit code remains null when command execution supplies none. Transport and
collection errors are retained in the private receipt, without exposing their
contents through Ansible output.

The audit runs through native Ansible async polling, with its job data under
the private audit directory. A finished result supplies the scanner return code;
missing, unreachable or unfinished results supply none. This separates scanner
status from the SSH polling process without an encoded stdout protocol.
The console is still redirected to a private file, preserving partial output.

Completion requires exit 0 (without warnings) or 78 (with warnings), successful
retention of all existing evidence, and a fully fetched report containing
`finish=true` and a positive integer `lynis_tests_done`. An accepted exit alone
does not establish completion. Other statuses, unfinished or zero-test reports,
and collection failures fail the role after available evidence is collected.
Timeout statuses 124 or 137 indicate incomplete evidence; available partial
artifacts are retained. Empty files can indicate failure before scanner output.

The command retains `audit system` with all test groups and `--cronjob`. In the
pinned source, `--quick` only skips input pauses; it does not select fewer tests.
`--cronjob` already enables that behavior and noninteractive handling of an
existing scanner PID file, so a separate `--quick` is redundant. It does not
install a cron job. Individual checks still depend on host applicability and
the profile exclusions below.

The command uses the extracted source and profiles through `--usecwd`, disables
plugins, uploads, and scanner update checks, and disables package repository
refreshes. The packaged source contains no runnable plugins, only a plugin
template and its README. The shared profile also skips DNF/YUM checks that can
mutate caches independently of `refresh-repositories=no`:

- `PKGS-7350`: installed-package inventory through DNF.
- `PKGS-7352`: DNF security-update advisories and vulnerable-package reporting.
- `PKGS-7354`: DNF package-database integrity checks.
- `PKGS-7383`: YUM repository/package-update management checks.
- `PKGS-7384`: YUM utilities, duplicate packages, and package-database problems.
- `PKGS-7386`: YUM security support and missing security-update checks.
- `PKGS-7387`: enabled YUM repositories' GPG-signature configuration.

Consequently, this audit omits those DNF/YUM package inventory, security advisory,
database, repository, and signature checks; other applicable package checks can
still run. Package maintenance belongs to the caller's dedicated workflow.
Normal system and DNS checks remain enabled; this is not an offline sandbox.
The role installs no package, service, schedule, or automatic remediation.
Lynis assesses system hardening and does not establish absence of compromise.
The pinned `include/tests_malware` primarily checks for installed security tools
and related processes; it does not run a comprehensive antivirus or rootkit
scan. A completed audit therefore does not certify the host free of malware.

Lynis itself creates private `/tmp/lynis.*` files and `/var/run/lynis.pid`, which
its normal exit handler removes. A forced termination can leave these behind.
The role removes its allocated remote directory only after all evidence paths
have been inspected, every existing evidence file has been fetched successfully,
the collection receipt has been saved, and terminal scanner status is known.
This permits cleanup after a failed or timed-out scan when its partial artifacts
were retained. A stat or fetch
failure leaves the remote directory in place and fails the role; consult
`result.json` for the private `remote_directory` recovery path and diagnostics.
The recorded path may already be absent after successful collection and cleanup.
An unreachable host or interrupted controller can prevent collection or cleanup;
the initial receipt and recorded recovery path survive. The role never removes
an existing scanner's PID or temporary files. Reports contain host details and
should remain private.

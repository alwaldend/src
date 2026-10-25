## Purpose

Make shared host configuration preserve access, activate updates reliably,
and maintain the intended separation between configuration and service data.

## ADDED Requirements

### Requirement: Native automatic security updates

The host role SHALL configure automatic_updates using the distribution's native
DNF or APT package and packaged systemd timers. Automatic reboots SHALL remain
disabled. Unsupported package managers and immutable OSTree hosts SHALL be
skipped with an explicit message. Repeated unchanged deployment SHALL converge.

#### Scenario: Enable a supported mutable host

- **WHEN** the host role runs with automatic updates enabled
- **THEN** the native package is installed and its update timer is enabled and active
- **AND** DNF selects security updates or APT retains packaged allowed origins
- **AND** no custom cron job or systemd timer is installed

#### Scenario: Disable automatic installation

- **WHEN** automatic_updates_enabled is false
- **THEN** the managed installation timer is stopped and disabled if present
- **AND** disabling does not install a missing package

### Requirement: Distribution-aware update activation

The update role SHALL detect required reboots using the target distribution's
supported mechanism and SHALL reject detection errors instead of treating them
as a negative result.

#### Scenario: Fedora has installed updates requiring reboot

- **WHEN** DNF reports a required reboot without a Debian reboot-marker file
- **THEN** the role schedules the reboot according to its configured policy

### Requirement: Safe proxy activation

An unchanged Traefik role run SHALL preserve the service process. Changes to
the executable, unit, or startup configuration SHALL activate only when changed, with recovery on startup failure.
Malformed candidate configuration SHALL be rejected before replacing active
files. Failed activation SHALL restore available prior configuration and report
failure rather than leave a previously working proxy stopped.

#### Scenario: Repeat an unchanged proxy deployment

- **WHEN** configuration, executable, unit, and enrollment settings are unchanged
- **THEN** the running proxy is neither restarted nor sent a reload signal

#### Scenario: Create enrollment state on the first deployment

- **WHEN** the configured issuer needs enrollment and its private cache is missing
- **THEN** the requested binding is saved with administrator-only permissions
- **AND** the next unchanged deployment reuses it without requesting another binding

#### Scenario: Reject an invalid candidate

- **WHEN** candidate configuration is malformed or cannot start the proxy
- **THEN** deployment fails and the previous configuration remains available
- **AND** a previously running proxy is restored when prior files exist

### Requirement: Read-only proxy configuration

The service identity SHALL NOT be able to replace configuration, its parent
directories, or its executable. Only declared runtime data paths under the proxy directory SHALL be
writable inside the service sandbox.

#### Scenario: Inspect proxy filesystem permissions

- **WHEN** the role configures a system service
- **THEN** the service user cannot write configuration or its ancestors
- **AND** it can write its required ACME data and log directories

### Requirement: Activate renewed SSH certificates

Replacing a signed SSH host certificate SHALL notify a validated SSH daemon
reload, including when the daemon configuration itself is unchanged.

#### Scenario: Renew an existing host identity

- **WHEN** a signed host certificate changes
- **THEN** subsequent connections receive the renewed certificate after handlers

### Requirement: Independent operating-system recovery

The download content mount SHALL NOT be required for OS boot. Content-serving
and deduplication services SHALL require the mount and SHALL NOT fall back to
the underlying system-disk directory.

#### Scenario: Content disk is missing during boot

- **WHEN** the content device cannot be mounted
- **THEN** OS startup and administrative access can continue
- **AND** content-dependent services cannot start against the bare mount point

### Requirement: Explicit firewall rule retirement

The firewall role SHALL delete only rules explicitly listed for retirement,
preserve unlisted rules, and install desired access before deleting retired
rules. It SHALL use the standard UFW Ansible module without custom output
parsing or a complete firewall reset.
An enabled policy SHALL also enable the native UFW service so the retained
rules load at boot before subsequent policy updates.

#### Scenario: Reboot a host with an enabled firewall

- **WHEN** the configured host reboots
- **THEN** its retained policy is active without another deployment
- **AND** a subsequent unchanged firewall deployment succeeds

#### Scenario: Retire a previously allowed port

- **WHEN** an operator adds the port's rule to `firewall_delete_rules`
- **THEN** its allowance is removed on the next run
- **AND** desired management access and unlisted rules remain present

#### Scenario: Remove an allow declaration without retiring it

- **WHEN** an allow entry is removed without a corresponding deletion entry
- **THEN** the existing rule remains in place
- **AND** arbitrary comments do not affect rule management

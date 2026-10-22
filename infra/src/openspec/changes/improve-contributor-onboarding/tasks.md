Deferred; implementation has not started. These unchecked tasks are future
work; this planning request authorizes no shared or live host configuration change.

## 1. Establish the supported setup path

- [ ] 1.1 Locate authoritative setup source and existing host trust overrides;
      record a supported host baseline, private-CA compatibility needs, and a
      small public smoke target with an inspectable expected output.
- [ ] 1.2 Prepare an isolated setup exercise before configuration changes for
      a supported host lacking the hardcoded truststore path; record its current
      result and the explicit custom-trust case needed for acceptance.

## 2. Improve and verify contributor setup

- [ ] 2.1 Add the concise root entry path and explicit host trust configuration,
      preserving the owning setup source and any required host_bot IaC; verify
      documentation links and inspect generated configuration without deployment.
- [ ] 2.2 Execute the documented smoke workflow on the exact candidate in the
      isolated supported environments; save a repeatable report with host,
      revision, commands, TLS verification outcome, and representative output.

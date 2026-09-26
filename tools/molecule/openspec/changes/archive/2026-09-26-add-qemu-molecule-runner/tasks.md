## 1. Define runner acceptance and tool inputs

- [x] 1.1 Write smoke, failure, cancellation, and concurrent-run acceptance
      scenarios before runner implementation. Verify each runner requirement has
      an observable result and retained artifact.
- [x] 1.2 Resolve compatible Molecule, QEMU, and seed-image tools through
      repository dependency workflows and reuse the Fedora image pin. Verify
      tool versions, image format and checksum, and a disposable guest's boot,
      SSH, package access, and selected accelerator; retain the smoke artifact.

## 2. Implement shared Molecule and QEMU integration

- [x] 2.1 Add the shared rule and Go runner using declared package mappings,
      tools, scenarios, and guest inputs. Verify the staged collection preserves
      executable binaries and loads the exact candidate with isolated Ansible
      configuration and no inherited infrastructure credentials.
- [x] 2.2 Implement shared create/destroy playbooks and owned QEMU lifecycle
      commands with fresh overlays, seed material, virtual disks, and loopback
      ports. Verify successful teardown, repeatable destroy, bounded boot
      failures, and distinct resources for two concurrent smoke runs.
- [x] 2.3 Add supervision, handled-signal cleanup, restricted manifests, and
      sanitized artifact collection. Inject a verification failure and a
      cancellation; verify non-success status, retained diagnostics, no surviving
      owned process, no leaked test credentials, and documented abrupt-exit
      recovery.
- [x] 2.4 Declare opt-in execution, cache, platform, timeout, and any proven
      runtime network/sandbox requirements. Verify two explicit invocations
      execute twice, missing prerequisites fail clearly, and ordinary build
      actions remain sandboxed and offline.

## 3. Integrate and deliver the tooling contract

- [x] 3.1 Add the owning README, Bazel packages, and OpenSpec validation
      registration. Verify the tool and its specs are discoverable, and that
      the public test interface and runtime exceptions are documented.
- [x] 3.2 Run complete runner acceptance and the linked collection consumers
      against the same candidate. Verify retained outputs, clean teardown,
      package mapping fidelity, and no credential leakage.
- [x] 3.3 Run OpenSpec, affected semantic lint, and repository quality checks;
      commit and publish through repo-delivery. Archive only after runner
      acceptance passes and verify consumer links point to the stable spec or
      archived change.

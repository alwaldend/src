## 1. Update the installer

- [x] 1.1 Pin the newer installer and verify its checksum and shell syntax.
- [x] 1.2 Verify the Ansible package and exercise the installer with a simulated systemd service manager.

The Harbor Ansible syntax test passed. Its runfiles mapped the pinned script
to the K3s role's installer path. The packaged script generated a systemd
unit in an isolated fixture, skipped the restart on an unchanged run, and
restarted when forced.

# Forgejo runner acceptance

Run `bazel test //projects/ansible_collection/test/forgejo_runner:registration_test`.
The packaged Ansible registration tasks run as the current user in disposable
directories. A local executable records dummy registration and restart events;
it never contacts Forgejo or Vault. The fixture renders the actual configuration
and service templates and retains them with `registration.json` in Bazel's test
outputs.

Cases cover initial registration, repeated deployment without bootstrap inputs,
missing inputs, credential-bearing command failure, bounded registration,
file permissions, default and advanced configuration, and protected paths.
API-issued identity cases cover import without registration, preservation on
repeat deployment, invalid ID/UUID/token rejection and name/instance mismatch.
The real package manager, account creation, systemd, cross-account ACL transfer
and remote Forgejo protocol are outside this fixture.

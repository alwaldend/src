## 1. Managed configuration

- [x] 1.1 Install a repository-scoped helper and connect it to the dev1 Git configuration.
- [x] 1.2 Validate the packaged helper and Git credential protocol with an isolated fixture.

Git's credential exchange accepted the configured repository and rejected a
different host, repository, and protocol, a readable credential file, and a
symlink. The packaged Ansible playbook passed syntax validation. The native
login authorized a Git dry run and publication of the theme branch.

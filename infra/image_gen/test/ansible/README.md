# Profile workflow installation

`bazel test //infra/image_gen/test/ansible:assets_test` runs the real packaged
deployment playbook with the `workflows` tag against a disposable local-only
inventory for every configured profile. The selected NoobAI profile also runs
the actual playbook JSON translation with tiny local model/Git substitutes, real packaged
first-party sources and a fixture API/restart command. Python dependency
installation and Vulkan discovery are disabled only in this fixture. It verifies that the installed native
JSON and App contents exactly match their generated sources, and emits a
`profile-workflows.json` receipt. All untagged tasks, including downloads, Vault
reads and service operations, are excluded. The shared role's separate local
fixture covers model destinations, Git pins, source updates and idle protection.

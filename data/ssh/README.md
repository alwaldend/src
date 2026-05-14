---
title: ssh
description: Ssh keys
---

## Links

- CA generation: https://developer.hashicorp.com/vault/docs/secrets/ssh/signed-ssh-certificates
- OpenPGP ssh key: https://github.com/drduh/YubiKey-Guide/tree/master

## ca_clients

CA for SSH clients

Update:
```sh
bazel run //infra/dc1/vault/tf:update_ssh_ca_clients
```

{{< readfile file="ca_clients.pub" code="true" >}}

## ca_servers

CA for SSH servers

Update:
```sh
bazel run //infra/dc1/vault/tf:update_ssh_ca_servers
```

{{< readfile file="ca_clients.pub" code="true" >}}

## simeonwarren

Update: manual export from the yubikey

{{< readfile file="simeonwarren@alwaldend.com-2026-03-22-openpgp.pub" code="true" >}}


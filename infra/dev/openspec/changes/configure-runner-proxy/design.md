## Decision

Proceed with reuse of the historical T3 Code proxy on
`yc1.nodes.threexui.alwaldend.com:40201`, account `t3code`, stored at
`secrets/alwaldend.com/vault1/approles/src_infra_threexui/http_proxies`.
Historical source: commit `a6c803838`,
`users/simeonwarren/t3code/ansible/group_vars/all.yaml`.

Bazel's Java HTTP proxy is a `java.net.Proxy.Type.HTTP` connection even for
an `https://` proxy URL. Supply a loopback HTTP endpoint with a systemd
socket-activated OpenSSL TLS forwarder. This reuses software already required
by the host, validates the remote certificate, carries no credentials in the
unit and avoids a new proxy implementation or dependency. Direct HTTPS proxy
variables alone do not cover Java's downloader; changing remote listeners or
deploying a new VPN has greater scope and is unnecessary.

Inject only the selected password into the controller, URL-encode it, and
merge job environment into the existing runner configuration. Preserve the
role's default timeout, cache and trust settings. Suppress task output and
diffs when serializing the credential-bearing runner configuration.

## Failure cases and verification

Before implementation, consider absent credentials, encoding of reserved URL
characters, a stopped forwarder, untrusted/mismatched TLS certificates, public
downloads accidentally bypassing the proxy, local domains accidentally using
it, changed runner defaults, and secrets leaking through logs or arguments.
Verify rendered Ansible with dummy credentials; syntax-check and build the
packaged playbook; validate unit files and exercise the TLS stream forwarding
with local fixtures. Live acceptance must fetch an external dependency and
contact the local forge from the runner account. Policy and deployment are
not implied by successful source validation.

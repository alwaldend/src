---
title: Firewall
description: Setup firewall
tags:
  - ansible_role
---

The role installs incoming and routed allows, deletes explicitly retired rules,
and enables UFW. All rule operations use `community.general.ufw`. Rule comments
are descriptive and do not determine which rules are deleted.

For an enabled policy, the role enables and starts `ufw.service` before
updating rules. This restores an existing policy after reboot even when an
older deployment left the service disabled. Firewalld is stopped first because
its service conflicts with UFW. The role also masks firewalld when its package
is absent, so later dependency installation cannot enable the conflicting
service. A disabled policy also disables UFW startup.

Declare retired rules in `firewall_delete_rules`, which defaults to an empty
list. Removing an entry from `firewall_allow_port_from` or
`firewall_allow_forward` alone leaves its installed rule in place. Add a matching
deletion entry and retain it for hosts that may have missed earlier deployments.
Unlisted rules remain; the role does not reset the firewall.

Each deletion requires `rule` (`allow`, `deny`, `reject`, or `limit`) and accepts
the module's `name`, `proto`, `port`, `from_port`, `src`, `dest`, `direction`, and
`route` parameters. Sources and destinations default to `any`; `route` defaults
to `false`. Omitted protocol, port, direction, and application parameters retain
the module defaults. Match the original rule specification, including any
destination or routing constraint. Application profiles use `name` without
`proto` or `port`.

```yaml
firewall_delete_rules:
  - rule: allow
    port: "8080"
    proto: tcp
  - rule: allow
    name: mDNS
    dest: 224.0.0.251
  - rule: allow
    route: true
    src: 192.0.2.0/24
    dest: 198.51.100.0/24
```

Install replacement management access in `firewall_allow_port_from` before
retiring an SSH profile; desired allows run before deletions. Do not also declare
a retired rule as a desired allow. The download consumer declares TCP 22 and
explicitly retires inherited `SSH` and `mDNS` profile rules, including the IPv4
and IPv6 multicast destinations.

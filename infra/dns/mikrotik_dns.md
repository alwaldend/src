---
title: Mikrotik DNS
linkTitle: Mikrotik DNS
description: Records declared for the internal dc1 view served by RouterOS
---

Every record this repository declares for the internal `dc1` view, which RouterOS serves for that network.

This page is generated from the declarations; run `bazel run //infra/dns/cmd/dump -- --write` after changing them. The authoritative source is each owner's `dnsconfig.json`.

| Domain | Type | TTL | Value | Declaration |
| --- | --- | --- | --- | --- |
| @ | A | 300 | 185.199.108.153 | infra/dns/dnsconfig.json |
| @ | A | 300 | 185.199.109.153 | infra/dns/dnsconfig.json |
| @ | A | 300 | 185.199.110.153 | infra/dns/dnsconfig.json |
| @ | A | 300 | 185.199.111.153 | infra/dns/dnsconfig.json |
| @ | AAAA | 300 | 2606:50c0:8000::153 | infra/dns/dnsconfig.json |
| @ | AAAA | 300 | 2606:50c0:8001::153 | infra/dns/dnsconfig.json |
| @ | AAAA | 300 | 2606:50c0:8002::153 | infra/dns/dnsconfig.json |
| @ | AAAA | 300 | 2606:50c0:8003::153 | infra/dns/dnsconfig.json |
| @ | MX | 300 | 10 mail.protonmail.ch. | infra/dns/dnsconfig.json |
| @ | MX | 300 | 20 mailsec.protonmail.ch. | infra/dns/dnsconfig.json |
| @ | TXT | 300 | _globalsign-domain-verification=0QBJgVV_uwcFLTi1Rot3bb1LyJ5uW1WD0ygvIS4OM5 | infra/dns/dnsconfig.json |
| @ | TXT | 300 | protonmail-verification=bdcd133d3f472fa17f66328950d02fbeae1bef75 | infra/dns/dnsconfig.json |
| @ | TXT | 300 | v=spf1 include:_spf.protonmail.ch ~all | infra/dns/dnsconfig.json |
| _dmarc | TXT | 300 | v=DMARC1; p=quarantine; adkim=s | infra/dns/dnsconfig.json |
| _dmarc.simplelogin | TXT | 10800 | v=DMARC1; p=quarantine; pct=100; adkim=s; aspf=s | infra/dns/dnsconfig.json |
| automation.openhands | A | default | 192.168.10.92 | infra/openhands/dnsconfig.json |
| bm1.dc1 | A | default | 192.168.1.222 | infra/mikrotik/dnsconfig.json |
| bm2.dc1 | A | default | 192.168.1.216 | infra/pve/dnsconfig.json |
| bm2.dc1 | AAAA | default | fd2e:546d:5738:0:365a:60ff:fe08:6daa | infra/pve/dnsconfig.json |
| bm3.dc1 | A | default | 192.168.1.218 | infra/vault/dnsconfig.json |
| bm3.dc1 | AAAA | default | fd2e:546d:5738:0:e2be:3ff:fe2b:9a1a | infra/vault/dnsconfig.json |
| canvas.openhands | A | default | 192.168.10.90 | infra/openhands/dnsconfig.json |
| cloud | CNAME | default | nas.alwaldend.com. | infra/nas/dnsconfig.json |
| cloudinit-test.vm.pve1.dc1 | A | default | 192.168.10.10 | infra/pve/dnsconfig.json |
| dc1.automation.openhands | A | default | 192.168.10.92 | infra/openhands/dnsconfig.json |
| dc1.canvas.openhands | A | default | 192.168.10.90 | infra/openhands/dnsconfig.json |
| dc1.cloud | A | default | 192.168.1.209 | infra/nas/dnsconfig.json |
| dc1.host-bot.simeonwarren.users | A | default | 192.168.1.210 | users/simeonwarren/host_bot/dnsconfig.json |
| dc1.server.openhands | A | default | 192.168.10.91 | infra/openhands/dnsconfig.json |
| dc1.t3code.host-bot.simeonwarren.users | A | default | 192.168.1.210 | users/simeonwarren/host_bot/dnsconfig.json |
| dkim._domainkey.simplelogin | CNAME | 10800 | dkim._domainkey.simplelogin.co. | infra/dns/dnsconfig.json |
| dkim02._domainkey.simplelogin | CNAME | 10800 | dkim02._domainkey.simplelogin.co. | infra/dns/dnsconfig.json |
| dkim03._domainkey.simplelogin | CNAME | 10800 | dkim03._domainkey.simplelogin.co. | infra/dns/dnsconfig.json |
| flux | A | default | 192.168.10.60 | infra/flux/dnsconfig.json |
| forgejo | A | default | 192.168.10.40 | infra/forgejo/dnsconfig.json |
| git | A | default | 192.168.10.40 | infra/forgejo/dnsconfig.json |
| harbor | A | default | 192.168.10.50 | infra/harbor/dnsconfig.json |
| host-bot.simeonwarren.users | A | default | 192.168.1.210 | users/simeonwarren/host_bot/dnsconfig.json |
| host1.automation.openhands | A | default | 192.168.10.92 | infra/openhands/dnsconfig.json |
| host1.canvas.openhands | A | default | 192.168.10.90 | infra/openhands/dnsconfig.json |
| host1.cloud | CNAME | default | host1.nas.alwaldend.com. | infra/nas/dnsconfig.json |
| host1.flux | A | default | 192.168.10.60 | infra/flux/dnsconfig.json |
| host1.forgejo | A | default | 192.168.10.40 | infra/forgejo/dnsconfig.json |
| host1.harbor | A | default | 192.168.10.50 | infra/harbor/dnsconfig.json |
| host1.ingress | A | default | 81.26.185.118 | infra/ingress/dnsconfig.json |
| host1.nas | A | default | 192.168.1.209 | infra/nas/dnsconfig.json |
| host1.pve1.dc1 | CNAME | default | bm2.dc1.alwaldend.com. | infra/pve/dnsconfig.json |
| host1.server.openhands | A | default | 192.168.10.91 | infra/openhands/dnsconfig.json |
| host1.threexui | A | default | 192.168.10.80 | infra/threexui/dnsconfig.json |
| host1.vault.dc1 | A | default | 192.168.1.218 | infra/vault/dnsconfig.json |
| host1.xcp-ng | A | default | 192.168.1.213 | infra/xcp_ng/dnsconfig.json |
| host1.xoa.xcp-ng | A | default | 192.168.1.206 | infra/xcp_ng/dnsconfig.json |
| ingress | A | default | 81.26.185.118 | infra/ingress/dnsconfig.json |
| mail._domainkey.yandex | TXT | 300 | v=DKIM1; k=rsa; t=s; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCcYzFVgkeDOhaIIkWM8gNQjxVsv0/aXfU+ax5urB5y6hA6lSjRnjRo6tm0bXbkOJf41GmiwMNgdXpwRtzgzAlX1i2aJbtEr4b9jzibEGLQ7Cvqs44bOYES9f/K3ueQpnvdTOJmFqlRReFL7ZrUyDFCoQ7f4+7h4i8s01cCcRrt5wIDAQAB | infra/dns/dnsconfig.json |
| nas | A | default | 192.168.1.209 | infra/nas/dnsconfig.json |
| njalla1.nodes.threexui | A | default | 45.142.141.133 | infra/threexui/dnsconfig.json |
| njalla1.nodes.threexui | AAAA | default | 2a0a:3840:8078:141:0:2d8e:8d85:1337 | infra/threexui/dnsconfig.json |
| openid.flux | CNAME | default | flux.alwaldend.com. | infra/flux/dnsconfig.json |
| operator.flux | CNAME | default | flux.alwaldend.com. | infra/flux/dnsconfig.json |
| protonmail._domainkey | CNAME | 300 | protonmail.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch. | infra/dns/dnsconfig.json |
| protonmail2._domainkey | CNAME | 300 | protonmail2.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch. | infra/dns/dnsconfig.json |
| protonmail3._domainkey | CNAME | 300 | protonmail3.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch. | infra/dns/dnsconfig.json |
| pve | A | default | 192.168.1.216 | infra/pve/dnsconfig.json |
| router1.dc1 | A | default | 192.168.1.1 | infra/mikrotik/dnsconfig.json |
| router1.dc1 | AAAA | default | fd2e:546d:5738::1 | infra/mikrotik/dnsconfig.json |
| runner1.forgejo-runner | A | default | 192.168.10.100 | infra/forgejo_runner/dnsconfig.json |
| server.openhands | A | default | 192.168.10.91 | infra/openhands/dnsconfig.json |
| simplelogin | MX | 10800 | 10 mx1.simplelogin.co. | infra/dns/dnsconfig.json |
| simplelogin | MX | 10800 | 20 mx2.simplelogin.co. | infra/dns/dnsconfig.json |
| simplelogin | TXT | 10800 | sl-verification=bxfzzfjiggzsxyzxhhmkmjqkaskjgy | infra/dns/dnsconfig.json |
| simplelogin | TXT | 10800 | v=spf1 include:simplelogin.co ~all | infra/dns/dnsconfig.json |
| switch1.dc1 | A | default | 192.168.1.254 | infra/mikrotik/dnsconfig.json |
| t3code.host-bot.simeonwarren.users | A | default | 192.168.1.210 | users/simeonwarren/host_bot/dnsconfig.json |
| threexui | A | default | 192.168.10.80 | infra/threexui/dnsconfig.json |
| vault | A | default | 192.168.1.218 | infra/vault/dnsconfig.json |
| vault.dc1 | A | default | 192.168.1.218 | infra/vault/dnsconfig.json |
| www | CNAME | 300 | alwaldend.com. | infra/dns/dnsconfig.json |
| www-staging | CNAME | 300 | alwaldend.github.io. | infra/dns/dnsconfig.json |
| xcp-ng | A | default | 192.168.1.213 | infra/xcp_ng/dnsconfig.json |
| xoa.xcp-ng | A | default | 192.168.1.206 | infra/xcp_ng/dnsconfig.json |
| yandex | MX | 21600 | 10 mx.yandex.net. | infra/dns/dnsconfig.json |
| yandex | TXT | 300 | v=spf1 redirect=_spf.yandex.net | infra/dns/dnsconfig.json |
| yandex | TXT | 300 | yandex-verification: b83672f59b3dbe16 | infra/dns/dnsconfig.json |
| yc.threexui | NS | default | ns1.yandexcloud.net. | infra/threexui/dnsconfig.json |
| yc.threexui | NS | default | ns2.yandexcloud.net. | infra/threexui/dnsconfig.json |
| yc1.nodes.threexui | CNAME | default | host1.nodes.yc.threexui.alwaldend.com. | infra/threexui/dnsconfig.json |

---
title: Cloudflare DNS
linkTitle: Cloudflare DNS
description: Records declared for the public zone served by Cloudflare
---

Every record this repository declares for the public zone, which Cloudflare serves.

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
| canvas.openhands | CNAME | default | ingress.alwaldend.com. | infra/openhands/dnsconfig.json |
| cloud | CNAME | default | ingress.alwaldend.com. | infra/nas/dnsconfig.json |
| dc1.automation.openhands | A | default | 192.168.10.92 | infra/openhands/dnsconfig.json |
| dc1.canvas.openhands | A | default | 192.168.10.90 | infra/openhands/dnsconfig.json |
| dc1.cloud | A | default | 192.168.1.209 | infra/nas/dnsconfig.json |
| dc1.host-bot.simeonwarren.users | A | default | 192.168.1.210 | users/simeonwarren/host_bot/dnsconfig.json |
| dc1.server.openhands | A | default | 192.168.10.91 | infra/openhands/dnsconfig.json |
| dc1.t3code.host-bot.simeonwarren.users | A | default | 192.168.1.210 | users/simeonwarren/host_bot/dnsconfig.json |
| dkim._domainkey.simplelogin | CNAME | 10800 | dkim._domainkey.simplelogin.co. | infra/dns/dnsconfig.json |
| dkim02._domainkey.simplelogin | CNAME | 10800 | dkim02._domainkey.simplelogin.co. | infra/dns/dnsconfig.json |
| dkim03._domainkey.simplelogin | CNAME | 10800 | dkim03._domainkey.simplelogin.co. | infra/dns/dnsconfig.json |
| forgejo | CNAME | default | ingress.alwaldend.com. | infra/forgejo/dnsconfig.json |
| git | CNAME | default | ingress.alwaldend.com. | infra/forgejo/dnsconfig.json |
| host1.ingress | A | default | 81.26.185.118 | infra/ingress/dnsconfig.json |
| ingress | A | default | 81.26.185.118 | infra/ingress/dnsconfig.json |
| int.forgejo | A | default | 192.168.10.40 | infra/forgejo/dnsconfig.json |
| int.vault | A | default | 192.168.1.218 | infra/vault/dnsconfig.json |
| mail._domainkey.yandex | TXT | 300 | v=DKIM1; k=rsa; t=s; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCcYzFVgkeDOhaIIkWM8gNQjxVsv0/aXfU+ax5urB5y6hA6lSjRnjRo6tm0bXbkOJf41GmiwMNgdXpwRtzgzAlX1i2aJbtEr4b9jzibEGLQ7Cvqs44bOYES9f/K3ueQpnvdTOJmFqlRReFL7ZrUyDFCoQ7f4+7h4i8s01cCcRrt5wIDAQAB | infra/dns/dnsconfig.json |
| njalla1.nodes.threexui | A | default | 45.142.141.133 | infra/threexui/dnsconfig.json |
| njalla1.nodes.threexui | AAAA | default | 2a0a:3840:8078:141:0:2d8e:8d85:1337 | infra/threexui/dnsconfig.json |
| pages | A | default | 185.199.108.153 | projects/alwaldend.com/dnsconfig.json |
| pages | AAAA | default | 2606:50c0:8000::153 | projects/alwaldend.com/dnsconfig.json |
| protonmail._domainkey | CNAME | 300 | protonmail.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch. | infra/dns/dnsconfig.json |
| protonmail2._domainkey | CNAME | 300 | protonmail2.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch. | infra/dns/dnsconfig.json |
| protonmail3._domainkey | CNAME | 300 | protonmail3.domainkey.djgwfzcu5fgjtpoijqqomgifmqj6zeiuwdd4mzim4hrxab3zsgwkq.domains.proton.ch. | infra/dns/dnsconfig.json |
| simplelogin | MX | 10800 | 10 mx1.simplelogin.co. | infra/dns/dnsconfig.json |
| simplelogin | MX | 10800 | 20 mx2.simplelogin.co. | infra/dns/dnsconfig.json |
| simplelogin | TXT | 10800 | sl-verification=bxfzzfjiggzsxyzxhhmkmjqkaskjgy | infra/dns/dnsconfig.json |
| simplelogin | TXT | 10800 | v=spf1 include:simplelogin.co ~all | infra/dns/dnsconfig.json |
| t3code.host-bot.simeonwarren.users | CNAME | default | ingress.alwaldend.com. | users/simeonwarren/host_bot/dnsconfig.json |
| vault | CNAME | default | ingress.alwaldend.com. | infra/vault/dnsconfig.json |
| www | CNAME | 300 | alwaldend.com. | infra/dns/dnsconfig.json |
| yandex | MX | 21600 | 10 mx.yandex.net. | infra/dns/dnsconfig.json |
| yandex | TXT | 300 | v=spf1 redirect=_spf.yandex.net | infra/dns/dnsconfig.json |
| yandex | TXT | 300 | yandex-verification: b83672f59b3dbe16 | infra/dns/dnsconfig.json |
| yc.threexui | NS | default | ns1.yandexcloud.net. | infra/threexui/dnsconfig.json |
| yc.threexui | NS | default | ns2.yandexcloud.net. | infra/threexui/dnsconfig.json |
| yc1.nodes.threexui | CNAME | default | host1.nodes.yc.threexui.alwaldend.com. | infra/threexui/dnsconfig.json |

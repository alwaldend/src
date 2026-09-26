---
title: Download DNS
description: Cloudflare and RouterOS aliases for static hosting
---

Owns the component's `dnsconfig.json` through the shared DNS module. Public
download requests follow the delegated Yandex host. The local host CNAME
follows the apex owner's canonical local A record. The apex and `www` keep
their existing owner; DNS cutover is a separately reviewed live operation.

This root selects only DNS credential injection and its own state backend.
It does not authenticate with either VM provider or require either VM's state.

## Local DNS preparation

Before configuring the local host by its inventory FQDN, prepare and commit a
reviewed intermediate `infra/dns/dnsconfig.json` revision that installs the
final `download_dc1` A record in the dc1 view and removes the old dc1 Pages
A/AAAA members. Keep the old public Pages A/AAAA records and all unrelated
records unchanged; split prior `dsp: ["all"]` declarations by view as needed.
Plan and apply that exact revision through `infra/dns/tf`, after separately
authorizing this local DNS change. Review that only the intended dc1 apex
records change and verify the local answer before continuing.

Apply this component's aliases after that preparation, then confirm
`host1.dc1.download.alwaldend.com` resolves to the address used by the VM.
Ansible connects using that unique hostname. The local website is in a
maintenance interval until host configuration and content publication finish;
the public website continues using Pages during this local preparation.

The local VM consumes the same apex A value from source. The host CNAME and
local download alias follow it, so renumbering has one maintained address.
The dc1 apex also declares the unusable AAAA address `::ffff`, following
[RouterOS's documented override](https://help.mikrotik.com/docs/spaces/ROS/pages/37748767/DNS).
An A record alone does not suppress upstream AAAA answers: while public DNS
still points to Pages, IPv6 clients would reach that site instead of the local
VM. This local override keeps clients on the configured IPv4 endpoint without
changing public DNS or relying on a guest's temporary IPv6 address.
The public apex follows `download.alwaldend.com`, which in turn follows the
component-owned Yandex host. Renaming that host therefore requires no copied
endpoint edit in the apex owner.

## Staged public apex cutover

The apex belongs to `infra/dns/tf`, separately from this component DNS root.
Its old GitHub Pages A/AAAA records and the new public CNAME have different
Terraform instance keys. Applying the final declaration directly can race
CNAME creation against address-record deletion. Cloudflare rejects that
coexistence; use two sequential, reviewed IaC revisions and saved plans.

Before either phase, complete the DNS owner's
[adoption procedure](../../dns/README.md#migration-and-recovery), verify both
hosts and content using address overrides, configure certificate issuers, and
authorize the apex cutover explicitly.
Record the prior source revision and recovery plan. Schedule a maintenance
window: the public apex has an intentional address gap between the phases,
and resolver caches can extend that interval.

1. Prepare and commit a reviewed intermediate `infra/dns/dnsconfig.json`
   revision with the old public apex A/AAAA members removed from the global view and
   `records.download_global` absent. Keep local records, `www`, mail, TXT,
   and unrelated records at their prior values. Through
   `//infra/dns/tf:tf.plan`, save and review a plan that deletes only the old
   public apex A/AAAA instances and creates no apex CNAME. Apply that exact
   saved plan through `//infra/dns/tf:tf.apply`. Verify completion and the
   authoritative Cloudflare inventory shows those A/AAAA records are gone.
2. Only after phase 1 succeeds, select the reviewed final declaration from
   this change, including `records.download_global`. Create a fresh saved
   plan through the same owner; review the CNAME creation and intended local
   apex changes while preserving mail/TXT/`www` and unrelated records. Apply
   that plan, then verify authoritative records and both DNS views, website
   responses, and redirects.

Keep all source/input files unchanged between each saved plan and its apply.
If a phase fails, stop and inspect the actual record/state inventory before
planning recovery. Restore the recorded prior declaration only through a
reviewed owning-root plan; remove any conflicting CNAME first if returning
to the old A/AAAA records. Never begin phase 2 based only on a planned deletion.

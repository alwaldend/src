---
title: Global DNS records
description: Cloudflare-only entrypoint for canonical DNS declarations
---

Global-only owners use this entrypoint with the [shared module API](../README.md)
to avoid requiring RouterOS configuration. The combined module also delegates
its Cloudflare resources here, so both entrypoints share one implementation.
Normalization returns all declared views, but this entrypoint owns only global
records. Mixed-view owners use the parent module to manage dc1 as well.

Enabled global records use the explicit `cloudflare_zone_id` when supplied.
Otherwise, the pinned `cloudflare_zone` data source resolves the configured
`zone` name and rejects missing or ambiguous matches. This requires zone-list
and zone-read access in addition to record permissions. Disabled ownership and
declarations without global records perform no zone lookup.

The optional zone ID must be known during planning so Terraform can select
whether discovery is needed. The repository's AL injection supplies a known
value or an empty string before Terraform starts. Zone-name discovery uses the
normalizer's canonical name, including case and final-dot handling.

`import_addresses` contains global addresses relative to this entrypoint.
`record_ids` exposes bound Cloudflare IDs for adoption and identity checks.

---
title: Airflow
date: 2026-11-26
description: Airflow workflow scheduling infrastructure in DC1
---

This project owns the Vault identity and DNS for `airflow.alwaldend.com`.
The [AppRole root](approle/README.md) declares the component's DNS access,
guest SSH signing role, Vault ACME role, and browser OIDC client. The
[DNS declarations](dns/dnsconfig.json) define the shared ingress name and
the DC1 guest address.

## Deployment

1. Apply `//infra/airflow/approle:tf.apply` and reconcile shared Vault
   memberships with `//infra/vault/tf:tf.apply`.
2. Apply `//infra/airflow/dns:tf.apply` to configure the DC1 and global DNS
   views. The ingress and guest serve the names when the Airflow service
   is deployed.

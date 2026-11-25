---
title: MLflow
date: 2026-11-25
description: Infrastructure for MLflow tracking and Vast.ai project execution
---

This project owns the Vault identity and DNS for `mlflow.alwaldend.com`. The
standalone [AppRole root](approle/README.md) declares the component's DNS
access, guest SSH signing role and Vault ACME role. The [DNS
declarations](dns/dnsconfig.json) define the public ingress name and the
DC1 guest address.

## Deployment

1. Apply `//infra/mlflow/approle:tf.apply` and reconcile shared Vault
   memberships with `//infra/vault/tf:tf.apply`.
2. Apply `//infra/mlflow/dns:tf.apply` when the guest and ingress are ready to
   receive traffic. The DNS project manages both the DC1 and public views.

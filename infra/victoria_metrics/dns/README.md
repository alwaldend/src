---
title: VictoriaMetrics DNS
---

The [DNS declaration](dnsconfig.json) gives the service hostname two records:
a public CNAME to shared ingress and a DC1 A record for the guest. A third
record, the public `dc1.victoria-metrics.alwaldend.com` A record, lets shared
ingress resolve the guest's private address through public DNS. Terraform
provisioning uses the same guest address.
Verify the address is available on the DC1 network before deployment.

Review DNS changes with `bazel run //infra/victoria_metrics/dns:tf.plan`.

# OpenCanary DNS

This Terraform root manages the site-local address, upstream alias, and public ingress alias in
[`dnsconfig.json`](dnsconfig.json). The guest uses the same declared address.
Review with `bazel run //infra/opencanary/dns:tf.plan`.

---
title: Grafana configuration
---

This Terraform root manages the VictoriaMetrics Prometheus data source and
repository dashboards through the Grafana API. It uses its own Vault HTTP
backend and requires DC1 access to the Grafana hostname.

## Dashboards

Place dashboard JSON at `dashboards/<section>/<project>/<dashboard>.json`.
For example, dashboards under `dashboards/infra/ingress/` appear in Grafana
under **infra → ingress**. Use a distinct, stable dashboard UID and the
`victoria-metrics` Prometheus data source UID. The relative file path is the
Terraform instance key. Preserve the UID and move the Terraform address when
renaming a file. Folder IDs derive from repository paths.

The dashboard under **infra → victoria_metrics** renders every panel in
VictoriaMetrics' existing Monitoring overview. Its queries and layout come
from the [VMUI dashboard source](../../victoria_metrics/ansible/files/monitoring-overview.json).

## Deployment

1. Apply the Grafana AppRole changes and deploy the Grafana Ansible playbook.
   The playbook enables nested folders and allows the Grafana process to reach
   the VictoriaMetrics HTTPS data source under SELinux.
2. Review `bazel run //infra/grafana/tf:tf.plan`, then run
   `bazel run //infra/grafana/tf:tf.apply`. The Grafana AppRole supplies the
   existing Grafana local administrator password and the VictoriaMetrics Basic
   Auth credentials. Terraform stores the data source password in the protected
   backend state.
3. Verify the data source and open the dashboards under **infra**. The Grafana
   dashboard requires the collector's native `grafana` scrape job.

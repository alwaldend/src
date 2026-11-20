---
title: VictoriaMetrics role
---

The `alwaldend.main.victoria_metrics` role installs and configures the
VictoriaMetrics single-node service on a systemd host. Run it with privilege
escalation and set `victoria_metrics_binary_src` to the executable path on
the Ansible controller.

The service binds to loopback and retains data for 30 days. To configure
scrape jobs, set `victoria_metrics_scrape_config_src` to a YAML template on
the Ansible controller.

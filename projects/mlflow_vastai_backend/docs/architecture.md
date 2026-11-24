---
title: Plugin architecture
description: Submission, execution, and result tracking across MLflow and Vast.ai
---

[![MLflow Vast.ai plugin architecture](assets/architecture.svg)](assets/architecture.svg)

The MLflow run ID identifies the workload across the plugin, coordinator, and
worker. The coordinator records its rental assignment before dispatch. The
worker records the outcome before reporting it.

The worker receives the project snapshot from MLflow's artifact store and runs
the entry point on the rental. The workload logs metrics and artifacts to the
run. The coordinator publishes confirmed terminal status to MLflow and retries
failed updates after the submitting client disconnects.

Terraform provisions and releases rentals. The coordinator discovers existing
rentals through the Vast.ai API.

See the [application state machines](../state-machines/) for process lifecycles
and job transitions.

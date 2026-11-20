## Context

The deployment playbook configures VictoriaMetrics after the guest is ready.

## Goals / Non-Goals

Move service configuration to a reusable role. Preserve the service settings
and keep guest provisioning and artifact packaging with the deployment.

## Decisions

The role accepts a controller-side executable path. The deployment supplies
the verified Bazel artifact through group variables, so using other roles in
the collection does not download VictoriaMetrics.

## Risks / Trade-offs

The role expects a Linux host running systemd and an executable that matches
the target architecture.

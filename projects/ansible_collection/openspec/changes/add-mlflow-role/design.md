## Decisions

Extract the service tasks, templates, and restart handlers from the MLflow
deployment playbook. The role receives its inputs through Ansible variables.
The caller provisions the host and configures the HTTPS route. The role
checks the local tracking server and coordinator after applying changes.

Callers with dedicated volumes can require those mounts in the systemd units.
The role owns the paths and startup checks. The caller owns the disks and
mounts.

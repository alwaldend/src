## Decisions

- Use the existing Helm provider with the Forgejo mirrors. The third-party
  inventories pin the chart archives and OCI image indexes for both controllers.
- Enable cert-manager CRDs. Use a ClusterIP service for Traefik because the
  cluster does not run the K3s ServiceLB controller.
- Keep issuer and route configuration with the workloads that need them.

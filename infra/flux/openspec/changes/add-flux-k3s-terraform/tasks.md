## 1. Separate bootstrap and workloads

- [x] 1.1 Package the XCP-ng and Ansible bootstrap and keep Flux workloads out of the K3s auto-deploy directory.
- [x] 1.2 Package the Vault OIDC bootstrap and inspect the Helm provider's token and TLS configuration offline.
- [x] 1.3 Pin the Helm provider and package the ordered workload releases. Validate the workload root and render the FluxInstance chart.
- [x] 1.4 Publish the API server CA during bootstrap and inject it for Terraform runs.
- [x] 1.5 Open and close the SSH forward automatically during Terraform runs.

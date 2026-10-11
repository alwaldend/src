## Source acceptance

The standalone architecture change has no Terraform, inventory, DNS, or secret
changes. Njalla contains njalla1 and work1; Yandex Cloud contains yc1. The shared
node-detail subgraph describes Traefik, 3x-ui, Xray, and Mullvad relationships.
Control-plane ACME uses Vault; data-plane ACME uses Let's Encrypt.

Retained Drawio pages are generated projections of the same owner source.
Validation receipts and rendered inspection artifacts remain under ignored
out/architecture. Live deployment belongs to the separately authorized work1 task.

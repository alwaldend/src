## 1. Traefik role

- [x] 1.1 Set the shared port default and declare listeners and firewall rules in deployment configuration.
- [x] 1.2 Validate the deployed configuration shape offline.

The packaged inventories resolved the metrics port and collector address
for all opted-in hosts. A local VictoriaMetrics process scraped the pinned
Traefik binary's dedicated listener successfully.

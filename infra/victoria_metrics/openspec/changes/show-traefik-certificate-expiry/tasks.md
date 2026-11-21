## 1. Dashboard

- [x] 1.1 Add the Traefik certificate expiration panel and verify its query with representative samples from multiple scrape jobs.

The pinned VictoriaMetrics server served the candidate dashboard and returned
the expected lifetimes for four certificate series from three scrape jobs,
including two serial numbers for one common name and a negative value for an
expired certificate.

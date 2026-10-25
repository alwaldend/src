## 1. Restore the maintained public declaration

- [x] 1.1 Restore the historical GitHub Pages address keys and values for the global view only; preserve every other declaration and update the deployment documentation.
- [x] 1.2 Validate the DNS declaration and prepare a committed candidate before live rollback.

## 2. Apply and verify the authorized rollback

- [x] 2.1 Review and apply a saved plan deleting only the public apex CNAME, then verify its authoritative removal.
- [x] 2.2 Review and apply a fresh plan creating only the eight public GitHub Pages address records; verify convergence, public HTTPS, and preservation of other records.
- [x] 2.3 Record the observed outcome and verification evidence for source handoff.

## 1. Establish the deployment contract and fixture

- [ ] 1.1 Write the behavioral failure matrix before implementation, including host trust, path safety, interrupted uploads, conflicting versions, malformed archives, concurrency, and disk exhaustion; verify it covers each requirement.
- [ ] 1.2 Prepare an isolated SSH/Nginx end-to-end fixture with repeatable invocation and result artifacts; verify it can observe uploaded bytes, active-site content, and unchanged paths after failure.
- [ ] 1.3 Add owner OpenSpec packaging and aggregate validation registration; verify strict validation discovers this workspace through repository targets.

## 2. Add SSH metadata and ordinary publication

- [ ] 2.1 Extend deployment metadata and Bazel inputs, regenerate bindings through their owner, and update command dependency selection; verify SSH-only invocation needs no Soras while an existing OCI fixture still works.
- [ ] 2.2 Implement explicit SSH target selection with established authentication and host verification; verify an unknown/mismatched host is refused and credentials do not enter generated public metadata.
- [ ] 2.3 Implement temporary upload, byte verification, and final rename into the public project/version directory; verify ordinary archives remain files and no additional `files/` path appears.
- [ ] 2.4 Implement idempotent identical uploads and conflicting-content rejection; verify interrupted/multi-file publication reports exact results and preserves existing bytes.

## 3. Activate and redeploy website releases

- [ ] 3.1 Implement designated `.tar.gz`/`.zip` site extraction with the documented archive-root contract; verify valid archives publish their HTML/assets and remain downloadable.
- [ ] 3.2 Implement validated staging, per-project coordination, archive identity bookkeeping, and atomic site selection; verify malformed/escaping/link archives, full disks, and competing deployments preserve the selected site.
- [ ] 3.3 Implement existing-release selection through the same deploy operation; verify selecting an older validated release changes only the link and introduces no rollback command or automatic deletion.

## 4. Integrate and validate

- [ ] 4.1 Connect public artifact URLs and project/version mapping to existing release metadata; verify generated release pages address the exact uploaded directories.
- [ ] 4.2 Run the end-to-end matrix plus existing relevant OCI checks; retain candidate identity, exact commands, remote-tree results, checksums, and HTTP evidence.
- [ ] 4.3 Document explicit environment selection and deployment errors, then run semantic lint and required repository checks; verify the delivered command and examples match the final implementation.

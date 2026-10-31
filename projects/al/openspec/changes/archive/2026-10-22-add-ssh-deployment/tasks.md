## 1. Establish the deployment contract and fixture

- [x] 1.1 Write the behavioral failure matrix before implementation, including host trust, path safety, interrupted uploads, ordinary file updates, malformed archives, concurrency, and disk exhaustion; verify it covers each requirement.
- [x] 1.2 Prepare an isolated SSH/Nginx end-to-end fixture with repeatable invocation and result artifacts; verify it can observe uploaded bytes, active-site content, and unchanged paths after failure.
- [x] 1.3 Add owner OpenSpec packaging and aggregate validation registration; verify strict validation discovers this workspace through repository targets.

## 2. Add SSH metadata and ordinary publication

- [x] 2.1 Extend deployment metadata and Bazel inputs, regenerate bindings through their owner, and update command dependency selection; verify SSH-only invocation needs no Soras while an existing OCI fixture still works.
- [x] 2.2 Implement explicit SSH target selection with established authentication and host verification; verify an unknown/mismatched host is refused and credentials do not enter generated public metadata.
- [x] 2.3 Implement temporary upload, rsync transfer verification and final rename into the public project/version directory; verify ordinary archives remain files and no additional `files/` path appears.
- [x] 2.4 Implement rsync updates and preservation of unrelated retained files; verify interrupted/multi-file publication reports its itemized results and preserves the selected site.

## 3. Activate and redeploy website releases

- [x] 3.1 Implement designated `.tar.gz`/`.zip` site extraction with the documented archive-root contract; verify valid archives publish their HTML/assets and remain downloadable.
- [x] 3.2 Implement validated staging, per-project coordination, retained version selection, and atomic site selection; verify malformed/escaping/link archives, full disks, and competing deployments preserve the selected site.
- [x] 3.3 Implement existing-release selection through the same deploy operation; verify selecting an older validated release changes only the link and introduces no rollback command or automatic deletion.

## 4. Integrate and validate

- [x] 4.1 Connect public artifact URLs and project/version mapping to existing release metadata; verify generated release pages address the exact uploaded directories.
- [x] 4.2 Run the end-to-end matrix plus existing relevant OCI checks; retain candidate identity, exact commands, remote-tree results, checksums, and HTTP evidence.
- [x] 4.3 Document explicit environment selection and deployment errors, then run semantic lint and required repository checks; verify the delivered command and examples match the final implementation.

## 5. Review follow-up

- [x] 5.1 Preserve literal SSH and rsync arguments with the existing POSIX quoting approach, following the user correction to remove the unapproved library.
- [x] 5.2 Make upload, extraction from published archives, and link selection independently selectable; verify the expanded SSH/HTTP failure matrix.
- [x] 5.3 Wrap propagated errors with operation context and preserved causes, and package the repo-go skill with offline evaluation and discovery checks.
- [x] 5.4 Record the local-dependency-first rule and Go skill routing in AGENTS.md; verify the host branch reaches its merged trunk and publication remains an independent trunk PR.

Exact final validation and review-thread state are recorded in the delivery
receipts referenced by evidence.md.

## 6. Dependency approval correction

- [x] 6.1 Require explicit approval for each new external dependency in AGENTS.md and provide a self-contained approval procedure in the dependency skill.
- [x] 6.2 Remove the unapproved quoting dependency and preserve the existing literal-argument behavior without a new library.
- [x] 6.3 Recheck the SSH/HTTP fixture against the corrected quoting implementation; final skill and delivery results remain in the candidate-bound receipts.

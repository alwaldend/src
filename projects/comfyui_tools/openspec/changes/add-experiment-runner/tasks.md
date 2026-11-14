## 1. Contracts and failure cases

- [x] 1.1 Record the generic behavior, ownership boundary and failure scenarios before implementation; verify proposal, design and capability delta exist.
- [x] 1.2 Add command-level E2E fixtures for buffering, duplicate identity, restart, ambiguous acknowledgment, ownership-scoped cancellation and escaping media; verify they emit a repeatable result artifact.

## 2. Generic runner

- [x] 2.1 Add the producer/CLI configuration, workflow snapshots, identity exclusions, private journal and run lock; verify duplicate and restart E2E scenarios.
- [x] 2.2 Add buffered submission, independent completion observation, media collection and local reports; verify slow retrieval does not empty the execution buffer.
- [x] 2.3 Add status, graceful stop and explicit owned-job cancellation; verify unrelated fixture jobs remain untouched.
- [x] 2.4 Add the packaged command and neutral usage documentation; verify its CLI help and affected Bazel targets.

## 3. Skill and acceptance

- [x] 3.1 Add the model-agnostic ComfyUI experiments skill, UI metadata and offline Promptfoo coverage; verify skill validation and exact discovery.
- [x] 3.2 Exercise the generic command on the already-authorized live endpoint with an external task-local producer; retain prompt IDs, queue observations, distinct workflows and original outputs as acceptance evidence.
- [x] 3.3 Complete formatting, semantic checks, affected E2E tests and OpenSpec validation against the final candidate; record limitations.

Publication state is recorded by the AL delivery receipt and pull request;
those records own the candidate, remote head and publication result.

## 4. Review corrections

- [x] 4.1 Preserve observable and cancellable ownership after observer transport failures; verify recovery and owned cancellation E2E cases.
- [x] 4.2 Resume valid journals larger than 64 MiB without truncating their JSON; verify a large persisted-journal E2E case.
- [x] 4.3 Recheck stop requests after slow producer completion and before submission; verify no extra native job is accepted.
- [x] 4.4 Preserve exact numeric producer notes across initial decoding, restart and gallery generation; verify a full-width integer seed.

- [x] 4.5 Bound transient artifact retrieval retries without resubmitting generation; verify a command E2E case matching the observed live download interruption.

- [x] 4.6 Decode terminal native outcomes from unsuccessful wait envelopes; distinguish execution failure, external cancellation and nonterminal timeout in command E2E cases.
- [x] 4.7 Reconcile best-effort cancellation against the exact native identifier; retain ownership after a rejected interrupt and collect outputs from a completion race.
- [x] 4.8 Preserve an external cancellation failure when a later observer failure is recovered; verify no replacement submission occurs on restart.

## 5. Direction development

- [x] 5.1 Replace endless refill guidance with one job per ready direction and explicit scrap/refine/fork evaluation, independent lineage and pause persistence.
- [x] 5.2 Update generic documentation and offline skill cases; validate packaging and report the limits of offline behavioral coverage.

# TigerStyle for Go

## Safety

- Bound input sizes, queues, concurrent work, retries, and traversal depth
  where they can grow with inputs or load. Define what happens when each limit
  is reached. Keep control flow explicit and use fixed-width integers where
  representation is part of a protocol or storage format.
- Preallocate slices and size maps when their expected size is known and
  bounded. Enforce resource limits directly. For example, check a batch limit
  before appending to a slice because append can grow past its initial
  capacity. Choose channel buffering deliberately and account for shared
  backing arrays when passing slices across ownership boundaries.
- Check invariants where they protect correctness. Return contextual errors
  for invalid input and recoverable failures. Reserve assertions for
  programmer errors and keep required runtime checks active in every build.
  Use Go fuzzing to exercise invariants when appropriate under the repository
  testing policy.

## Performance

Estimate the expected resource costs of hot paths early, including bandwidth,
latency, memory, and compute. Use measurements to evaluate batching,
allocation, and copying decisions. Batch work when it reduces per-operation
overhead within the applicable memory and latency limits.

## Readability

Keep interfaces small, keep state close to its use, and name values for their
role and units using Go identifier conventions.

## Sources

- [TigerStyle](https://tigerstyle.dev/)
- [Go Tiger Style](https://github.com/solvicode/Go-Tiger-Style)

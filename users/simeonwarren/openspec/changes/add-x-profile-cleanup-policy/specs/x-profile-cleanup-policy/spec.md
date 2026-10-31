## Purpose

Provide a reusable user-owned cleanup policy that combines a 30-day age
threshold with preservation of explicitly named X threads.

## ADDED Requirements

### Requirement: Age-based cleanup with thread preservation

The policy SHALL select posts, replies, and reposts strictly older than 30 days
for ordinary deletion. It SHALL preserve the six user-supplied threads named
main, mlops, prompt factory, anki, literally me, and books, with comments naming
each seed. It SHALL use the cleaner's existing thread traversal and preservation
precedence, without forced deletion or a restricted discovery window.

#### Scenario: Old unpreserved activity

- **WHEN** an owned post, reply, or repost is older than 30 days and is not preserved
- **THEN** the policy selects it for deletion or repost removal

#### Scenario: Recent activity or preserved thread

- **WHEN** activity is at most 30 days old, or a post belongs to a preserved thread
- **THEN** the policy does not select it for removal

### Requirement: Isolated single-post policy

A separate policy SHALL select only the exact user-supplied post for the live
browser demonstration, without any age rules or thread expansion. Discovery
SHALL be limited to the target day to avoid unrelated historical pages.

#### Scenario: Single-post plan

- **WHEN** the demonstration policy is planned against an inventory containing its target
- **THEN** the plan contains exactly one ordinary post deletion for that target

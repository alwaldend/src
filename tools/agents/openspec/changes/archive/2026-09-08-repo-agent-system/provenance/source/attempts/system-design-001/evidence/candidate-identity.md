# Exact candidate identity

## Subject

- Baseline source revision:
  `1423dce5fab45ce5223caeb6a24791bf1a2cc3ff`
- Subject: 20 maintained non-goal files changed by the design phase
- Manifest algorithm: SHA-256 over each file's exact bytes; the manifest is
  sorted in the intentional repository-entry order below
- Manifest digest:
  `d6e29c5644ca9f834f3a70d2856f3069673d88a19f59ee38dee599a06c40c8ca`

The checkpoint-owned files below
`projects/agents/goals/repo-agent-system/` are excluded from this manifest
because this checkpoint atomically changes them. Their identity and integrity
are owned by the goal resource versions and artifact digests and are checked
by `goal validate` after closure.

## Manifest

```text
6bd1e4046d9cbd8911c6a45ac02f788ebc24d587bb2bf97d06d9fb72d60be564  .gitattributes
672af72139da188a627445574b08716bdb4c0a4625229523d8172ed1e756626d  AGENTS.md
06689711fae574d4ed6018ef1380fa48ac5129a3b93cac6ab1d2b50804230969  README.md
8f96ac5685d3ac36f4eb4e7e0194ce49fcbfae780f62ab0cbab0feb9b8db9798  projects/README.md
3390d9faa88a756c624d78bee15aa0717b2aa90c96dd8bd341757f859a1c1923  infra/README.md
e421087e7114b20cf52dd82c559fc94d6891dfc9a9e6effe7699ece3330734e0  tools/README.md
4006a430089d9670bb9760daaa24d62b5c53dafc2d8239ad831745d261d4dbfa  data/README.md
66134cc8748a7b3e7a6bfeb9eaa0fa18dcaa84c85325b0620197f03417304737  third_party/README.md
bc1b93cd9d6165d3923f8ef0991d1dee434932921929ff62f262cefae19eb45a  users/README.md
8aa245bba451d782d649d8011e48bdb6ac92dc3199e2f0d1e816292d8dfece6f  projects/agents/BUILD.bazel
f6a662eaa4b25cd59fd649efbcdc2def18ee23f117dc39f1ac8e0f9e02c571cf  projects/agents/README.md
74c6eefe2d0f89fea89becf387ad1e75da61cd0f1dbb517a98b1d1907de5ca52  projects/agents/docs/BUILD.bazel
2a5f81bbffe4709e960a229f7c2d9e727b607384646712d6b87a256a324d1f91  projects/agents/docs/current-state.md
6af4166e7ed5b9bdc44856b2c7c23d30dc95fd10ff19f0fced1ecaa6698ecdcd  projects/agents/docs/architecture.md
fc1dc0197888108323b6ac3dfcc0f268db6f78a04e1fdd08ef51cbc1f3b4af6c  projects/agents/docs/roadmap.md
8a646a2001d6ecd699a2414a7191e1e17a4b06f311526c1e7affc18fea17b8d8  projects/agents/goals/BUILD.bazel
60bd3c558eeacf9b06bd0a0d052fe1d189b2223defe52d711d5ca597a9271f01  projects/agents/goals/README.md
9186b527517eabc6f9218ec4ad74ca123dbc2805f509474eb5fcad41cab1718d  projects/agents/skills/answer-question/SKILL.md
c7195ce9902e958acb002d4cca78cc7df030d42f94c4e67857e768b9f71dfc89  projects/agents/skills/repo-secrets/SKILL.md
a2df07718dc1743c0e8ec3afd0b1c39cc4e12408bd4208c658d3034950f8e3ad  projects/alwaldend.com/hugo.toml
```

`sha256sum --check` reported `OK` for all 20 entries immediately before this
attempt was closed.

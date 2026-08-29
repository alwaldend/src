# Evidence manifest

[Back to goal](README.md)

## Frozen inputs

- Tracked Sisyphus blend SHA-256:
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.
- Tracked standalone Reimu blend SHA-256:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Original Sisyphus image SHA-256:
  `3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`.

## Attempt 01

| Candidate | Blend SHA-256 | Render SHA-256 | Verdict |
|---|---|---|---|
| V01 | `e482952bff46e3fea6b6d67b90ffc360bada6f45f280d4283f26db647c38e9d0` | `7346705b4290c61bebd3e8e1ab3e72cdaca1ecafcff88e7476d96a4233f41bdb` | Reject |
| V02 | `477fb8a5f38461247c65622c02f4cb64c95da8113ebdab5378537cb8cd244209` | `07ec1137586d7a0bfd9870fc2e3e49c0c6c115b23568a9739db3a762bbb829c5` | Reject |

Both review images were inspected directly. Saved-file hierarchy and protected
hash regressions remain unverified, but the pixel failures already reject both
candidates.

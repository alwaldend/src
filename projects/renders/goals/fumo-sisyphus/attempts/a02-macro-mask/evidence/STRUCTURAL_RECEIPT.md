# A02 structural and preservation receipt

Pinned Blender 5.2.1 clean-opened candidate
`c4b2b1118e8c215e0787703ad700dac1a665742cc9862315c6dadd48cccf5bd2`.
The scene used `CAM_SCENE_QUOTE_FREE` at 512 by 529 and exposed the named
environment, lighting, placeholder-only, reference-only, and rockwork
collections. `REFERENCE_ONLY` was render-disabled and contained the packed
`REF_Sisyphus_Original_Packed` image. No Reimu-named object existed; the only
Fumo-named object was `PLACEHOLDER_FUTURE_APPROVED_FUMO`.

The V01 input, controlling reference, tracked Sisyphus Blend, and tracked
Reimu Blend all matched their frozen SHA-256 values after execution:

- V01: `e482952bff46e3fea6b6d67b90ffc360bada6f45f280d4283f26db647c38e9d0`
- reference: `3d40e2726ae5ff84983f642e20809bb6689c77ecffe5060c6aa760bdee121519`
- tracked Sisyphus: `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`
- tracked Reimu: `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`

Protected V01 boulder, placeholder, reference, camera, and lighting object
signatures were unchanged. This establishes clean reopen, naming, and source
preservation only; it does not establish qualitative module reusability.

# Pinned attempt 017 failure evidence

Repository-pinned Blender 5.2.1 independently clean-opened the exact frozen
baseline and partial candidate and reproduced the writer's measurements.

## Exact inputs

- Baseline SHA-256:
  `6f49bd4e0a8af6b45870d9d4224a520c1398e52e6e9c42f8fb5bee7b8c17118e`.
- Partial candidate SHA-256:
  `2428de0a0b65e572de9437a8d3ef35f1ee21c18bd9dbf27ff01de1816418c0bd`.
- Baseline coordinate digest:
  `41ee23670f67335ac070d95bd782436f53405034f1e24efdebdad709f7d47df2`.
- Partial coordinate digest:
  `f2fdcbdaea90335b9de47861e8ce59b64896ec81853d19ee2ea675725d9fb16e`.
- Topology digest:
  `91e80c085bc52d306949bf2aa62540ee18209d1d4261fa3e57c9c7d55d13d741`.
- Writer receipt SHA-256:
  `4506d8be13155bd2ca2a3434d3080b0eeb998f69749c3c5de88f6e18d6c67fe2`.

The exact input file hashes were unchanged after auditing and rendering.

## Independent audit

Pinned invocation `9ce11acc-6a90-4f1f-92a2-3600ba92ddc7` produced
`failure_audit_017_final.json`, SHA-256
`435e0fb91302363be4fd129eea1123031940e3b5f06a2f6e8dc98eb0aa8be8cb`.

The audit reconstructed the frozen central measurement window and found:

- plane height range `0.45792272686958313` to
  `0.4257808327674866`, a `7.019065055325502%` reduction;
- plane variance `0.014215173048885853` to
  `0.01352948526060966`, a `4.8236330709314545%` reduction;
- 958 changed vertices, all inside `REGION_plane`;
- maximum `CONTROL_plane` displacement exactly `0.0`; and
- the candidate mask unmasked exactly the 1,548 plane-region vertices and
  kept 7,958 others masked.

The unchanged gate is 35 percent plane-variance reduction. The audit therefore
records `plane_gate_passed=false` and `full_coupon_passed=false`.

## Frozen visual packet

Pinned invocation `45bf89ec-9475-4984-816f-03785f563115` wrote a fresh packet
under `pinned_render_packet_017`. Its manifest SHA-256 is
`cbcb3c6e2f8763729e1edb234c9efacbdff528e5965b0c9d113de0a5e34242a5`;
`READY` binds that manifest and has SHA-256
`0c6e618805ef39d22c8e40845fdad0c2f26b1d9c0c59ad1fc01983e356078668`.

Image SHA-256 values are:

- baseline front:
  `54b88efdcaa062737e5f20b98043bdbb0a0befc869aea59eee0174e0682fbded`;
- partial front:
  `97439d62459e7e2f112af20f7ad202c6942e619b617caa8acf7152e00b342f23`;
- baseline three-quarter:
  `0073ff802f667a4330996039f630ca619e860c893dc224a4da215dc7033de8bb`;
  and
- partial three-quarter:
  `2087d1ae2b025d9cf4d37383c96638cadeeee16928aed650bcbf65abace7e713`.

The frozen pixels show a real but subtle central flattening. The root, tip, and
outer silhouette are effectively unchanged, consistent with the measured
failure. Failed binary and image artifacts remain ignored task scratch in this
worktree; this compact receipt is the durable evidence.

## Causal next decision

The positive, isolated response supports an under-dose hypothesis but does not
pass the capability gate. One new plan may hold every setup variable fixed and
vary only cumulative identical Flatten passes, measured in fixed blocks. A
same-dose replay, strength or radius change, path change, fixture redesign, or
threshold change would not be this causal test.

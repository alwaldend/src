# Reimu Fumo evidence manifest

[Back to goal](README.md)

For human-viewable renders, comparison sheets, and reports, use the
[artifact log](artifacts.md). This manifest keeps hashes, measurements, and
command evidence; a hash is not a substitute for an inspectable artifact.

## Evidence manifest

### 2026-08-30 Attempts 58–59 checkpoint

- Attempt 58 candidate A/B/C blend SHA-256 values are
  `36972e42519bf6c0844ac6efbadddaacbc61c435e935d0b8a77617614114995f`,
  `853b0753be5732ddc52f596d61d08f3ae3de71705e35aef192955c4534242964`,
  and
  `f7a1e6664bcb42ca85bf5e4bf3fcb7d907cf9cd69bb7288e7539ce34c1129fdb`.
  Reopen and technical checks pass for all three. The anonymous
  [beauty matrix](../../../../out/reimu_fumo_attempt_058_layered_head/blind_comparison/beauty_matrix.png),
  [form matrix](../../../../out/reimu_fumo_attempt_058_layered_head/blind_comparison/form_matrix.png),
  and [review](../../../../out/reimu_fumo_attempt_058_layered_head/blind_comparison/BLIND_REVIEW.md)
  SHA-256 values are
  `3f40d5278bb64b2062a1560c2a008e502d488e02c8b15e7030aaac6e0440b703`,
  `feee5f8b5a234617aaa40dd40236624a97bdd2f413889c527e9e00bd87b7780f`,
  and
  `31ffca011420463d90eb1a3aa4a6bde000dc9f0b8275efe4a42aba3d9e537751`.
  Hash mapping is `X=B`, `Y=C`, `Z=A`. Reviews rank `X>Z>Y` and `Z>Y>X`,
  respectively, while rejecting all three; only the all-reject verdict is
  robust.

- Attempt 59's [primary contract](../../../../out/reimu_fumo_attempt_059_method_reset/ATTEMPT.md),
  [method audit](../../../../out/reimu_fumo_attempt_059_method_reset/METHOD_AUDIT.md),
  [source-pattern specification](../../../../out/reimu_fumo_attempt_059_method_reset/SOURCE_PATTERN_SPEC.md),
  and [context specification](../../../../out/reimu_fumo_attempt_059_method_reset/CONTEXT_SPEC.md)
  SHA-256 values are
  `48cc51c0e6b652da9292112668589149666ddce8355e7ded63db99fb92128891`,
  `8b23517e07261dbddb2bcc709a8d6249b203524858b4c6ca37fdb58873a5b885`,
  `c0584c0a71e9e59214b27278c99781254ede601de169678790e7682a84c20407`,
  and
  `12a01807eaca8e2f33406c2896ae5f39fc483e017365a4ce33209e7ab0f8d95c`.
  The corrected `Y=±.19 Wh` structural seed has minimum edge/rest ratio
  `.632478`, zero material occurrences below `.60`, minimum triangle-area
  ratio `.645180`, signed volume `.31600305`, and zero non-adjacent
  intersections. Its first cycle also clears every hard guard. The ignored
  kernel diagnostic manifest is
  `17d20fa0064e5e55901135e94d2886af3c94ea2bf87ef1fc2d093961469a943c`;
  it is solver evidence, not a likeness candidate.

### 2026-08-30 Attempts 52–57 checkpoint

- Attempt 52 loose visual-hull guide SHA-256 is
  `dfb0f1cc9517a1055c15eb5efab8bb273388d4c2768327e9d44d37c3c0e4537e`;
  its [source comparison](../../../../out/reimu_fumo_attempt_052_visual_hull/evidence/source_vs_blender_hull.png)
  is
  `9f86a81d053a01484e861bf8daa670fd3a6fd12dabcb7bf1785510e137f2c2d1`.
  The loose/strict hulls reach mean precision/recall `.884/.976` and
  `.945/.947`; the loose worst-view precision/recall are `.812/.909`.
  Strict-core volume is `79.3%` of loose volume. The per-view gate fails, so
  the guide is uncertainty evidence only.

- Attempt 53 head decision, head curve sheet, lower contract resolution,
  lower interface diagram, adversarial verdict, and construction-board
  SHA-256 values are
  `9aeb2cf650575a736e410b7df4d33fad342ccd10d10b141f19dd93b7d5ccd93f`,
  `c530d23189f0a7f44de6071de836f8771c61d925855ada01395feb85b7a03cc1`,
  `6b0579c0206fe5ea3134dab1ec97ec56c5eeb9e7582ecc187acfb000b7b04936`,
  `76f23a85d6085253c1f9526a350718881fbf876e4bb1252df8074ac4bfb24e88`,
  `fa3c0bd31f796c4337de6615ac65b47022cbf6c7d4cc63aac67c9cc6ae20ab0c`,
  and
  `73c472a010686a9d4620e08a3038788d84d3e6a4b1b0cf087fac54e508eb3bc2`.
  The rung-3 source snapshot remains
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
  Its blind absolute scores are `5.5/5.4/3.4/5.8/3.2/3.1/6.3`; the packet
  freezes contracts and selective-copy boundaries, not candidate approval.

- Attempt 54 face configuration B, rejected head-receiver B, and provisional
  hidden lower-receiver B SHA-256 values are
  `9f9aafc5cbfa60372baa95c2db6e27ef4cd7af23493282f93848619f58a38320`,
  `5209c97a7e3770a53bd2d0070e79aa8b52465107b5857ab15493ab4368501ed2`,
  and
  `e4b20c48614fd89eb6cf852170dfbc6855cf42551bac0d1b7f66190e9126b584`.
  The [face sheet](../../../../out/reimu_fumo_attempt_054_face_applique/blind_review_sheet.png)
  is
  `3652e4156fd2c5efb4cdfdc0d3c40ba68a542840c3b548ff1d5488ad443dd385`.
  Face B's maximum nominal projection is `1.35168 mm`, but its blind scores
  remain `6.0–7.5/10`. Head B misses minimum precision at `.852`; lower B
  passes its hidden-support envelope at `.98545` minimum allowed precision,
  a `.758 × .720 Wh` floor span, and three rear-profile direction changes.
  None is a complete or approved model candidate.

- Attempt 55 rejected head Cycle 2 and lower coupon SHA-256 values are
  `37abb814145fb54ff73ddb7e216ec93b1d0a9ff57fef9b070ec9152445de86ae`
  and
  `a8984b2caccd6a95182f6237673b4a74dc174eafe8e473794641a8702b6e2645`.
  Their [head gate](../../../../out/reimu_fumo_attempt_055_constructed_head_receiver/cycle_2/cheap_gate/contact_sheet.png)
  and [lower sheet](../../../../out/reimu_fumo_attempt_055_constructed_lower_receiver/packet/front_side_xray_review.png)
  hashes are
  `6460bc079f12af3b989df92dad8ba8c610e6593d1064317e889fea629da55cab`
  and
  `8de928c420c4e9fab652633a3dabec1e69317b69b4d9f7e7400fcd371f401b0f`.
  Head Cycle 2 fails at profile `1.250 Wh`, additional rear reach `.415 Wh`,
  rear width `.939 Wh`, and the visible helmet/root vetoes. Lower `C0`
  passes `.520 × .460 Wh` contact, `.1873 Wh²` area, and `.7832` box fill;
  visible panels fail at front stretch `5.85%/5.19%` max/p95 and gore stretch
  about `25.57%/8.55%`.

- Attempt 56A rejected State 2, 56B rejected sewn configuration B, and 56C
  rejected State B SHA-256 values are
  `920cae18003df3589f8c8f3cd04abf72f215e3b4618ec06e88416708869eb774`,
  `87a6667439d1f29a74ac4ff494d7b39740f73cebe452476bad72ce07f1226fed`,
  and
  `4c42c1e037ded3b8e00a2e272eb01a777ae1ff31c925598fa38f2f79fc599c1b`.
  Their review-sheet SHA-256 values are
  `b3d2b4c7fb994aec0572742bc1cb2fbc3db85b99b71b7fc6d922f4314968e133`,
  `6fbf302ff095dd6799cf9ba82b701f15c91a6b32928f1de844212cd4acbb20ee`,
  and
  `7a9a2d3a44f78c3d3fb555479c9cb108a5404f5b952c6955438a1e9cb5868076`.
  State 56A-2 reaches `.7503 Wh` depth but scores `2/10` for construction,
  gusset, and lower closure. 56B-B is stable to `.065 mm` RMS but scores only
  `4/10` construction and `5/10` medium. 56C-B reaches `72.9%/65.1%` worst
  pattern strain max/p95 and fails the bench/prow image gate.

- Attempt 57 source board, frozen strategy decision, and maquette-contract
  SHA-256 values are
  `0e4d2543b9f9a735a47c9a60e958c298ac2196058a0496e488c63f27273cc0be`,
  `e72330a188de345f865c784b717c1b803fe6e17d52c22724062b57746492a21b`,
  and
  `462f6116b5f0efb949cf12c2500a8b352f44e7dae83ec205da2afcc96fa374e8`.
  The frozen targets are brown-core depth `.74 ± .05 Wh`, core-plus-leaf
  depth `1.02 ± .06 Wh`, visible beige `.603 ± .03 Wh` in both axes, and
  beige projection at most `.015 Wh`. The only State 1 candidate has SHA-256
  `3acadcc30987872d214ec9bf73cd57dfe28e34935674085ec4458af815a247fd`.
  It passes reopen validation and all frozen numeric bands: core depth `.719
Wh`, complete depth `1.061 Wh`, leaf-owned reach `.342 Wh`, visible beige
  `.603 × .603 Wh`, beige projection `.0092 Wh`, and fringe relief `.0379
Wh`. The preserved
  [review packet](../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/review_packet.png)
  and [blind review](../../../../out/reimu_fumo_attempt_057_hybrid_head_maquette/state_1/BLIND_REVIEW.md)
  reject the helmet/plaque/slab representation; no correction or State 2 was
  made.

- Attempt 58 is frozen in the
  [layered-head contract](../../../../out/reimu_fumo_attempt_058_layered_head/ATTEMPT.md).
  Three isolated representation copies are in progress. No candidate hash,
  pixel verdict, selection, or promotion is yet claimed.

- Attempt 36 rejected W0 wire sheet, cage manifest, machine metrics, and
  deterministic builder SHA-256 values are
  `9ecc6042fe15daa1dc439aa06972ebf38b94211e83d08aeeed70459ffc68ec92`,
  `34824e4e0440be7beeb07105457d56de06416f54899b9f759f8675ba53ca1aa7`,
  `5ad2ea85c3ab85403001c5263feee0009ea3d93bae08259ee07a4f7f2715d967`,
  and
  `687d18bb335ed33d51f987fcd61953e5c3efd7052d764336e22103a009258b23`.
  The disposable cage has the declared `117 V / 212 E / 96 Q`, Euler `1`, one
  `40`-edge boundary, exact hairline controls, one brown and one beige
  component, and no reported 3D self-intersection. It still fails at minimum
  scaled Jacobian `-.097250`, edge-ratio p95/maximum
  `9.503721/15.750680`, front/`-48°`/`+48°` projected crossings `4/52/77`,
  crown-contour RMS/maximum `.116900/.142886 Wh`, and all six crown widths.
  Human review independently vetoes the pointed pinched crown, folded
  three-quarter flow, peaked depth, and wall-like return. The packet completed
  in `6.763 s`; no persistent blend, G0, or G1 was created. The protected
  standalone Reimu remains
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.

- Attempt 35 rejected coupled outer-skin blend, deterministic builder,
  first-gate sheet, G1 beauty, semantic ID, boundary witness, and metrics
  SHA-256 values are
  `ab940623f608d7419c3b72523416de816049e3ee1587141c92d2a4dfbd1e1813`,
  `8bfbc6e2204af00f7908900fa100de65dbe94e6703f081501af2932fa4f35899`,
  `42cfe3b675b126aaf5e1f0182477e7bfaf9bd9f0595229ffe871f371144756bf`,
  `a48d692b9d752d281a313d94b9da9e48d0d9375f8947e541dd098b9aa846718a`,
  `6b9a7424866d924889b19202485b5be48b9631c8016f387be494a08e981ae5f6`,
  `0bfaa291d019e6ca26b2631cc97d14ee72840fe627793ab29af9d1fe6f059ed5`,
  and
  `11aa0249e35a3ba412f647667ea29542ca17795853b7a05c0731f9f6145d5f67`.
  The candidate is one open component with Euler `1`, `32,769` vertices,
  `32,768` faces, one `256`-edge exterior boundary, and no reported interior
  topology defect, but those checks describe the wrong high-valence polar
  representation. The visual gate rejects its banded helmet/cap and broad
  opening: beige exposure is `.802 × .602 Wh` versus the intended roughly
  `.603 × .603 Wh`; the semantic raster reports three brown components.
  Builder execution through four outputs and audits took `15.0902 s`.
  Protected standalone hashes remain Reimu
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`
  and Sisyphus
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.

- Current goal-skill SHA-256 is
  `b954c1678137ac16471b12ee507c87e62a0e6b99aaeb1ae01eeb42ad4a04c5d7`.
  It now prioritizes deliverable evidence, small stable modules, whole-result
  context, early vetoes, dependency-driven concurrency, prompt stale-work
  cancellation, immutable candidate-copy promotion before asset splitting, and
  a verified commit-and-push checkpoint after every active goal turn.
  `bazel_agent build //.agents/skills/goal:skill` and
  `bazel_agent test //:buildifier_test` pass; the latter executes one test.

- Attempt 34 rejected front crown/fringe blend, deterministic builder,
  first-gate sheet, and metrics SHA-256 values are
  `6a564216687bf9955befef44a5279ad6115bbb59c6dea6908c499f475d6a0463`,
  `fb3179a35dea010f88aad6c758b6141e4cfea43693916e02a2165728f9a4d122`,
  `a244cbf2a6eb37053c3580b12aa5fd115177826499c241fd0a5df06f068da3e2`,
  and
  `6f65c0b4b0a9d729475e0cd5d81e538e65a0586982275ec3f51f5f27a215cf6d`.
  The exact free-edge anchors did not rescue the candidate: absolute review
  scores front identity `3/10`, bilateral wrap `1/10`, contact `1/10`, and
  padded-felt read `1/10`. Raw checks report `210` solid no-go pixels,
  `686` guide-halo pixels, hidden-closure exposure `925/843` pixels at
  `+48/-48`, `70` boundary edges, and `124` invalid-incidence edges. The
  center puncture is a guide intersection, not an intentional opening. The
  mesh, Coons-like mapping, faceted guide, and closure layout are rejected;
  measured source controls and fixed cameras remain evidence only.

- Attempt 33 rejected connected visual-cage blend, deterministic builder,
  metrics, build/review notes, and first-pixel contact-sheet SHA-256 values are
  `aac5ee8412eb26953596c1a7eb30517404b94c73e89657ac888032a78738b0ea`,
  `eb681eb2837d781efe0e3b0eb55af5d35574084d6231d4cd7b8333f41ef92be5`,
  `bfcb960281c9aeddf6cd0be6564df979c4092eba2d3846160b4cbbc54af6f1df`,
  `cb79ad01284fcc56a467ace25822d3810ce2383b15ea7473f076d85dc8f42d39`,
  and
  `9fdc89b72f00e7b28320e51e10e92fd5bd0e5306cf6e1d94c46a58664a8fe8df`.
  The candidate matched its declared scalar bands but failed absolutely at
  `1/10` likeness: exposed beige crown, raised W fringe rail, full-width rear
  cape, under-chin rail, and floating pill locks. It is rejection evidence and
  not a geometry parent.

- Corrected canonical head-target and turntable-role dossiers have SHA-256
  values
  `0f5cc5926a5daeed6f09e5abaf9b5781ff4643e19f395f6b2374bda4193b914d`
  and
  `4e789c6acf1b862f8e4c4032c3b341e5f2efa47a5622fc208840e858485ce15b`.
  The controlling front datum is `Wh = 368 ± 4 px`, center `x=485`, crown top
  `y=231`, and central tip `(.588,.677)`; the superseded `395 px` box included
  red/background. The `.603 × .603 Wh` datum is visible beige exposure, not a
  hidden opening. `1.098 Wh` is the lock-inclusive crown-to-lowest-lock
  composite, while the lock-excluded Attempt 36 cage ended at `v=.990`.
  Independent rear-leaf reach requires `.35 ± .05 Wh` of reservation.
  Fixed turn roles are front `03`, three-quarters `07/29`, side brackets
  `10/11` and `25/26`, rear `18`, and rear three-quarters `14/15` and `22`.

- Attempt 29 source-exact chart script, metrics, SVG, and PNG SHA-256 values
  are
  `b96abd732366215f465aa534f0bbfdae2199567626dd48554aa22d1cf0f8fb44`,
  `133ff968a3e023bef509b793e3f533e5f8247a3b91033162091cef81157f4567`,
  `db212543fffca38c7239f0a59f58938ab2ddb85fa9eb1688eb655f48e7bf62cb`,
  and
  `7fd6efca521ca2ffa590074492a4bfde44e2f17ab737b97e114af53ff0dbbb59`.
  The deterministic chart is `5,389 V / 10,602 E / 5,214 Q`, one boundary
  component, Euler `1`, zero non-manifold edges, and exact at all `39` O/F
  source keys. It passes the fixed pre-result gates at minimum scaled Jacobian
  `.193588`, maximum/p95 edge ratios `9.408576/5.788490`, minimum edge
  `.41231 px`, and diagnostic maximum 3D condition `21.37697`. An intermediate
  result at edge ratio `10.0996` was correctly treated as a failure; the final
  sampler changed rather than weakening the `<=10` gate.

- Attempt 29 independent interface-preflight rejection SHA-256 is
  `531adb36fa0c4bf57b355dea09d90ce08ecd7e05f39b0aa60b489982fb8d4fb8`.
  Two implementation-blind reviews independently reject the hidden-`U`
  depth/occlusion contract while retaining the chart. The old successor
  terminal band at support `+0.002..0.006 BU` is not authorized: at the same
  footprint it can enter V4, which occupies approximately support
  `+0.006..0.020 BU`. Promotion now requires an explicit valid layer course,
  complete-shell signed crossing/clearance tests, real occluder context,
  calibrated integrated masks, and evaluated continuity gates.

- Attempt 29 corrected exterior-lap contract, rendered diagnostic, and
  editable diagnostic SHA-256 values are
  `9459ab3e57512bd7fa9bda66842bfdde4262d4271ceaa1dbe0c83718b35e7eb8`,
  `2f0195612fa0e33f5c9d337f2986835b88bce23e73878682d72236508760acf8`,
  and
  `94d3501742c3375e1f7441cc6742c418f8a1d74a6ca316ebbbbff544b6859e49`.
  The BUILDABLE reset freezes support < V4 < successor exterior lap < future
  occluders, uses a first-test `0.060 BU` lap at `0.002..0.008 BU` outside V4,
  and requires complete-shell collision, bounded clearance, continuity, and
  integrated-mask gates. It explicitly cannot promote the complete interface
  until accepted bow/side-band geometry proves zero exposed terminal pixels.

- Attempt 29 rejected exterior-lap Candidate V2 blend, validation JSON,
  validation report, and progress contact sheet SHA-256 values are
  `c76b23e167a169d4462806a988825111d7aef62a8468191ad061d5ad17e6d747`,
  `aef1d27c6626c2dc0016ea6cdde8068b7aca9e3dbde4b08947999952a1be40f1`,
  `2a375a47fd7dbd96ba98c5bb2c589729259089ce99681af86fd3c5e7b1d5bd6a`,
  and
  `6501eb7f283ebd008c470a92908578e865b80bd566f41756d337c2790fcef53e`.
  V2 is one closed all-quad `31,534 V / 31,532 Q` diagnostic shell and its
  terminal V4 clearance is in-band at `.003838..004001 BU`, but it fails with
  `4,078/33` successor-support/successor-V4 triangle overlaps, minimum signed
  support clearance `-.016672 BU`, and O normal mismatch
  `110.70°/92.95°` max/p95. The fixed pixels confirm temple perforations,
  crown slit/ridge, hard seam, and helmet mass. Candidate and parent promotion
  are false; protected hashes remain unchanged.

- Attempt 29 accepted-parent interface contract, actual-geometry diagnostic,
  measured JSON, and preflight scorecard SHA-256 values are
  `e149f6130d53726cfe5e28778030329b3d6db114ec509f1fcf8465735f4c40c4`,
  `c3156703fd1cb990385284210ed014d7533a19aa2c5e91e86b68a5a75d0a2562`,
  `caf4182cefe78da54dd4ccc3569d49141e9e962dead8d77c55189993c561c2d4`,
  and
  `2b84794e82622a1c7b061edfb3c41dc3f54d765fa1245094de5c131e74549be7`.
  Read-only Blender inspection verified both parent hashes unchanged. V4's
  crown back/outward skins are only support `+0.006 BU` and
  `+0.01947..0.02006 BU`, so V4 cannot be pulled inward. A later independent
  review used these measurements to reject, rather than authorize, the
  successor-only support `+0.002..0.006 BU` terminal; the lower
  grazing-left hook is world `-X`, column `84`, rows `17..24`, and remains
  reserved for later lock occlusion.

- Attempt 29 all-relevant-reference board and source-only role matrix SHA-256
  values are
  `9e4a5e9c96214a7d45f22ae895b8c9c9e209fd433409737345cb5c2dbef810ee`
  and
  `0b7bdecaa17dc92e11e63c1100f9eb73cd8bbb4fc51f450a41a7e6ad6c334ed2`.
  All five attachment hashes, all four turn frames, and all eight selected sofa
  frames were verified. The audit assigns explicit per-view authority, proves
  that O00..O16 is a projected silhouette rather than a photographed seam,
  confines hidden overlap to localized bow/ribbon-covered regions, and
  corrects that supporting physical-front variant's center-span/depth ratios
  to `.476/.275 Wh` under its `Wh = 189 px` datum.

- Accepted Attempt 28 rear felt-panel V4 candidate, accepted review sheet,
  implementation-blind verdict, visual metrics, technical validation JSON,
  validation report, and source contract SHA-256 values are
  `d2c090deeffc29f98abd57d60f1e8155d4ad7b5511de2cf4c0de30bde9fe874a`,
  `8f9b5b8b59cf38af11ab9dd3c06f8932d38ca6397ead6915fb716542b10196ce`,
  `4e0a37614d00fc3a2f0d7d220a99eefa64e9e9017b12d8fd32a1dff7b8dea725`,
  `c2b247410b7478611aad68552f38abdaa28f09147064fcc405c66234aeb710a4`,
  `098b316d9a0de04fe871e9dd1eed7dd49c6245900c882a36c6a2cd28b2978cca`,
  `39be062963fcb8f48e3484ffcf5bea0d2950fd330ded643cc16f827d997d92ba`,
  and
  `e4093b94723a86a6835c425079cc99632cb8a00b44a6aad3f348e938f1de98b1`.
  The saved file reopens as exactly `4,898 V / 9,792 E / 4,896 Q`, one
  all-quad closed component with Euler `2`, all nine semantic groups, and no
  modifiers. Registered rear silhouette passes at IoU `.977983716`, contour
  RMS `.0050595 Wh`, contour maximum `.0112695 Wh`, and zero missing contact
  pixels after erosion. Blind scores `9.5/8.8/8.1/8.0/8.1/8.2` accept it as an
  isolated parent while freezing its rear XY contour and retaining the crown-
  tab, front-rim, lower-hook, sharp-edge, and material findings as integration
  conditions.
- Earlier whole-result-guard goal-skill checkpoint SHA-256 is
  `b9701801a94d94fa5dbd629b3912477f8d046c79924301092fc349b993fd6f74`.
  Its module workflow now requires checking the whole requested result,
  composition, shared proportions, neighboring interfaces, downstream use,
  dominant remaining failure, and a lightweight integrated regression where
  available. `quick_validate.py`, the packaged goal-skill Bazel build, and the
  root Buildifier test pass.

- Attempt 27 hair-base V3.4 candidate, validation, comparison sheet,
  source-honest boundary contract, and current deterministic builder SHA-256
  values are
  `30eadec363ac151b76eee482a9111196d59c6e43255740b24829d4de0d3980a4`,
  `6d9d6c2593b4862fbfdc9e675789c62c12062a54382a51f9e2017ea626cda209`,
  `66e53926ea549e66fd11ca0aaccfca3f0ba42e494d6f8b9b7f3547c3cae947d2`,
  `e6e1da4e02e6355a1ff54acff85fded3690fa735aa9aefea7fd8e91bcba65cb1`,
  and
  `77357e15ecd2cb268acece1731e3d9189f58e78859a0ef49fabc4775ca9ef799`.
  The candidate is exactly `10,530 V / 10,528 Q`, one component, Euler `2`,
  with zero boundary, wire, or non-manifold edges, minimum face area
  `.00010135 BU²`, positive volume `1.50937 BU³`, and minimum supported
  inner-skin radial clearance `.0275848 BU`. The frozen cushion, tracked Reimu,
  and tracked Sisyphus hashes remain exactly unchanged. Absolute blind scores
  of `4/4/3/4/5/4/8` reject the rectangular dome, helmet medium, invalid U-root
  architecture, and sawtooth hem despite that technical improvement.
- Rejected V3, V3.1, V3.2, and V3.3 verdict SHA-256 values are
  `0e2d157c6cd853c8a7e7bd088eb85bcab3a7a346a36ab73b8aa3f0ee4a8216d6`,
  `40dd88be38c53dfa0d988c48b78a52a21725a48847feade78a8f1747b2976e96`,
  `37fa932aa2268f9977ed39b7282d837d19fe5cde0fd538d555b7c281252dd76d`,
  and
  `68472d919d78f6af9c0e7b4db6d4887f71e3ef5e566c288a04cec9c39b491a36`.
  They isolate, respectively, collapsed analytic side radius, a widened but
  still floating analytic shell, direct support-wrap interface issues, and a
  convex-chord support crossing. None is an accepted asset.

- Earlier small-module goal-skill checkpoint SHA-256:
  `51357a623c64798a8e5b71e9bbdf4f122003420f7778ce5922e8e33a117ea66d`.
  The skill now defaults to small well-defined modules and explicitly permits
  a turn to return early when inspectable evidence falsifies the active
  approach, module boundary, or interface, provided the rejection artifact and
  reset are surfaced and goal completion is not claimed. `quick_validate.py`
  and the Bazel `//.agents/skills/goal:skill` build both pass.
- Rejected crown-wrap coupon V1 blend, validation, review sheet, reference
  comparison, verdict, and deterministic builder SHA-256 values:
  `2be2215ae78b885270c8f8c432f89f196948dfd238f3a21b2e15a3ce69f1f327`,
  `9ecc8c3fd13502389c2d1b74bc324c0f750203435e251d6afe682b8ec64e18f3`,
  `7e1cf49b7229ffde94db814953f4e3baf8f95b2a8a876425f34f83a0149be6b1`,
  `df7f835563e2b03d943a247887c823820af93ab29d3315fc1c480ea6c25ff535`,
  `646f02999247e3daede74970d40c6e4ecc12180405fceda54998611190d82606`,
  and `54016572eeb02f622dc72666ba91e24285ba8dfb25c249928ba90a2a510288f9`.
  The `1,179`-vertex coupon is one closed manifold with zero boundary, wire, or
  non-manifold edges, no modifiers, thickness `.045–.070 BU`, maximum designed
  clearance `.0549 BU`, and an unchanged cushion signature. Blind pixels
  score thickness `6.5`, contact `4`, apex continuity `6`, three-quarter
  continuity `5.5`, stuffed-fabric construction `3.5`, and presentation `7.5`.
  A continuous crown air channel and rigid rounded-rail read are major visible
  failures, so the narrow strip representation is rejected before expansion.
- Rejected front-hair V4 blend, validation, review sheet, comparison, and
  verdict SHA-256 values:
  `b0731d653b716d8ff6e7a06410145495162182bf1603fb30f23e27593e523c67`,
  `488932be3328d1937516f12851af63c3d05c129fc3fbd967d2c2f7c026a6e371`,
  `d20fdc7497a88b748bdc2a18127459c5ee308afca12964d358220e11a8fe1f19`,
  `280227add4dec19eb80d6bacba20ff1915ce24c1974107040962e796a2f6f7f0`,
  and `f2cbe839b0d2bd4f494eb6294413eeb9a06e6899c669f68f0d93723dddae9219`.
  The explicit `96 × 18` paired mesh is one component with zero boundary,
  wire, non-manifold, or cushion-overlap pairs; minimum sampled clearance is
  `.0276371`, it has no modifiers, and the accepted cushion signature is
  unchanged. Independent visual scores of `4/2/2/3/7` reject the front-only
  offset crescent as a crown-leaking floating visor. This falsifies its module
  interface despite technical success.
- Accepted head-support cushion V3 blend SHA-256:
  `b6b8b84742607d66f01f87362a69f3fa48bcdbb28f5a491192ae4141d4648328`.
  Its [validation report](../../../../out/reimu_fumo_attempt_025_modular_parts/head_support_cushion_v3/validation.json)
  passes one-component closed-manifold topology, zero non-manifold, boundary,
  or wire edges, all named groups and interfaces, and in-band depth
  `2.6450179`. Blind role scores are suitability `8.5`, visible-face softness
  `8`, depth/profile `8.5`, sewn construction `8`, and presentation `8.5`,
  with no major blocker.
- Rejected front-hair V3 [front render](../../../../out/reimu_fumo_attempt_025_modular_parts/front_hair_field_v3/renders/front.png)
  SHA-256: `cabf4f6c523d426a0d6bae2e8748122c20e4903aeafde7191fb23b07305ea690`.
  The visible needle and evaluated bounds far outside the registered panel
  prove that modifier-derived thickness is invalid even after analytic surface
  attachment and bevel removal. This is rejection evidence, not an accepted
  model hash.
- Protected tracked assets remain unchanged: standalone Reimu
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`
  and Sisyphus
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.

- Recovered hair checkpoint 02 candidate, manifest, contact sheet, and verdict
  SHA-256 values:
  `95ef7ac2a14139f29902dbd53970844bac03276feb1b70ee7c54609744e3fd4a`,
  `aac994a9ade0916bba1bac615ec11680615243bf7a1a95eac4a00fd6dfec08d5`,
  `62b3696a9e24170335cb87ee6891251811b423b0e7a5e2a76385b0f82a8889ca`,
  and
  `fc82a5095abfc935a32df6d8bc32e0c33d5e252b3bcb1027e8f248c5a3113609`.
  The exact hashed parent reproduced the candidate after the T3 crash.
  Modeling took `.033 s`, snapshot save/hash `.020705 s`, and four renders
  `12.79977 s`. Front face centers reach `-.0434134 Wh` signed clearance,
  rear centers `-.002454 Wh`, and two components have no supported vertices;
  construction and visual gates reject it.
- Physical-side fall-calibration image, measurements, and parent-review
  SHA-256 values:
  `fcf635e9ed3e8b1d1e4c0b9749183c54568b7c44e5f2ed9b7148ff705f16695e`,
  `fcf1f9f2a4fc5a575d03ffbab7263492bb4d2ee40b428500164def8aa00bf563`,
  and
  `41c9ba7cdb2d93e5a304b4a5b12c4de1a97b4b9cca6ade177102f7a6fc562597`.
  The source trace measures `27.8/147 = .1891` vertical fall at the first
  quarter of `216 px` rear depth. The invented `.60` gate is retired; the
  registered all-column source roof plus `4 px` is authoritative.
- Pattern-topology preflight report, source board, flat-cut diagram, sleeve and
  first-mask gate, and machine contract SHA-256 values:
  `62e88da4833d70d96df1da4cb820266e497d2ac88f3bfda50cfb8d3eaa4f529b`,
  `60bbe0a156cb1d22878382e62234d7f2a8f1f4ca2406e40a30d287c4ae52bcf7`,
  `0a77376ebbbaefabaa5e6fc075e2714e87d6577db550088fc609c46b84a724be`,
  `73f759ae70be189cade8be92457de711cd3d8228382c274ea3a590b4a09e28d5`,
  and
  `6abaec1348e17c07292efda9c60ba9f01c3f1242c68b5ace1f2a4544d15a2c21`.
  Parent review accepts its seam graph, waist-only pins, and free-hem
  topology; the provisional fall ratio must be calibrated from the physical-
  side source trace before it becomes a gate.
- Superseded live-checkpoint smoke candidate, live manifest, verdict template,
  and pinned-reopen manifest SHA-256 values (retained as **invalid negative
  evidence**, not an accepted tooling smoke):
  `5a002771b32be1620211a8159ba9c4893199168fc6ff0dcf9fb566a3a0039d27`,
  `5978edc5559c03835dac6f007ab8e742330ce451d960901e723d573ad7202a99`,
  `5390a227c191f9c4160e67a7396d4be91fdd771524e54f391720f8ab09790a7c`,
  and
  `8e18307172ef196d0865ef7bad9f6720fcf1d2e1358b5f74698168c23dc9e585`.
  Snapshot and hash took `.046471 s`; four 64 px live renders completed in
  `13.134157 s`. Independent audit found that `context.source_filepath` names
  an unrelated cloth-sack file while the manifest claims the hair parent.
  Reopen success therefore does not prove ancestry and this packet is invalid.
- Authoritative split checkpoint candidate, request, READY marker, manifest,
  verdict template, front/rear beauty, and front/rear silhouette SHA-256
  values:
  `b17ea1afec125dcefe262945b17726a4bbcae6a5bb96670d716f5ca68e038d35`,
  `bf9fbc7b52f0bb1e5b0dde31d19b8da2c7b0b18b681db7307e17e79ae6367a68`,
  `b50273906ea063af001ca00e0246038bb7d1b27efc75d8171a11d5d7b90bd6da`,
  `b6fa270f0ff78846dd79e68d2003703be77433412ce22f6e9e80b1a9740731cd`,
  `7ee9fb28181a646f9fca488cdb2297a01b62ac6c04d0a2acd401a3a33afadb77`,
  `e0d35764a8eed9229ff729dd73ac0d9e987a7d48dac68f46d7237b3a05b3dc57`,
  `84600ca86b9744ecf261d21ccc5e369686cc682f6b56fd82e38536c5a288a5d4`,
  `796b807273916081d51e578800a5c1696f10df171d6c905cc97bb0b58b78f1b3`,
  and
  `186eacf681554c902700c3c9dbb0adeb10bc57352dc81c223d73ceaf01f250fd`.
  Blender 5.1.1 controlled-opened the exact `52447e...` parent and returned
  the durable snapshot in `.047686 s`; candidate save plus hash/reverification
  took `.024597 s`. Blender 5.2.1 then reopened scripts-disabled, reverified
  every frozen input, rendered four 320 px files, and published `READY` last;
  render work through manifest inputs took `8.534518 s`. This is a tooling
  smoke around a known rejected model parent, not model acceptance.
- Constructed-panel hair candidate, contact sheet, reference comparison, trace
  alignment, verdict, metrics, manifest, and cycle report SHA-256 values:
  `2b5d74bdfaeb309a0dd8a89c1661ad605c919fea942570ab1cda93f41a3b73d0`,
  `a42c28fe03847bb6d5732d50e254d7910a50709d3564bd8984e5f0034b477fb4`,
  `1c66b3313660ad07733ac79a7b2a6649abff343f940aad47ce773665d348d43f`,
  `bbfe39c21bf237b13bbd4454cde7c75778b817a54ed9bd00949d0b0631eb6844`,
  `06a972bf040c1403c85709e72686b8a435231619d9df70605cc04e2d03b5441e`,
  `a4a9d1747b2972b12ea775c48225da408d5bc63ff8462c0943c0e6f0525e5193`,
  `d0d4af4b30267d83239dad45041938ad7db2ae870a30e602e51eb4e1f6ef59d5`,
  and
  `a2a77cc615489f0d59896ae67f65f2ad31c449d87a8b1ca04857dd91762c2ec3`.
  The cap contains `6,912` all-quad faces; registered front/rear edge error is
  `0.0 Wh`; minimum evaluated clearances are `.0306785`, `.0076700`,
  `.0077496`, `.0168446`, and `.0185450 Wh`; and all `14,425` rear-facing
  support samples are covered. Those technical passes cannot rescue the
  visible box/equator/V/curtain/teeth/fused-lock failure, so the candidate and
  rectangular-field representation are rejected before side views.
- Front source-reset README, schema-2 contours, source board, comparison board,
  physical overlay, clean-front cross-check, rejected-candidate overlay, and
  integrity-manifest SHA-256 values:
  `86bcccec3a15d538c1828d17251043e2f2afd54839c628401ca5de80d0e90191`,
  `bce01ccb0a942f66abbcb839d5fe9c878919c08946cf6b7aa54e2723e2747639`,
  `e8cca64c3ca0272dd37496983f017df3258ed3c7fb682ebd1a3b031d097ec17a`,
  `ac6687198897d670f5c96c2d4fab8436d3a43ca233b67b597a4ccabf12e59d20`,
  `990c204feff5582b95a030453ecc9aaffddcd18f9d5cb98b742384249fd1f1ab`,
  `ed5ac7f344bb278f37e308beb8a9eabb8eb291f8b766f4d6ee8c638991c8b720`,
  `85d26218a840de32d96550a9282b5fe73aa6354d7e31678fe2adc38ebbaa0c66`,
  and
  `8f985432e0ef9be5915fcca71651456d27dfa25924d2aae13ae927038fdee096`.
  All `24/24` listed artifacts verify; maximum stored normalized-coordinate
  error is `.000004973545`. Its prose now pins the clean tracked
  `LANDMARKS.md` hash
  `c7ea9fd19f077d6dce2055275d4c9b1dbfe1b9a029cd085e09c5ae1dda8b7ab7`;
  the earlier stale embedded hash was rejected before geometry. The packet
  accepts brown visible ownership, one
  continuous asymmetric three-span front edge, the broad `(235, 211)` center
  band, and full rear coverage only. Literal one-piece brown-cushion,
  beige-applique, and occluded seam/root claims are rejected as unproven.
- Conditional sewn-panel preflight README, editable SVG, rendered diagram,
  and integrity-manifest SHA-256 values:
  `b40b84fccb97d5eac4c1fdafeded13d927ac22f233dcb4990cb2eea8a139645c`,
  `6c413dc3aa352625c2ff0be4cca5fafd40d78788545b6962e802a85a656f7d3c`,
  `9829f63e5ecefa049d2500c15708c5fa8d807dc0f20c1aeb087a18d175ae59b9`,
  and
  `49d6b202256f8d1cf6c9f23e879a15fd669ccc4a9d2f44158582f4256b52ffd4`.
  All `14/14` entries verify, including the exact schema-2 contours. The
  `A/F/R` brown-annulus, nearly flush visible-face insert, and brown multi-patch
  rear-disk topology is authorized only as one reversible front/rear macro
  trial; every hidden closure remains an explicit inference.
- Pre-pixel coupled-head absolute review contract SHA-256:
  `673e7a1c744e1d635e30c787cbaab30c63782afc4c2d173f4e010e1f521cd3f9`.
  It was corrected against the final source audit before candidate geometry
  and cannot be weakened after viewing pixels.
- Flat-pattern gravity candidate, true-side beauty, silhouette, source overlay,
  verdict, metrics, technical gates, clean-reopen report, manifest, and cycle
  report SHA-256 values:
  `aa81bd5932dcab3baee96acc0ad51c258c138f0a729f959af7e6b23e6b6db1fa`,
  `7ee3cf4fe9731d9644a8d7ad599d89a1cce4b2be074fb50d45490b1eeda80112`,
  `27bd335eda48078224c777a111bd490e744cf6b83f0ebd61d1bec1a18bf98dfd`,
  `bdff0b0fd36fd0555d2fc883d4ce0c06d3d006d5cb2d03b6d69c272cd458559b`,
  `0330fc646c21e21f3fe0cc19811cc7c69077e0df86a9ff2ea5825dfe4accea16`,
  `71664250bdf844670a0acc47fd304e5d318262d37fc355ce0e3a5414a948d343`,
  `7c54c7b245529735b192063e6052132adf3417a9e3e8346058ec400a09358b6d`,
  `ab0aa99a40d7ae60d4cf47ff412081bb7398865a5af0bf7770f9bf5716daa22a`,
  `470c9d12826750cd573132399770b0870162f7a6cee01c1d9c04848129634712`,
  and
  `f8ddd52a2940b974e98952f46afb7b5340d8b6fff7249b9825514cc34a2672c3`.
  Planarity, waist-only pins, face-free seams, open hem, convergence,
  support/ground contact, frozen ancestry, and clean reopen all passed. The
  independent first-pixel scores were `3/10` silhouette, `3/10` construction,
  and `2/10` intended medium, rejecting the short lampshade silhouette and the
  fully waist-pinned homogeneous gravity-only representation.
- Tracked standalone structure-audit report and JSON SHA-256 values:
  `c053eeadae47bba8620e13c53834d95869cedc0530baf9791b38c966b8955aad`
  and
  `ca6278c67cc422bb6da2b4f50d98867546d60cdff93f272d2b3597a2ea9fdda7`.
  The unchanged tracked blend is
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
  Its `FUMO` collection clean-appends with all IDs local, but it has zero real
  rig/animation data and `35` evaluated primary-form crossing pairs; reuse is
  limited to packaging, hierarchy, scale, naming, and temporary neutral-review
  conventions.
- Calibrated sector-skirt candidate, metrics, beauty, overlay, and mask SHA-256
  values:
  `8389e932c0f111950864b9be154d7dd3f50c31e6edc71579f5a61bb8714971fa`,
  `29147d9082e1f90753b22a657c33c5b1518b8ca0ae393101607b4bfeab9acdb6`,
  `6f7c4054e435c371e8fbbdc83a7706d899a0b43b7b052d8d296884bb8329e0d6`,
  `77662c57b27a9d15d8130122f97f064ab22647bf78939d3cead4cc92a11ec71a`,
  and
  `925b698f522e425d7573869a166e30650a8bd0d406d1c8505fe5ffe168f7088f`.
  Pinned Blender 5.2.1 cleanly reopened the exact candidate. IoU `.85987`,
  maximum roof excess `4 px`, zero violating columns, and zero ground gap all
  pass, but the controlling beauty remains a visibly taut diagonal cape.
- Support-derived hair candidate, contact sheet, reference comparison, trace
  alignment, verdict, and manifest SHA-256 values:
  `6c27030778fabddfb3936d7e51139b41709329e74f5a2a29a0ae31e6b9d38dc8`,
  `373a5bda1c3e9bdc4f9d481a8730b27f23335a0778d63d85814e706c994d5717`,
  `2702969d2942cbcc145e5567ce99c1069ad7cba183440b14a37dd791421770ef`,
  `373f16d0cbe3e377dc97eb3c5d5a6419f5c12ee37ceac35eab091e8cddebee8a`,
  `3682f20492afc9c654044ec2cc65d1e87ba0e27d9b3aa04608730917902b9c5f`,
  and
  `377bfba4c6760bb2dfb3d0b3afe2c4b04ae3efd0b53a172f043a822793eb8e56`.
  The shell has one component/one lower loop and minimum evaluated clearances
  `.0088761 Wh` at shell, `.1173574 Wh` at left lock, and `.0265978 Wh` at
  right lock; front trace delta is `.000022 Wh`. Pixels still reject it as a
  helmet with pinched sails, floating locks, centered V, and narrow rear bundle.
- Approved hair-trace report, front trace, front source-and-mask comparison,
  and rear trace SHA-256 values:
  `22dfb00362a505d7d59bfcb647a3e5207e812a7ca574020db7b623f3241218e8`,
  `b00026b5efd792b3b15d3cbdce71907bfd46f91f81b7a44f6f6b85a39e231e0d`,
  `2c341fe1a47fe09c39dd5cb0e4d39210bac6e120a83ad3b867e5a4e3024d10db`,
  and
  `8fb5d3c2ece1d8d896a233c2746e9b41760bfed38160d22f2e74eae059a1377d`.
  The registered front root span is `.506 Wh`, center depth `.292 Wh`, and
  lock widths `.146/.140 Wh`; rear crown-to-lowest over width is `1.254` with
  four-to-five unequal lobes. These paths control the next 3D mapping.
- Gathered-cloth checkpoint B candidate, metrics, and contact-sheet SHA-256
  values:
  `0508f29ff2d563a96f7da3e8e9268ed6d15d6799cce92c3408697659346c349c`,
  `81e6aa5a8fb0ace38a76ceacfee969f9d44686d066fe9e4d8791770754ebe19f`,
  and
  `18fbab442ba46e36acf623c24b37373923fec82ce263d6a3c6ea00d41d564e33`.
  The actual `21 × 19` simulated sheet gathers and its centerline pools, but
  off-center columns keep a cape; simulation took `2.490 s` and the complete
  six-view checkpoint `11.146 s`.
- Macro-drape checkpoint C candidate, metrics, and side-sheet SHA-256 values:
  `a955f8b9207c8f5e2303ac863997e3533810cdc69fbace4f15353ce0c9855763`,
  `fce0ad2a8b4de23f209389e53fb15bb4dba3f3f2027d39a5a9d3f3b2dfcac267`,
  and
  `d44fafa2d3b57c1fe4ad3053695197e8bf1d701476f94e707a88c6ac05c89730`.
  Saved and reopened candidate hashes agree. The side-only gate stopped in
  `4.20 s`; `12` columns exceed the source roof by up to `15 px`, so the
  rectangular rear-sheet family is retired despite its `.6523` silhouette
  IoU and `30.48 px` center direction-change margin.
- Whole-process iteration-audit report SHA-256:
  `3e0518b158463806f3c1596d35aa78af7ddcb8b6e45cf4144372bf1f206e749f`.
  Its direct benchmarks measure warm MCP `.253 s`, pinned beauty plus ID
  `5.47 s`, warm Bazel audit `2.50 s`, scoped validation `6.88 s`, local Git
  checks `.05 s`, and remote read latency `1.65 s`; source interpretation and
  handoff remain the observed critical path.
- Hair Stage 3C checkpoint 01 metrics, front, rear, front silhouette, rear
  silhouette, and front hair-only SHA-256 values:
  `61df3a85f64af4894ff2f515d5b48df717c4c1d847aa47ebb5e95e3752ae3659`,
  `63ade4b5db165d021e90332cce73dee502f9b5680e40d3df2d487f08a1fe8579`,
  `b1eab952d1933d3e1f8ed3cbf8c96488aa530c745fc0e11c89e18c3049f36596`,
  `b6946481017463bdd812c2a4a83771e7ef87aedbf9ae8ec9a905467642db1f64`,
  `22f3214e7c6e5a3ee3656688444c7fcb70970206a92af07a33f6850869251dcd`,
  and
  `443e9175d7330570af514aa855dc7fdcce7611753f600ef2d0b563a541480cc4`.
  The rear ratio passes at `1.31630` and the new outline improves, while
  support exposure and the front-field visual veto reject the scratch state
  before a blend save.
- Source-owned garment checkpoint A Blender file, contact sheet, and metrics
  SHA-256 values:
  `5bcce547df55de47e17bab0fffc128165e8105520c09b34024abf57e6540c247`,
  `c7c5f25be07df36c87fc3cabb7c758e2cae99d9c264423eba4bd2fc0a3c18e40`,
  and
  `9c5b1c14d55562b2fa7970fcdf057ad98b851d7d079595bd71948cb84d414548`.
  It created six 320 px views and reopened in `8.63 s`, then failed the true-
  side cape veto before masks or a full packet.
- Hair Stage 3B checkpoint 02 front, rear, front-silhouette, and rear-
  silhouette SHA-256 values:
  `5b0a740fa884a1a7f9832c64739bff89c38b993ee76bded4549d4a56c3b2ae92`,
  `61a1a7c69c60a05dc9543ae52c4cfefb7095b3e2b5a06cf901954345eee57633`,
  `9abcf0db7f7bf25362371947884cd8e915ae555a6e11e037f92d62b77045a740`,
  and
  `97a4c50352931d8aa01bb6d93db2cb8520caf7a8e915f91da5123c9dd5576d1a`.
  Rear width is `1.06252`, crown-to-lowest is `1.38560`, and their ratio
  passes at `1.30407`. Construction took `.01685 s`, beauty renders
  `8.12112 s`, silhouettes `6.04747 s`, and the checkpoint `14.20053 s`; the
  visual veto rejected it before a candidate save.
- Hair Stage 3B checkpoint 01 metrics, front, rear, front-silhouette, and rear-
  silhouette SHA-256 values:
  `24013edb24564abc0ee11478bf5eec7477b94eca97443e7bccd9f8ef5449468d`,
  `668cf546973709b81c3c0c2d374f17ec986df4e3ce5f6b2c18855cfb7941d981`,
  `da0e5562874058734014aec6d3c5033400e09ff13c85b8e44e7aef85650bb525`,
  `297409d4027d115f3fdcc570be5ee1b9257fd2a0a8b198a64fc88795fe4dfe12`,
  and
  `2c2957be19bca74ecb71a16f1f1487cba9a924eac5251fbf42f8a3f918d4b9cc`.
  The rear ratio passes at `1.2993`; crown exposure, curtain sides, and the
  incorrect front field fail the controlling image gate. No blend was saved.
- Independent hair audit report and comparison-board SHA-256 values:
  `d33aa63f8f0f430b1eff21b63f24bec31b029156a360b1b292f43542bde0bb72`
  and
  `94185d13626ea3dcd5402878228ab33bd274ae7bdcf7e64483758422653ae150`.
  These are source-owned observations from every supplied view, not
  measurements fit to a candidate.
- Hair Stage 3A Blender file, contact sheet, and deterministic review-result
  SHA-256 values:
  `52447e6982733851d099de8eeff476385fb29aa9d4088d3bf596a4e32bbd274a`,
  `d11de5d802e0c4f08c16f1773c145934a19f17775737d7897038805d076efd8c`,
  and
  `ed61533daee4727e32d52f0294b98deecbde0002bbb6408c0ee2fddf0f5a6413`.
  Structural ancestry and clearance pass; absolute hair likeness fails.
- Source-owned body preflight report, construction sheet, front annotation,
  and side annotation SHA-256 values:
  `45bf0b11211ae1ffc399afe36b590cbc2b1fe37ca7d68619a33bf7935e0cebab`,
  `cbe1957791909041d8576ba9a5ddb8273929f1b87942bccb9e730e6640a73919`,
  `55e861dbadf03a5e7c257115ef0ec5517bdb5cdf336c03aea9b364da4ec9c19b`,
  and
  `37245d00b370d4812583d761e7ced79af6d64799d34b89b3474a46e23ae82a12`.
  The images support the compact-core plus hanging/pooling garment ownership
  hypothesis; no Blender candidate was created during preflight.
- Seated-core sculpt Attempts A/B Blender and checkpoint-sheet SHA-256 values:
  `9f1c8b0876737a58693f11a39321e87555a4ee83101192521cce7ee2d6e1ab33`,
  `c8b180677fb44214ee9458a492ae504b3618640e60e154d97505dca75706b2b5`,
  `61233b789417660907d0dd2779e1e004fc12b00d8d6056d26176ce3b2dfa62c7`,
  and
  `1fe9ba8cc50c5968309f5e00a358872b730f84bec9cb475d79239c7d52d5fa78`.
  A is one remeshed but visibly stacked volume; B is one continuous cage but
  a visibly tapered/reclining mass. Both fail before garment derivation.
- Seated-body assembly Attempt B Blender file and two-view checkpoint SHA-256
  values:
  `e96599314e85ab44a353f010d6762030c450b33b95b3425fa024024727b48410`
  and
  `210a236dd58e1df6fdec206120e9d370257c5d3317a34ad4deacdf06264a96d9`.
  The low support exists, but the visible analytic panels ignore its form; the
  cape/collar/wedge family is rejected before diagnostics.
- Pinned-Blender packet smoke manifest, beauty, and component-ID SHA-256
  values:
  `d230851c1be176f212280e5332abe413b8ac208de350ca788808c7c003fca6b0`,
  `d200855b1064a4b94433bede2012b2f8262d03d67bb67de62376c70c28ac6336`,
  and
  `f30295c8d07ca658a64a1c97688970b633a35c1cae57f07f086ef781338c6ec9`.
  The manifest reports Blender 5.2.1 LTS and hashes both the source blend and
  render spec. A repeat produced byte-identical PNGs and canonical manifest;
  the ID image contains `104` exact `(128, 64, 32)` interior pixels. A separate
  root rerun produced two valid outputs in `5.82 s` end to end.
- Seated-body assembly Attempt A Blender file and two-view checkpoint SHA-256
  values:
  `a13903a27330a8e227ada67bb90dfa01806c36cc98368a39be09d1a6b6bd94f7`
  and
  `64d252484a471af960e1d1979ff47c2ae1b47ffd08e77f4800888d98cee063fc`.
  The front has grounded joined feet, while the true side rejects the cape,
  exposed collar neck, and pill sleeve before later packet work.
- Isolated red-mass Attempt F Blender, front, side, and three-quarter SHA-256
  values:
  `c2a2374887218b8d25fc211522ead50d7f9197b5da59633f0bd92dfc9cba01a2`,
  `931ab9106f633c554a8489fea1c7e1f0a087ec28f77a1cc02b8e62700e2af312`,
  `fa2a0bcc064893cf6f7cd03f0923b766cce9b3da3d2a95e193bb08a3494b950b`,
  and
  `8451f26f77a695473b92fd95c3e510ea3bd4c463ce771b74b2750977ec7e3583`.
  The open hem fixes the prior disc but remains a cape with detached shoes.
- Corrected coupled head-context Blender file and contact-sheet SHA-256 values:
  `30604f00ca959629ab79f3a0a8cda9235d46196bc420e06ef3fb339dd01feeeb`
  and
  `9278754dc2628484081b8c6a694ef455b2332cbae924962229c730aa1f06d35b`.
  The support remains a provisional improvement; the diagnostic hair is a
  rejected helmet/cardboard witness.
- Two-panel mochi checkpoint 01 Blender, front, side, three-quarter, and top
  SHA-256 values:
  `439ee7b03de3c441f9fbc06d674b1c15d60fe68f490acf58db092777da5f764a`,
  `01650f361bc99ba2cdceeeae4578e7ff75eb75de2b8b11e62c890e7e91cfb6f9`,
  `fa00ed668924ae674219c912262ec6b06969e967d97d957108b222f5266afd31`,
  `edb70105cebbf0ca79c2583767b378948dff16ce728598d7ad7a00cfbf5108c6`,
  and
  `4b1a730dd27b587698eba73f831033916e1f36447d25f2f69fa56f33f64cee3d`.
  Its `H/W=.89024`, `D/W=.77511`, and no-broad-gusset structure pass as a
  provisional construction parent; the exposed seam still fails this state.
- Retired Stage-2C checkpoint 02 Blender, side, three-quarter, and top SHA-256
  values:
  `3da55940a4f745bf59815475dae28224054941188a9ab11d4028002fda5f804e`,
  `44859f7f7e21607d28f19d9080879421400372bbc6103787ac03bc3ece234e81`,
  `08dcd01c7d295fae09408f09fe8783768128daaa177eb5f212f11e5f3c948ca7`,
  and
  `6f9092a58f052199a48be72197e4c4270089b1a1f4b4a434f9f25ff46bfc15ed`.
  Numeric D-profile gates pass, but its foam-block image gate and ownership of
  the outer-hair height fail.
- Known-bad body Attempt B Bazel audit HTML and JSON SHA-256 values:
  `cf89cac378cfc7485fa731a27dfbc50ec97f8379780868fda606f500ed102585`
  and
  `2ed6c05dda13f8384f0aaa50b3f0d58a1d52effe61ac9ee7f5e8bf9a1cde8409`.
  The packet has no file errors and fails the intended head, bodice, sleeve,
  and foot image gates. This validates rejection mechanics, not likeness.
- Stuffed-cushion v2 Blender file and contact-sheet SHA-256 values:
  `76210602f4b76bdadbe87d0c17b727acd59081d7f9177c741dcb947a0cc3c892`
  and
  `55b3cece888dc59e7629cf84dfc7347d801902f11269421b1f4a390cf932eb63`.
  Exact-parent, topology, envelope, crown-arch, face-plane, and seam-tension
  gates pass, while absolute top/side/three-quarter review rejects the form.
- Non-bald review proxy v3 Blender file and contact-sheet SHA-256 values:
  `e30e19166432bde63d5368aae43adac1436fea2879361c9a4d95ea869c363072`
  and
  `d5e6fc19a88df589a6923cca732908504acaf5b92aa9a550c335cf9ba36b3826`.
  The rear and crown are covered, but this is a rejected composite proxy, not
  a promoted hair or body checkpoint.
- Paper-pattern head cage Blender file and contact-sheet SHA-256 values:
  `652bdfe58f2e9766b6d49503c92f015e4a555fe8caedeb748afc4c1735f4f7df`
  and
  `68814b3ae9917d5ddd46447f48c0459ec229418eea5b3813a57981ada0de03b0`.
  The bounded planar front/rear/gusset representation test passed, but the
  angular cage is not an accepted plush surface.
- Stuffed-cushion v1 Blender file and contact-sheet SHA-256 values:
  `467da28e44ab3f83e672449ac6c340336835220feaa9226e0023a887fe01cd26`
  and
  `95aaeab40332880c529fac34473ceb984df2e0f675dc5063e07eca97e8c04b30`.
  Every frozen topology and envelope gate passed, while absolute review
  rejected the mattress-like side and three-quarter read.
- Body v5 Attempt F Blender file and contact-sheet SHA-256 values:
  `310488600ce325ff03dde175b0018bf8b95f170c57c028e76fdb0f18369f7ad7`
  and
  `3719976fddac8f4ab905814bc3747952eca095132fd350f47b02dcd2db56e343`.
  Crossings, boxes, angles, topology, and median foot occlusion passed. Foot
  top maximum and the controlling absolute construction image gate failed.
- Goal-directory migration manifests:
  `out/reimu_goal_split_manifest.json` and
  `out/reimu_attempt_split_manifest.json`. The tracked monolith had 3,751
  lines and SHA-256
  `9ee3913058f322cba31e37533fe3cb6362b7051b5a130a080cab6fbfdee095bb`.
  The migration preserved all 17 top-level sections and all 20 numbered
  attempt entries, produced the required 96-line README control plane and
  purpose-specific records, and split history into one linked file per
  attempt. Relative-link and content-preservation audits passed.
- Goal-skill directory and session-artifact guidance SHA-256:
  `664c2e9602cbb5ed96b2e958b9553afae2c71c80ecc8832917feb8123a343315`.
  It now also requires whole-process critical-path audits, progressive-fidelity
  speedups, and proof that each workflow change reduced time or rework without
  weakening acceptance. `quick_validate.py .agents/skills/goal` passed; the
  host-bot Bazel rc now supplies the installed Android NDK path.
- Body v5 Attempt A candidate, result, and contact-sheet SHA-256 values:
  `a379eccf328cb8d6bae1eb8a43598505dadac83cdd5f9a328e5390c984d5d7a6`,
  `7c8caeab5d3fb4fa8af45d6d5f8e65b9bf79ec4362168db78a8ea3157394a846`,
  and
  `c8d636ab063e87d203d062f2561f83eea07b335a9de5d2d30d8be3e45fcf6485`.
  Its box/topology gates do not override the rejected hard assembled read.
- Body v5 Attempt B candidate, result, and contact-sheet SHA-256 values:
  `6905d61643619f0cacd6fe6a4daab2e21e6ab34493e6e357f2cf22b1e8a49bc8`,
  `dd7373325c5b862e234e1db1d63916aa75ddeb442139b87d4cab28665a4ee230`,
  and
  `9ce0711d7033c8c476808f21f80e6854804a5b3880d9d26ba9f2a8de3486edf6`.
  The visible review and diagnostics both reject it: `32/32` ruffle/foot
  crossings, failed foot mask and size gates, and absolute category scores of
  `4.5/2/2/2`.
- Body v5 Attempt C candidate, result, and contact-sheet SHA-256 values:
  `98fa1f1f52737cb089411e3896b782c47e766b1cefb18f8cae1b619488979cc3`,
  `3fbb2bab7f5dc67217a5f4ef48bf3a18347370003156030c77225e06a08bff17`,
  and
  `3cbe605f21ad45410ef82222404c63f8d6bc5fdc2363f29d440e77e35824b0f4`.
  Every box passed, but visual construction failed and crossings regressed to
  `806` skirt/ruffle and `150/166` ruffle/foot samples.
- Independent goal-skill forward-test final audit SHA-256:
  `e08805e622753df1c3f8fae82ab32c4db1af035aeb7581b5b7af9fc3bd364c07`.
  The realistic six-stage fixture passed goal layout, record separation,
  rejected-artifact preservation, exact final-byte verification, and manual
  session-link visibility. Its nonblocking findings were incorporated into
  the current skill hash
  `664c2e9602cbb5ed96b2e958b9553afae2c71c80ecc8832917feb8123a343315`.
- Head/hair representation-reset report SHA-256:
  `02c202e6fa5640900cdf431c7db66840eb4fa64307c209305b4dada2e25be6b0`.
  It replaces the sphere/egg family with planar front/rear pattern grids and a
  variable-width continuous gusset, with hair withheld until the cushion gate.
- Face v4 candidate and five-view contact-sheet SHA-256 values:
  `99e63b7fc8c398910f04e6a8eaf41e00f9cfaa8a6d66e5f34ce2c9b4f68d68f2`
  and
  `26373fe49d3d6294407f7314453219635aaef3402820ebe4d7b050e5bcb760b3`.
  This is a preserved visual rejection, not an accepted checkpoint.
- Continuous-shell body v4 candidate, result, and latest contact-sheet SHA-256
  values:
  `d7752976db29d79bb2c43f1ea526d5ae649dded74939a69a1aa15faf9a5fc946`,
  `aceef25b1029b783b54a56acbd4741348bef6a7a0145b294e3308014800a2c3e`,
  and
  `e7a5eb3e804675dc2a11c2766160d210672e7d326dfa9e337a376a863752e95e`.
  Its foot/garment BVH overlaps remain `112/114`, so it is rejected.

- Migrated rejected baseline SHA-256:
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- Baseline five-view contact sheet SHA-256:
  `cf044a81699e1833ec5ba0e233814c02c8991214f7427fef6f995eb94192cd87`.
- Rejected v31 candidate SHA-256:
  `9d0ac214676d0ca61e971c26998da4698b1d7645f506cbdb05b946e9dba41839`.
- v31 render packet hashes:
  front `fa3efc0104da47315cb23fa49a9788b8f06ab645d314805f589c3bc036e7bca4`,
  side `4782d53ab58283b93b6385f6601c42bb41384bfc4ac3ba516d14962451103e77`,
  rear `e0c02087c5a4a318f4856b4a3639ea96b82e1b55842d695cf84fa110c76d4187`,
  three-quarter
  `164beb8483350e13ae4ea7aad3301894e0acd832692826d8820da58627df3685`,
  and presentation
  `002f760fd6087a6e85c632d1d907d7b356df83c41365f7c1b7e03bf257a1ed83`.
- The five raw reference hashes are recorded in the reference dossier.
- Frozen `LANDMARKS.md` SHA-256:
  `c7ea9fd19f077d6dce2055275d4c9b1dbfe1b9a029cd085e09c5ae1dda8b7ab7`.
  The ownership clarification preserves the outer `1.03 Wh` envelope while
  preventing its reuse as the hidden beige support height.
- Local reference-measurement dossier SHA-256:
  `f6d303c0c6c11befc0da286931e852e3c5c1906f326014ddf0eba955d222f0a1`.
- Turn and sofa contact-sheet SHA-256 values:
  `736884a1866aafe890f6c4bf748b1d1a3c2bef8cba25ef47e3bbcf91ba718500`
  and
  `e46c85843e7cd2314f66dcdcd374413233806f068d021834d06032b8ca106119`.
- Rejected Attempt 3 candidate SHA-256:
  `9aa12e0563d3ba71f1ecbb2181ebf8dc90c55c890fa333a82114b9a4cc40dbff`.
- Attempt 3 metric packet SHA-256:
  `752368bbc80a6e606bec977c91f5db5b09ec0f5674b28354a609a74b05fe5e3a`.
- Attempt 3 clay render hashes: front
  `3e34851c1623fb4a30a473f94fcfe05204542eee7ecdafaaa287ca85804ae067`,
  side `d9baaa8ea0274850c1056d8c6465f33d48bd3237f685a03d006014d8ebe4eb0b`,
  rear `d2ebcba6ed2d3129a7f129d5b9e8f1e0ee2a941ef26b1d4476cd947b22c232d1`,
  three-quarter
  `5980a9058ec66c70b30fb97d6d0a64a5ec126c3d01994713ae4ea7aad3301894`,
  and presentation
  `ee01026e2c0556ea197ccdb96e82f213d3ec38c8b0871dbeb376354acd8cfe19`.
- Attempt 3 aligned 42% front overlay and reference contact-sheet hashes:
  `146fefcfb778983783472f3ce7328a8a3c3f1a4f5358f6df651fe0577c1f903f`
  and
  `289b06ea112bf8735cccbbd2551e1bdc97f353f908d3feb202e3735185736809`.
- Rejected Attempt 4 candidate and metric packet SHA-256 values:
  `bf75e1b01d65a68066c41d0f7a070cae12f961dd814eb42183272208a9384262`
  and
  `177f52f5101323f12fd14e746fc1f562bb729afdd20c11fb0ed2c28996727e76`.
- Attempt 4 diagnostic render hashes: front
  `b0d1c05378e8580ba53f64945ff633fac746401bad6d743d54e73cc35f35b183`,
  side `bd3a0a7af5be1565cbaee4ac160ebe74eca342c5c2c0d4d89fbf2a3ffa2366c4`,
  front silhouette
  `364e7a2cdefc6687a08546c969b54efb2a656e35d6825810002cbac1ebc6739e`,
  and side silhouette
  `c87faee50162c9c2b6cba22e41b70b7c9568288a4748dc3aaf45ca9a6ff9e8d5`.
- Attempt 4 aligned 42% front overlay and reference contact-sheet hashes:
  `033f38ea6646b73334577d6761daf5d6ef40d348406ab467822d21629a752ee3`
  and
  `d54c59551a1c89bd3507a73c84fd1c520688aa2d138fef1578d0d13d02db3838`.
- Rejected Attempt 5 candidate and metric packet SHA-256 values:
  `57c4d7e4bb5a7b65299e8363d637a227944c474685ad16cf5434f117ea6ef528`
  and
  `a4daf8563e23b2e48347dbb0c9763d998cff64f14a50b6c79f0f86352f01cafd`.
- Attempt 5 diagnostic render hashes: front
  `f884bdeccbbd1f06db0af7e065586a1d82270ea1f36536eb0e1838d16295352e`,
  side `d659814f71bbfb4aa19ab1546efb32a7d9788fbf69763108ec047d4926d0180e`,
  front silhouette
  `24ec07ee002672a83cf78e1bec2b52382b9f40ef9610e1a841031e63efe40330`,
  and side silhouette
  `2b7163f5d8be8857ae31a7e35a76d9f518ee35f2eff6cbe104078478a2777ee9`.
- Attempt 5 aligned 42% front overlay and reference contact-sheet hashes:
  `10efa0c640981d4b21e43108eecf9bfb9a032ac0c477b9a42259b8a42c84f762`
  and
  `50617ba423c85201f872ff34ad20ea5a340f3315f23403e720442c97f5fc3383`.
- Rejected Attempt 6 candidate and metric packet SHA-256 values:
  `3523bd9f7c01bd97274d2a9355be6f798c83f94cd2f59c877e86cd2dee6a8cfa`
  and
  `b893e352ed388c6f7d93111922b0853bf75b454c284baed956a12554df8d6c8a`.
- Attempt 6 diagnostic render hashes: front
  `f15524e636831dca7fa9076ef0903d37bb47ab0cd3acdc81c5b9edd755ba752e`,
  three-quarter
  `f7530596c456d5385b2ea2b0e74b1d731e28b5e555127e8a68407876c022141b`,
  grazing close-up
  `18b6b28ac4c22b1ad9711a04dd978debeab41b3e7d71acae3462657f0373f63f`,
  front silhouette
  `b4251b7203a2843cd5e9f18d05b86ead2c591cb240f664c460bfdf151ddfcc04`,
  and side silhouette
  `1eb4951ce0c210b353d7517eb3f7532707e44107e1c4d4a348c6f334544d702f`.
- Attempt 6 aligned 42% physical-front overlay and contact-sheet hashes:
  `6baae71e748a09031585d17c97ce5df3887c33ddb3439ab21396e493ac989cfc`
  and
  `54f33030fbe4933969e91986796bfbcbadce0300f1cd386b44f1fcd54a22a386`.
- Rejected Attempt 7 candidate and metric packet SHA-256 values:
  `b784a5190a9089bc0b997759934e5045d83dfecc5858743328ca37b8cc264d23`
  and
  `06bc148cd2702b0f8e6d0b149c1f90c05c39307464b8c92503714ed2d9396fbb`.
- Attempt 7 diagnostic render hashes: front
  `b6d7c490c5b63b5388f00c9c4dff793a94740a687d20c1781144c8c70ec4def9`,
  side `4547a996b622127778d057f898175b0dcf384625d1a83f0c2b6ffba99ed3eb94`,
  three-quarter
  `da691055a0bf7f4541d7eeae7893a8ca5201e3e80fbadc5d78f979d86c23196b`,
  and grazing close-up
  `890bfa04b86381376128d4fefcff4317d2c722e95e4f4f46ede7c48ebd8553f9`.
- Attempt 7 aligned 42% physical-front overlay, aligned reference, and
  contact-sheet hashes:
  `88009f7f82f769326462f5713553ba6893fb129cb294743283e4d79768118df3`,
  `3d138666c857d203832461d8e7d3b4d7f8a87879c8d6d30e8e87cc662a412e31`,
  and
  `80f7cd8b25f62507ed3dcb572b30ce255c2e64864e3c742a5d49d3a51978a499`.
- Attempt 7 front and side silhouette hashes:
  `00fa9ed7db8043f710bb7fd42445354cc359b4d20caa113167bdbc34e96de3dd`
  and
  `0beb93b9b1a9982f98b120b6cbc8197f4b3c04e1754dcab0b283e8143d284a7a`.
- Rejected Attempt 8 head candidate and metric packet SHA-256 values:
  `01ce5141f171550c74c18af801dde3f90e60389aacc6257d43bacf4c9fee4dc0`
  and
  `50eead12c019969d69e942252d0af75d508f0afffaeceeb20ddae6600d0e479e`.
- Attempt 8 head render hashes: front
  `75857997e133ceadbcd62d5001cb794afedef8eba169f4bc8d5af19a4b271e50`,
  side `abcdccc6c285b27c737de4c0c9db62811d26afb8bff48e3ed9d73febc35d089f`,
  rear `698fc88b7474e606d64202d45f9217548c4629456d547f646acb5eaa86262a92`,
  three-quarter
  `2e6efff92faf824c8ff2976b4d1aaf06a040c0d7a800fca530e20acd828c4449`,
  and close-up
  `5f88b8164805e77fd9092d11751649ef1d809c1ac0978b614075f6a5d4b34246`.
- Attempt 8 head contact-sheet SHA-256:
  `5e86fc3411f20724a7f24b636e32ad466a1f3cc373dd8722a5a276445a0f258a`.
- Rejected Attempt 8 right-loop candidate and metric packet SHA-256 values:
  `367e1c029dfcc440107edc5017bcb62d6cd4ae05cb9b6b912f9dc4535142858e`
  and
  `5d6f557f6b6ccaea5739505e1f2979ec6c99fd5ff2657977b9a8d688bfa45199`.
- Attempt 8 stable right-loop render hashes: front
  `f5f4474210c170ba98ea72d8786cc520879def993393ba7b41d23ebf12e1f9d0`,
  side `c7a5f8840dc6b117650a1aa048c89ad4a24757c27f64b8c2c4334b9dae55295b`,
  and three-quarter
  `d9d84fc746e39a0f7cdd07beba2ff114ccb97b9235a335aa94fa5c1a722f4a36`.
- Rejected Attempt 9 head candidate and Gate A/B/C metric SHA-256 values:
  `9161f289277ec3e76fd7090c102f01206c7ae7b90b937e7338a4ee46c5550f8f`,
  `4680e587e78e8d46ad41c7317d6ac1685fe9110a09d9c3ab1be6f38f89d2aa39`,
  `2f484015f26730e2619c77bbcf9391a42844b572eed3758ec056b6fdb086f0a7`,
  and
  `8fbf39505837c5438fc6f4fda7ae14694c69ad858a60647baf24683ccc83317c`.
- Attempt 9 final head render hashes: front
  `5bcd78d91365570e4da746e60156a103c4ab01f881afdfc8a386b1fa2ce7f5ea`,
  side `a496722f558fbf7eef9f12c5593e86813f086d3524067893b7e2a25a457aa090`,
  rear `d0151aa76dc269b77bbe834477852086afffac25f9b35daca2d37f74552ffe79`,
  three-quarter
  `1f4bb92875fc4a5b2093cff18faf91f7be3890517943e1b836f72c364b86f7cf`,
  grazing `0882752a1fa2f329f4932a8dc8051aa0a79648452db1bd7f73b2ac1436b72b25`,
  and contact sheet
  `0dd69259721e0cc479531f7e6535fe0e55cd55050cb75854a7481fd2a5d5cd90`.
- Rejected Attempt 9 right-loop candidate and metric SHA-256 values:
  `14d14faabd03739f0c7c0f415d96a6e27c44c2fb2474839ecdf9a73db39277fd`
  and
  `4a2f30a97f509de854923bab7c4396cbe6dae6502be2210153598affa0723150`.
- Attempt 9 final trimmed-loop render hashes: front
  `7ebf18b6cb8d6a2e3f90d038ebab7327f5419ca33399075aa54603c9dac0cfb5`,
  side `563d2abff4c175285a4f7f61aeb571bf61a3f81273cf52553e636e0f0bb6b100`,
  three-quarter
  `690c4f652ae536fd36d976ae2f6da623d0d28dd0d8c9d59e0c12b4db86a025c7`,
  close-up
  `c4ca71beecd97ea92a30114e83d1f3a627faca7741e32d508e87a4ed9a1d5801`,
  and contact sheet
  `db28f7e421cfbf862eb34374ab5b17a917c89e0b9eec5f13a2187050b65fb052`.
- Rejected Attempt 10 head candidate and metric SHA-256 values:
  `3b2b9af18aed018f2c10387d91a72e8ef34ea9d2e0789998f024e423dbf8f379`
  and
  `bd4e89017dc022f724e5c93278007c294e302860d10c508d9fdc074fd72cb048`.
- Attempt 10 head reference-comparison and scaffold-contact hashes:
  `a5e0803af587ca25eb63f3466296c13b5b50e485ec353811b907d1a72f0e6c51`
  and
  `4ead5e6786419d1d23f5fc6738b6554048a97ae80a9c5bcee0b91bd7ecea4b6d`.
- Rejected Attempt 10 right-loop candidate, metric, and contact-sheet hashes:
  `f335be3d109b772fefd154cb9fd5ba0b4b9ac8db415b53c1814e968100dd0e91`,
  `3649fdc204061bfbc498612257309ded8009047c1a0c1f60d9d6d49b13405d32`,
  and
  `156f6ed1ae3faefe18de9a68625759b7e0d76c01deb519661606778f0a1d4771`.
- Attempt 10 head and bow reference-board hashes:
  `7154552ad070dd65bacba8dfc9e844128c053d84442dc2345b978c3be2428083`
  and
  `e19c7e27f97cff0e88cb8f44f19f6840f67c7d835ee49df01ce5eb41d2bb3623`.
- Attempt 11 reference-only head macroform audit SHA-256:
  `381e8df5cc2f48bc896e8a1ead537a4540090c546cc33b786c531c80af0179a3`.
  The audit is temporary evidence under `out/`; it corrected the bare-cushion
  normalization and stroke order before candidate creation.
- Attempt 11 interactive Sculpt Mode probe V3 script, result, report, and
  disposable blend SHA-256 values:
  `b8c5136ad9a8da63715307d77f29f67587a17e02de8c5d70e5ff35d5a17e9ab7`,
  `fbf9e39c8d72de7b60b74e4d57ca8019480e130168e47f7cd269aaac31b9415d`,
  `1371780fbff03799801670fdd4d16705d1a9fd4d0fb0c816b8eac1e9a5e894e6`,
  and
  `7b09a13114004ba2b8a8ee3559d8f745369ab729cab20d5da7e1abd3e4dfcc76`.
  Required stroke execution passed; per-brush candidate restrictions are
  recorded under Attempt 11 rather than inferred from the aggregate verdict.
- Rejected Attempt 11 head family-0 builder, untouched setup, checkpoint,
  metrics, reviewed state, and contact-sheet SHA-256 values:
  `2ec205915051b2e047f19a254628b4fa4cb19753485e23da68e2ff908d869b74`,
  `917ffac9e7184e075f721956dbfc3128771e93b3dfb05be35f11aff1b9783168`,
  `1eab348b4f9b47513c3b8b0c52dadbcb819b3fb62184dec455c8431b46ad2dbf`,
  `d2b3fec2fdb8f98a373701f1047b47921c7c9b1cfa67df90b115a9910ecea021`,
  `8b0433203fc592329e71f4b31f1ff7d8bfa7e0f0aab467e8e1924e69227e8fd0`,
  and
  `bb44017ef63c422d6d6210ea323ffbc0eb65ec066cf4e13cb00d353849bb4687`.
- Rejected Attempt 11 sharp-cube Grab-reset builder, untouched setup,
  checkpoint, metrics, reviewed state, contact sheet, and scorecard SHA-256
  values:
  `a7365068fc837bdb5600731019f031c5523ad0dada42a9510fe78085fefeb755`,
  `027358982b62f389e6b839d6568fb83fcec8a26bd7435a82d0971a174f593d21`,
  `4e1c12e1a9bbc1516e58eb9cc29d47432a4c642af6477b7b23a9cc6dccfe9eec`,
  `146c6ca4c67b0e4cf7577a767b5017e477368363c232839b037085c2edbbf957`,
  `575efe20b61e7ce2f71a8bb22fc4362a33a74a052ecf9939be6702df35ddc818`,
  `6912501cd4c176ba1e4dd2a4bf05ce0e99ea9cd9ac01c04d1cf8b76c459837db`,
  and
  `8920e2c96173dc3860f05c6eb06996788348d4642083cf89a69e38be540a6098`.
- Rejected Attempt 12 builder, audited base, untouched seed checkpoint, seed
  metrics, reviewed state, contact sheet, and scorecard SHA-256 values:
  `251cc2bfaf27c067e0879141ebc0ff3bedc2d3d6e78fe857c5a55ccf6549ed8e`,
  `a7365068fc837bdb5600731019f031c5523ad0dada42a9510fe78085fefeb755`,
  `90d93bd9221b37914a9f47bc7da2d6cfe6e41abd2f7f6b9cbeb2d32a2ce0c8d6`,
  `c873bec4cc57af95d375001e9cfc1f6f93603f346736d3cd0d16e0d4466cab89`,
  `ac3a9f7763286299a8bb4219f57bdc78d5c1390037c894d1713c4fb15a2deeaf`,
  `9200edaf6243240a9d8a009702eba0a0415d2f125fb6ef87dc2fa7b6bb16d034`,
  and
  `831e992a4f148036a1aa1098777f818763f6eb1020f95c2690721c1d5930a5f0`.
- Rejected Attempt 13 builder, untouched seed checkpoint, seed metrics,
  reviewed state, and contact sheet SHA-256 values:
  `752e3c9725a4fd3f9bb9a6f36b2e64eeb9c14c31386d8b5cf4ad3155014b8ac0`,
  `2ff65a7136d67b28b0cbbb1f352d5f7ab0da31658516cd6ca6c40ec60c64e25c`,
  `4d5d85a91f4875c27e07700c2836a9e10e6f24d9b36b7485ceb125dae8255cab`,
  `652f617cc5c4e3f823527990e1a2d8aed9f7296538d783cd7b6e90741b447ad5`,
  and
  `80d850831407e3508869a2b033aeb224893b67f81a942968028943f01010b0b0`.
- Rejected Attempt 14 builder, simulation-setup checkpoint, terminal numeric-
  failure checkpoint, simulation probe, and pre-simulation diagnostic contact
  sheet SHA-256 values:
  `8bee397b83d2f2661e57c04e3be32465421893423a0d7abc129e8aecbd25fa06`,
  `fffab0271668ed7b6a78ddf0d0ed8bcf95d7d091164de9eefcfcc548c0022c75`,
  `defdf74884f017a9684645485854e673d5254c1019fff8151cfc884949d71004`,
  `02e8604aa71f990103acf9e3da8e58af320ac9ec417929eaf10d7200ccdd35ca`,
  and
  `e31831239ae73515d6d29582acf2d8a7ef888c7802b7b45061935d6f4655b779`.
- Attempt 15 dormant driver and terminal preflight-audit SHA-256 values:
  `b8f53ef5ac8344aa3e605d503d1a14f5e8d252b0af2ef77f6049775a67fdfa6f`
  and
  `348f39dc0d4961b5d716897452c8698cdf08f25b9f33fd1786a0d80d1ed24f68`.
- Attempt 16 front target, side target, reference board, and terminal guide-
  review SHA-256 values:
  `16d46eac0d4721b2f50fd70cc82b4ac096049e1e15970cbb3a6c4d050f98ad01`,
  `fcd8b6d8d2483f5a19fb05cd9931b2ea00f97e3a32d30432d3965f5240571d25`,
  `a06ddad6c763b84da4b5b418983ca3cb6ce8c05a0e4a07e76cddf7cf84c932e2`,
  and
  `dd95b347e554afd1e4fab18c4663a21187fba8098b57528cb766d10f07e45f92`.
- Attempt 17 manifest, review packet, and terminal blind-review SHA-256 values:
  `ef18d6cb3bcc174db406422ac2065b9e30a1ee29fa76fe602e7d36de5cfdb566`,
  `b6c19edfb977ad183cabf9eb84f0e5eaf6c852f58f966129a35c6d6ccceb9e41`,
  and
  `b3ab03b5eee0588725dc40a5699cafe386de3816ecaa85d8f20d2d36ce34571a`.
- Attempt 18 source-aligned panel overlay and terminal preflight-audit SHA-256
  values:
  `76e7ec23cc3e35ed5508ff5dc292e6ad9f745be451343298e72338deff453997`
  and
  `eb6f8213fceb5eb9a5befcd7d97d4f72f297650f40b1110b5467270f998b7a60`.
- Attempt 19 terminal support-preflight audit SHA-256:
  `6c5bc11057d6008dc412371ba94e6438a9f96446fbf6b47816ef44de47f6110d`.
  No driver, candidate, render, or Blender mutation was created.
- Attempt 20 independently reviewed source-ownership packet and exact raster
  metric SHA-256 values:
  `802fb8fa69b1a847c50a3a3f854b1f56b3482a792b3bc1e1f6eed3f960108f2b`
  and
  `b27c2cc2506f844cc30634b523adbf4c87bbe81287209e75ba2ac30174749810`.
  The source contract passed; its smooth witness established feasibility only.
- Attempt 20 final native Multires probe script and result SHA-256 values:
  `7da897287a71a62165c905e6437f95b69f8dc99f47f8e92827411e995d33cc44`
  and
  `9c6d31ac3c685b3adf214ff309cb3b7bc5bc76c3b6d29534ed2b956a7d7901ed`.
  Both Sculpt and Mesh X symmetry, fixed-index displacement symmetry, per-lobe
  direction, before-state locality, bilateral extrema, topology, and evaluated
  save/reopen gates passed.
- Attempt 20 rejected ordered-sculpt result, disposable blend, contact sheet,
  registered mask metrics, and registered reference overlay SHA-256 values:
  `3e7011f6e44951e55c9971d51436490a13e3d20736d6e8afa5a0219d7779bbe7`,
  `7c49e92fa841ff9d0148649d573a8aa8fbc25de9bbd1a90da94a75038bd8b3ef`,
  `5a897db5fa7c8c3c37bcc08367e0d904a9b9e2e48d0ddc80b15579594e0425cc`,
  `d27b2cab60823fccc17354346f9214282eccda769852c24c0ef564f62f2f759c`,
  and
  `4631c867cf379cdf26a8338c10259bbfbcfde78ca9b7c2b4cde536112a432b25`.
  The ordered sequence landed at `0.923 × 0.739 × 1.010 Wh`, and every stroke
  passed the mechanical audit. The independent visual review nevertheless
  scored front `1.5/10`, side/rear `2/10`, both three-quarters `1.5/10`,
  continuous stuffing `2/10`, soft manufacture `1/10`, later-panel
  suitability `1/10`, and presentation `4/10`. It identified a spiky carved
  creature mask or deformed foam ball, not a plush head cushion. Registered
  beige coverage was `0.8666`, crown coverage `0.7715`, and chin errors reached
  `0.10–0.13 Wh`; the candidate was rejected before plan freeze.

## 2026-08-29 Attempt 29 V3 mapping and Stage-A resets

- Native-support mapping audit SHA-256:
  `e41f3d8f91e9f8c5638588c6b663c4dda1cb608e1dbf8edd8b1ad9920474c06f`.
- Rejected radial/geodesic diagnostic, metrics, verdict, and candidate-stop
  evidence SHA-256 values:
  `ced1c97e4e26a197382b95078c45268ef3e1704b979322a83fa6c3de57039633`,
  `44beb2d183f9b7686a78af25d39d7078ed8333b55551835581e97ad2fb61d3da`,
  `62202cfa13b982f67427a2f26dd1fc140cc9f8c72fd31636fdfb137fee20499b`,
  and
  `2bc13aad09dd9bbf83dfa39f3b4966253f23cda889837e288ba00c241c50cf9e`.
- Alternative constant-native-course builder, rejected blend, metrics,
  verdict, and four-view contact sheet SHA-256 values:
  `4fbda2359b32431154d506b9465f176bac1241766b382412f5d45a94c5714c9f`,
  `b9fd054f119fd2171d40c443613754ce779c4329b8bb993d659952a31c8f57b9`,
  `b7855e31cb94e7fc0901d9f9feb094be0c59dc0c97dedfaf1012aef9632cbf29`,
  `132535792858dca9a3f00e383c2974b70ff3ce2b1a21ab63476f5de0cb275b06`,
  and
  `d52bfdd0cbc460986510f2fd4d4b2c38b358e5ed51c2c368b500a8e0b955428b`.
- Regularized-seating witness builder, rejected blend, validation, reset
  report, and four-view contact sheet SHA-256 values:
  `8662cd95596e11869c63663d7e02ab238cdae2336b4baca54aa62ae8139afc0b`,
  `8675249f405707ec26f58da1fe90e38161a3b4e2b50faceb87e33eba2e0b96c2`,
  `1032e8178a63bb412c68f09c9968b9e19c3c5f8aa7668e16d5c1f72ede5f3377`,
  `39ec30a35a071daedf9f58acd294b0c7641123071ebb74a48ca72218b8c41724`,
  and
  `495df9172e9a6d58c13c15fc843c64c2a50c3c2de5430cf9a4585b63661db1df`.
- Protected support, tracked Reimu, and Sisyphus hashes remain respectively
  `b6b8b84742607d66f01f87362a69f3fa48bcdbb28f5a491192ae4141d4648328`,
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`,
  and
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.

## 2026-08-29 front-crown crop coupon V2

- Builder:
  `8e296670acf3289eb0a7a05d7273553ce382295c3337479e8c6ead6255572490`.
- Candidate blend:
  `d406f298ca76d153f40c5d39123920fecdb6410130fe322f82862e0484091e37`.
- Validation JSON / report:
  `0b90081a5db4827202d94474f09812aa7f4952524930ed4e174ba2fcbc4ed8a0`
  and
  `0cef0bc51bdc4561d04c8bffbc65c08b55ca7a3311df98b3567dc9b729d9465a`.
- Candidate sheet / reference comparison:
  `5ccf90e415e63427cfdb0e5158d674232773072a56d4ffb1ade34d45917c1254`
  and
  `4e39933df37ca490028b47b8056955d692a2f07c09a979edc3f34f34fb57913e`.
- Decisive render hashes: edge profile
  `e1e774c04aba832ca9f09dd1c8c5359a7248584e1c7f11d6eafa7b1de7a1ffeb`,
  grazing left
  `95f63a1cbc8616ef5c16ea2b23a89f6ea96fcdaf15f55f5f790a662562ca4bbc`,
  grazing right
  `aea1c928a4ed41471664845dca32f0c0460820f9afd02be1a3811ed097bce400`,
  and top crop
  `8a6b7821e74ff6bc1ea30d881ed667e55835a6b51160f2bb9a3d9348e82419a0`.
- Geometry: `2,277` vertices, `2,272` polygons, one component, zero boundary,
  wire, or non-manifold edges, no modifiers, X
  `-0.888846..0.785416`, Y `-1.297007..-0.607824`, Z
  `5.950000..6.987741`, projection maximum `.062716 BU`, seated-row fraction
  `.9091`, main thickness `.046405..087000 BU`, and width/thickness ratio
  `18.391`.
- Protected hashes remain: cushion blend
  `b6b8b84742607d66f01f87362a69f3fa48bcdbb28f5a491192ae4141d4648328`,
  tracked Reimu
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`,
  and Sisyphus
  `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591`.
- Absolute blind verdict: reject at the decisive gate; seating/contact
  `1.5/10`, fabric read `2.5/10`, and crop closure `2/10`. The remaining front
  and neutral three-quarter renders were intentionally not produced.

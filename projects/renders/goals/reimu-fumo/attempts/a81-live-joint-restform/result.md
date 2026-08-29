# A81 result — live joint rest form reset at the controlling pair

## Decision

Reset A81. The direct-edit workflow changed the authoring modality, but the
joint support/yoke representation still failed the first visual gate and is
worse than protected rung 003. It is not a candidate and must not enter the
tracked reusable model.

## Observable result

The frozen live-edit snapshot is
`a81_live_joint_restform_pair.blend`, SHA-256
`bc1cf8fb4fb669076ac3199210fedcf43b534cea70d0048b12f4a1f4d6e197f2`.
Repository-pinned Blender 5.2.1 extracted exactly two 512 px views without
saving the source:

- exact front, SHA-256
  `145f0c4ba76581d4947b83759824f353404fa66d2271a59bc4aab494f82eebe8`;
- weak three-quarter, SHA-256
  `a1d26222a1ded976deb8e341dfcea9ed8f7650adf83b001855dcbf5f55100252`.

The front is dominated by a full-height beige bald support with only a small
central brown fringe and retained cheek locks. The three-quarter view is a
deep faceted beige rounded box or egg bordered by a narrow brown strip. It has
no continuous brown crown-to-nape field, no broad returning leaf, no readable
overlap, and no complete crown or rear coverage. This recurs `D-FACE-MASK`,
`D-HEAD-HELMET`, and `D-REAR-ROOT`, loses the canonical front identity, and is
plainly worse than rung 003. No repair or third view was authorized.

## Feedback timing and process result

Active authoring began at `16:34:55+03:00`; the controlling-pair deadline was
`16:46:55`. The pair blend was saved at `16:44:53.037`, but a foreground Eevee
render produced no PNG after more than 120 seconds. The coordinator stopped
authoring at `16:51:03`, 4 minutes 8 seconds after the pair deadline. The
separate pinned batch fallback later produced front at `16:53:21.301` and the
pair at `16:53:26.918`, 6 minutes 31.918 seconds late. Once invoked, the
reusable batch renderer needed only 17.906 seconds for both views.

The useful process change is therefore narrower than the modeling result:
direct sparse editing avoided another large generator, and the standing
reviewer applied an immediate pixel veto, but foreground rendering was the
wrong evidence path. Future live edits should save a frozen snapshot and hand
it immediately to the proven pinned batch renderer. The inherited
receiver/cap/yoke family is now closed; the next work unit must start from a
fresh atomic head-and-hair shell whose first complete blockout already assigns
crown, temple, and rear ownership to brown hair in the whole-plush view.

## Protected state

- Rung 003 remains SHA-256
  `c538a9aa070c4f0e127b6ace3b42220ae096c6e7a7fb1791b8906fd02f78bd3b`.
- The tracked reusable model remains SHA-256
  `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76`.
- No canonical Blender asset or reusable structure was modified.

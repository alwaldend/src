# 060 causal implementation repair 1 of 1

The first builder exited2 before saving: the free front lock's X coordinate
-0.0516019m has no head surface at the lower sampling clamp Z0.112m. The new
head narrows there while the hanging fabric legitimately remains outside it.
This is a wrong root-domain lookup, not missing head geometry to inflate.

A read-only probe of the exact planned cushion and050 lock coordinates found
275 ray misses with0.112m, zero with0.118/0.123/0.126m. All side-panel root rays
already pass. Source bytes unchanged; probe saved no candidate.
Evidence: out/reimu_fumo_finish/desktop_astra/macro_060_root_probe.log.

Use the smallest passing root clamp,0.118m, for front hanging locks only.
All XZ silhouettes, head geometry, side panels, eye domains and sleeve geometry
stay unchanged. Predict zero front-lock receiver misses and a saved candidate.
Rerun the same060 build to its fresh path; it does not yet exist. No further
implementation repair is authorized in this attempt. Preserve the initial log
and write the retry log separately. No visual acceptance claim.

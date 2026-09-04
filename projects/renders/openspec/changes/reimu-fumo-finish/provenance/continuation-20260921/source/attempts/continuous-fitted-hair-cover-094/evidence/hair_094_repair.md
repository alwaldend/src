# 094 sole causal data repair: source-descendant fiber mapping

First build failed only at pile transfer after full array/native geometry,
complete actual contacts and mutual exterior checks passed. No candidate was
saved. Builder fee44fa31aa00188c5b9dc890f5d25c69a4dad5eb406fe7b113b4dccd8850be4;
first preflight67ab475e4693c374172dd703cb684cc65c6354f64dbefe5df7d704fdc01280fd.
Exact frozen array packet492c89b457110d01273a1d8de0616dc266d18c5c441aa634c871395c54946b7f.

Independent static review found that the original code discarded source_ti
and ran a second global nearest query on an equivalent split mesh, losing
the promised ancestry. Read-only native diagnosis of all130000 source strands
found64 source-owner mismatches and up to2.875283um spurious nearest-point
displacement. The first failed witness is on unchanged graphic triangle6303,
with only2.241nm source-plane residual and positive dominant-axis double
barycentric weights(0.455604,0.090599,0.453797). It has only one descendant;
the second BVH query is both unnecessary and less accurate on skinny faces.

Use original source_ti and canonical source-barycentric ancestry to enumerate
ONLY its actual descendants. Compute double point/triangle plane projection
and, if outside a descendant, double closest edge; choose the nearest of those
descendants, retaining the original200nm spatial bound. Transfer through those
same triangle weights. Never switch to an unrelated material owner or increase
the spacing bound. Verify every descendant node references its source face.

This uses repair1of1; no further repair remains. Geometry helper and all gates
stay unchanged, and the rerun array JSON must hash exactly to the first packet.
Write new repair preflight/array/receipt files, preserve every original log and
failure report, same new094candidate path only if all gates pass. Source087
unchanged. This corrects data correspondence, not a geometry or visual variant.

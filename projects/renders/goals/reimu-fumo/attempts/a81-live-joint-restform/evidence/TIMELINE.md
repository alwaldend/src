# A81 feedback timing

Timezone: Europe/Moscow (`+03:00`). Filesystem timestamps record artifact
completion, not work start.

| Event | Timestamp | Evidence | Elapsed from author start |
| --- | --- | --- | ---: |
| Canonical attempt created | 2026-08-31 16:35:34.843199820 +03:00 | canonical `attempt.yaml` `creationTimestamp` | n/a |
| Active author started | 2026-08-31 16:34:55 +03:00 | explicit author signal from `/root/a80_constructed_cap_probe` | +0 |
| First pixel | 2026-08-31 16:53:21.301378600 +03:00 | post-decision pinned fallback `batch_render/packet/front.png` mtime | +18m 26.301s |
| Pair complete | 2026-08-31 16:53:26.918081264 +03:00 | post-decision pinned fallback `batch_render/packet/three_quarter.png` mtime | +18m 31.918s |
| Decision | 2026-08-31 16:51:03 +03:00 | coordinator: stop/reset authoring on missed pair SLO; preserve pair blend only | +16m 08s |

Budget: pair complete within 12 minutes of active author start. The deadline is
**2026-08-31 16:46:55 +03:00**. The canonical attempt was created 39.843s
after active authoring began. The isolated pre-edit snapshot first appeared at
16:36:41.239856360 +03:00; it is not a pixel milestone.

Status at 16:46:55.818 +03:00: **budget missed**. The pair blend had been
saved at 16:44:53.037 +03:00, but no first-pixel or pair PNG existed at the
deadline. The coordinator was warned at 16:44:21, 2m34s before the deadline.

Author status at 16:48:28 +03:00: exact-front 420 px foreground Eevee had
produced no PNG after more than 120s, and its in-flight wait was terminated on
the coordinator's SLO stop. At that point first pixel and pair were unproduced;
no synthetic live-render timestamps were assigned. Frozen pair blend SHA-256:
`bc1cf8fb4fb669076ac3199210fedcf43b534cea70d0048b12f4a1f4d6e197f2`.

Coordinator decision at 16:51:03 +03:00: stop/reset A81 authoring on the
missed SLO, with no repair or third view. A separate reusable batch-render
fallback may extract the already-frozen pair afterward; any such timestamps
are post-decision evidence and do not change the SLO verdict.

The fallback subsequently produced front at 16:53:21.301 and three-quarter at
16:53:26.918 +03:00, respectively **6m26.301s** and **6m31.918s late** against
the active-author deadline. They arrived 2m18.301s and 2m23.918s after the
stop/reset decision. The successful pinned command itself started at about
16:52:54 after one sandboxed output-base preflight failure; first/pair command
latencies were approximately 27.301s/32.918s. These are fallback extraction
times, not on-time live-author feedback.

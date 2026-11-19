# 059 metadata repair 1 of 1

The original builder saved head_059_candidate.blend, verified its protected
source, then exited 2 while constructing the receipt. The inherited boundary
walker overwrote global `previous` with a BMVert; `sha(previous)` consequently
failed. Candidate SHA256 is
b9ea8d91e6ea536d58ae5481dd24d17b492525d593c7103540c742a11599d44f.
Original builder SHA256 is
f98f1c394d801cf5d28759e6cce46bc2f127f512ec99b49e8a42ad1c23875acb.
Failure evidence: out/reimu_fumo_finish/desktop_astra/head_059_build.log.

Use the sole preauthorized repair to rename the wrapper's path variable to
`profile_source`. Do not rerun the builder, save again, or replace candidate
bytes. Clean-reopen the existing candidate and produce an explicitly recovered
verification receipt, retaining original writer identity, exit status and
unavailable runtime measurements. The corrected builder is a separate input,
not retroactively the writer of these bytes. Render only the verified candidate.
No visual acceptance or whole-asset success is claimed.

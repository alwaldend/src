# Continuous drape changed tips but did not remove the junction

Candidate SHA52fcbf7e39f0a54c6c4d7a223798a342884569f3527c119a63bff121c69c43ed.
Clean-reopen fixed front/side completed unchanged in21.03s. Longer tapered
ends and band seating changed, but the transverse dark junction persisted
at the same height. Reject before full review. The explanation that the
piecewise .113m drape transition caused this line was falsified.

Read-only visible inventory found only the intended new cream receiver,
front fabric, five rear panels and two side locks; no visible legacy hair
assembly explains it. A targeted pure geometry-function probe then found
the actual cause: crown_angles returns its no-intersection fallback before
applying the lower temple-wrap blend. At the front pattern's lower limit,
z0.1049851626m, angles jump from(-0.77,3.91159265) to(-0.3,3.44159265).
That is0.47rad per endpoint across only20micrometers in the probe.
The jump remains even after a continuous radial drape because it is a
separate angular parameter discontinuity. Evidence: hair_049_angle_probe.json.

This disconfirms the proposed representation-level explanation. Fix the
explicit fallback bug rather than repeat drape variations or retire a method
based on that misdiagnosis. Keep the2mm projected edge gate, support limits,
fixed cameras and visual rejection gate. No final criterion passes.

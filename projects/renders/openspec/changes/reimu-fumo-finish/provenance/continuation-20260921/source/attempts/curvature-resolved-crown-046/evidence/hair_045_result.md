# Rejected fitted-crown candidate045

Saved SHA8fdf926e19d7c441a8d77351f65ffe9380305196b098f4568707c9b1a450ab45.
Pinned clean-reopen fixed front and side render completed unchanged. Both
show pale scalp gaps and a hard horizontal crown break; reject, no full-view
or final pass. Front cf32dc84a0ecb570ebf6ddc1446631443181cd56a1af9619b7ef5f14f5b71a75;
side7101adf5f4baed3a9ce6e9673ef0e689b8e8cf6514e32a35bb59cd58c1155103.

Read-only mesh diagnosis narrows the causal fault: the projected front has
20.342mm edges despite subdivision,432 triangle centroids inside the receiver,
and maximum inward depth2.349mm. Finite vertices and depth0.712425Wh do not
prove fitted cloth. Rear lower edges now hang beyond the tapered receiver,
but upper cap coverage and temple coverage also fail. No acceptance retained.

Next attempt replaces fixed edge subdivision with adaptive edge-length-bounded
triangulation, reprojects after each topology refinement, insets the upper
stuffing behind fitted fabric, and wraps side panels forward only where the
observed temple gap requires it. Retain the same source044d, references,
camera and material settings. Do not repeat the same coarse projection.

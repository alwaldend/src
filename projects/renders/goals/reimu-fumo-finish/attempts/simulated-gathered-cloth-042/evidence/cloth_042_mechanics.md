# Actual cloth execution and causal repairs

Original planar042 aborted sleeve intersection before any simulation/save.
042b replaces it with a seated pre-shaped start, stationary pinned waist,
separate flat rest shape, and actual body/feet/sleeve/hand collision proxies.
First042b helper SHA527eb4e2c90c01272c3b64faee2b673fe48d7d0e6d776c4d031eec3687ad0190
failed conservative minimum-distance0.9mm at torso0.7637mm. Increasing only
torso radial clearance resolved that contact. Next observed foot0.3004mm
showed vertical offset was not normal clearance on curved sides. Convex-foot
normal separation raised it to0.6299mm. No actual intersection was found at
these tested contacts.

Root retained triangle-intersection rejection and a near-coincidence guard,
but converted the helper's arbitrary0.9mm preflight comfort margin into a
warning: the actual collision solver should separate small positive gaps.
This is a simulation-initialization choice, not a visual-acceptance waiver.
All visible mesh initial-overlap checks then passed, supported contact
distances were logged, and70 frames solved in11.77seconds with finite bounds.

Candidate042b SHA819207b0d96958e6e86110c30399db45db183772bf66cb03296671aef4ce764b
is saved as ordinary extracted mesh, not dependent on an in-memory cache.
041b source remained unchanged. Numerical completion is not retention;
fixed-view rendering and root visual decision follow separately.

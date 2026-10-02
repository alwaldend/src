# Fabric041 appearance diagnostic: insufficient

041 SHA2754c76c16192e256ac9835ded2ae86cf97703b93f7daa386c48db9c192dc866
saved and clean-reopened for five fixed views plus presentation. Existing
base mesh vertex/topology digests were unchanged. Local short pile, cloth
shader and studio light did not make the model convincingly plush. Fresh
implementation-blind image review identified washed-out color, rigid sewn
parts and helmet-like hair. Root agrees; no whole-character stage passes.

Read-only scene diagnosis found Cycles active, no material override and no
compositing node group. One relative-path diagnostic failed because Bazel
runs in runfiles; the absolute path opened. Blender5.2.1 uses the scene
compositing_node_group field rather than the removed node_tree field.

041b changes only light power to15percent, world strength0.10, sheen0.045,
and64 samples. Fixed front and side show restored red/brown/black color,
but the skirt remains a conical shell and the hair smooth. This is an
experimental source, not approved geometry or a finished material system.

Decision reset. Next causal change is real cloth simulation with a pinned,
gathering waist and body/foot/floor contacts, replacing only skirt and hem.
The current animation/export criteria remain failed or unverified. The040
structured review decision was reset; its prose word refine did not retain
or approve040b. Reusing unaccepted geometry for experiments is not acceptance.

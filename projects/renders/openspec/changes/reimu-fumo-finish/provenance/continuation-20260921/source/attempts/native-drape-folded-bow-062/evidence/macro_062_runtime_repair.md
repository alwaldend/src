# 062 bounded runtime repair

First builder execution ended with timeout exit124 after180 seconds. Preserved
macro_062_build.log SHA256
08f4f547e6fbe85f1892fbdac73af76a3c513ad8416d85564e23b840f111a5d8.
The log shows simulation frames12 and24, no Python exception and no save.
Separate stat checks confirm candidate and writer receipt do not exist.
This is evidence of an insufficient wall-clock allowance, not evidence that
the simulation's numeric or visual result passes.

Consume the plan's one pre-save implementation repair by changing only the
external timeout180 to360 seconds. The actual simulation vertex count must be
read from the eventual receipt rather than assumed. Geometry,36-frame solver,
self-contact, numerical gates and builder bytes remain exactly unchanged:
17e84f4cf5248460cadd305a1dab421c1134303cbfa6ba499b90b51a7663fbc6.
Predicted effect: the advancing solve can reach frame36 and report its actual
gates. Use a new macro_062_build_longer.log; do not overwrite failed evidence.
No further retry or geometry alteration is authorized inside this attempt if
this execution times out or fails. Root remains the sole writer.

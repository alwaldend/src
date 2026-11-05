# 053 repair 1: copied source-hash literal

The first pinned 5.2.1 writer exited 2 at the line 18 source guard, before
opening Reimu. No head_053 candidate or writer receipt exists. The actual
protected 050 SHA-256 remains
6098307ff3b44052bd31ab73e80f1f2df1c8dee5ae79391f1f0aa5fb685dd93b.

Cause: global replacement of label 052 with 053 also changed the embedded
hash substring b44052 to b44053. This was a script-generation error, not a
changed protected input. Correct only that literal; preserve the assertion
and all model parameters. The next run should pass this guard and reach the
bound operation. This consumes repair 1 of the plan's 2; no visual result or
model acceptance exists. Do not use global numeric label replacement again.

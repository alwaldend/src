# Camera-space coupon evidence

The local Blender MCP background interface opened protected A157 and placed
mesh `MCP_CAPABILITY_COUPON` on the ray from the subject center toward
`Review_front_Camera`. It saved a new candidate with SHA-256
`3d18bf556ef34f6ea5970f487a277f6eeca1eb824ae9ce2ee21defd92dc0c003`.

Pinned Blender 5.2.1 clean-reopened that exact file with automatic file
scripts disabled, found the marker, found no missing external images, and
rendered the frozen front camera. The render SHA-256 is
`728fa04f484c9d431b7a28e9abf183f5d506bd7f87fd5bc5bf6a49cd629a0ca2`.
Pixel inspection shows the magenta marker visibly covering the subject.

Afterward the protected A157 and A202 sources still hash to
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`
and
`6a9f3757facba526550e78817dc85f1d23cf85bcdad360228e113bb60d5f3aa0`.

This establishes only the local MCP save/reopen/render capability. It provides
no reference-fidelity, structure, or technical pass for a deliverable.


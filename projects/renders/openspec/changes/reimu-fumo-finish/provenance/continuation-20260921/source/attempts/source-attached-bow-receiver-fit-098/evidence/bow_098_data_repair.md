# 098 sole API/data repair

First build58d32bee630211ad1d727cb5dd42124fee2fc0b53249855a42c42ab53bb576f5
stopped in source evaluation before any model mutation: reused evaluated()
requires bmesh, omitted from the wrapper imports. Source native error is
NameError: bmesh is not defined. No arrays or candidate were written.

The sole repair adds only import bmesh before executing the hash-bound
original builder, and changes the preflight/array output names so original
failure evidence remains immutable. No shape, input, gate or parameter change.
All extracted helper dependencies inspected; existing globals supply their
other dependencies. One repaired native run; no further098 repair remains.

# Preserved continuation through attempt 115

This immutable historical snapshot preserves the Fumo branch and its saved
uncommitted work when rebasing onto the OpenSpec migration on 2026-09-21.
It is not an active goal-tool store. Current work is maintained in the owning
OpenSpec change. No candidate is accepted.

`manifest.json` identifies every source path, original Git object, byte count
and SHA-256. `source/` preserves the latest complete record at resource version
375, including the 60 uncommitted attempt directories 056–115. The original
migration's sibling `source/` and `manifest.json` remain unchanged.

`retired-goal-fix.patch` preserves the branch's empty-evidence-directory fix
for the deleted goal tool. The rebase retains upstream deletion of that tool;
this patch is historical evidence and does not reintroduce executable code.

Historical absolute scratch paths and links remain verbatim for provenance;
they do not guarantee availability in a fresh checkout. Use the current
change's design and tasks for continuation.

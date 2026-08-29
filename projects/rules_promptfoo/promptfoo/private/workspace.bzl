"""Stages selected SkillInfo bundles as a Codex skill workspace."""

load("@rules_skill//skill:defs.bzl", "SkillInfo")

PromptfooWorkspaceInfo = provider(
    doc = "An isolated workspace containing selected Codex skills.",
    fields = {
        "files": "The source skill files needed at runtime.",
        "files_by_path": (
            "Source files keyed by their destination below the workspace."
        ),
        "skill_names": "The staged skill names, in deterministic order.",
    },
)

def _validate_component(value, what):
    if not value or value in [".", ".."]:
        fail("{} must be a non-empty safe path component, got {}".format(
            what,
            repr(value),
        ))
    if "/" in value or "\\" in value:
        fail("{} must not contain a path separator: {}".format(
            what,
            repr(value),
        ))

def _validate_logical_path(path, skill_name):
    if not path or path.startswith("/") or "\\" in path:
        fail("skill {} has unsafe logical path {}".format(
            skill_name,
            repr(path),
        ))
    for component in path.split("/"):
        _validate_component(component, "skill {} path component".format(
            skill_name,
        ))

def _promptfoo_workspace_impl(ctx):
    files_by_path = {}
    names = {}
    for target in ctx.attr.skills:
        skill = target[SkillInfo]
        _validate_component(skill.name, "skill name")
        if skill.name in names:
            fail(
                (
                    "promptfoo workspace has duplicate skill name {} " +
                    "from {} and {}"
                ).format(
                    repr(skill.name),
                    names[skill.name],
                    target.label,
                ),
            )
        names[skill.name] = target.label

        for logical_path in sorted(skill.files_by_path.keys()):
            _validate_logical_path(logical_path, skill.name)
            destination = ".agents/skills/{}/{}".format(
                skill.name,
                logical_path,
            )
            files_by_path[destination] = skill.files_by_path[logical_path]

    files = depset(files_by_path.values())
    return [
        DefaultInfo(
            files = files,
            runfiles = ctx.runfiles(files = files.to_list()),
        ),
        PromptfooWorkspaceInfo(
            files = files,
            files_by_path = files_by_path,
            skill_names = sorted(names.keys()),
        ),
    ]

promptfoo_workspace = rule(
    implementation = _promptfoo_workspace_impl,
    attrs = {
        "skills": attr.label_list(
            doc = "rules_skill targets to stage below .agents/skills.",
            providers = [[SkillInfo]],
        ),
    },
    doc = "Stages selected skills in an isolated Promptfoo workspace.",
)

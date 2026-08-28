"""Rule for aggregating the files that make up one Codex skill."""

SkillInfo = provider(
    doc = "Files and distinguished metadata belonging to a Codex skill.",
    fields = {
        "files": "All files belonging to the skill.",
        "openai_yaml": "The optional agents/openai.yaml file.",
        "skill": "The required SKILL.md file.",
    },
)

def _package_relative_path(src, label):
    if (
        src.owner.repo_name != label.repo_name or
        src.owner.package != label.package
    ):
        return None

    short_path = src.short_path
    if label.repo_name:
        repository_prefix = "../{}/".format(label.repo_name)
        if not short_path.startswith(repository_prefix):
            return None
        short_path = short_path[len(repository_prefix):]

    package_prefix = "{}/".format(label.package) if label.package else ""
    if not short_path.startswith(package_prefix):
        return None
    return short_path[len(package_prefix):]

def _skill_library_impl(ctx):
    skill_files = [
        src
        for src in ctx.files.srcs
        if _package_relative_path(src, ctx.label) == "SKILL.md"
    ]
    if len(skill_files) != 1:
        fail("skill_library requires exactly one SKILL.md in srcs")

    openai_files = [
        src
        for src in ctx.files.srcs
        if _package_relative_path(src, ctx.label) == "agents/openai.yaml"
    ]
    if len(openai_files) > 1:
        fail("skill_library accepts at most one agents/openai.yaml in srcs")

    files = depset(ctx.files.srcs)
    return [
        DefaultInfo(files = files),
        SkillInfo(
            files = files,
            openai_yaml = openai_files[0] if openai_files else None,
            skill = skill_files[0],
        ),
    ]

skill_library = rule(
    implementation = _skill_library_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            mandatory = True,
        ),
    },
    doc = "Aggregates all files belonging to one Codex skill.",
    provides = [SkillInfo],
)

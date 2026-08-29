"""Rule for aggregating the files that make up one Codex skill."""

SkillInfo = provider(
    doc = "Files and distinguished metadata belonging to a Codex skill.",
    fields = {
        "files": "All files belonging to the skill.",
        "files_by_path": (
            "Skill files keyed by their slash-separated logical path relative " +
            "to root."
        ),
        "name": (
            "The non-empty logical skill name derived from the final segment " +
            "of root."
        ),
        "openai_yaml": "The optional agents/openai.yaml file.",
        "root": (
            "The non-empty owning Bazel package path within its repository, " +
            "without a repository or execution-path prefix."
        ),
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

def _files_by_path(srcs, label):
    files_by_path = {}
    for src in srcs:
        path = _package_relative_path(src, label)
        if path == None:
            fail(
                "skill_library src {} is outside skill root {}".format(
                    src.short_path,
                    label.package,
                ),
            )
        if path in files_by_path:
            fail("skill_library has duplicate logical path {}".format(path))
        files_by_path[path] = src
    return files_by_path

def _src_files(srcs):
    return [
        file
        for src in srcs
        for file in src[DefaultInfo].files.to_list()
    ]

def _skill_library_impl(ctx):
    if not ctx.label.package:
        fail(
            "skill_library cannot be declared in a repository root package; " +
            "move the skill to a named subpackage",
        )

    srcs = _src_files(ctx.attr.srcs)
    files_by_path = _files_by_path(srcs, ctx.label)
    skill_files = [
        src
        for path, src in files_by_path.items()
        if path == "SKILL.md"
    ]
    if len(skill_files) != 1:
        fail("skill_library requires exactly one SKILL.md in srcs")

    openai_files = [
        src
        for path, src in files_by_path.items()
        if path == "agents/openai.yaml"
    ]
    if len(openai_files) > 1:
        fail("skill_library accepts at most one agents/openai.yaml in srcs")

    files = depset(srcs)
    return [
        DefaultInfo(files = files),
        SkillInfo(
            files = files,
            files_by_path = files_by_path,
            name = ctx.label.package.split("/")[-1],
            openai_yaml = openai_files[0] if openai_files else None,
            root = ctx.label.package,
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
    doc = "Aggregates one Codex skill declared in a named Bazel package.",
    provides = [SkillInfo],
)

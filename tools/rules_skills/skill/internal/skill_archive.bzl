"""Rules for packaging skills that live in an external archive."""

load(":skill_library.bzl", "SkillInfo")

def _package_relative_path(file, label):
    """Returns the path of a file relative to its owning Bazel package."""
    short_path = file.short_path
    if short_path.startswith("../"):
        parts = short_path.split("/", 2)
        relative = parts[2] if len(parts) == 3 else ""
    else:
        relative = short_path

    if label.package:
        prefix = "{}/".format(label.package)
        if not relative.startswith(prefix):
            fail(
                "skill archive source {} is outside package {}".format(
                    file.short_path,
                    label.package,
                ),
            )
        relative = relative[len(prefix):]
    return relative

def _skill_archive_impl(ctx):
    prefix = "{}/".format(ctx.attr.root)
    skill_files = {}
    for file in ctx.files.srcs:
        relative = _package_relative_path(file, ctx.label)
        if not relative.startswith(prefix):
            continue
        logical_path = relative[len(prefix):]
        if logical_path in skill_files:
            fail("skill archive has duplicate logical path {}".format(logical_path))
        skill_files[logical_path] = file

    if "SKILL.md" not in skill_files:
        fail(
            "skill archive root {} contains no SKILL.md; available paths: {}".format(
                ctx.attr.root,
                ", ".join(sorted(skill_files.keys())) or "none",
            ),
        )

    return [
        DefaultInfo(files = depset(skill_files.values())),
        SkillInfo(
            files = depset(skill_files.values()),
            files_by_path = skill_files,
            name = ctx.label.name,
            openai_yaml = skill_files.get("agents/openai.yaml"),
            root = ctx.attr.root,
            skill = skill_files["SKILL.md"],
        ),
    ]

skill_archive = rule(
    implementation = _skill_archive_impl,
    attrs = {
        "root": attr.string(mandatory = True),
        "srcs": attr.label_list(
            allow_files = True,
            mandatory = True,
        ),
    },
    provides = [SkillInfo],
    doc = "Packages one skill directory from an archive-owned file set.",
)

def skill_archives(name, roots, srcs = None, visibility = None):
    """Declares one skill target per archive skill directory.

    Args:
        name: Name of the aggregate filegroup over the generated targets.
        roots: Package-relative skill directories, each containing SKILL.md.
        srcs: Optional complete file set. Defaults to one glob per root.
        visibility: Optional visibility applied to every generated target.
    """
    if not roots:
        fail("skill_archives requires at least one root")

    targets = []
    for root in roots:
        target = root.split("/")[-1]
        root_srcs = srcs if srcs != None else native.glob([root + "/**"])
        skill_archive(
            name = target,
            root = root,
            srcs = root_srcs,
            visibility = visibility,
        )
        targets.append(":" + target)

    native.filegroup(
        name = name,
        srcs = targets,
        visibility = visibility,
    )

"""Aspect for validating skill_library targets."""

load(":skill_library.bzl", "SkillInfo")

def _skill_validation_impl(target, ctx):
    skill = target[SkillInfo]
    output = ctx.actions.declare_file(
        "{}.skill_validation".format(target.label.name),
    )
    args = ctx.actions.args()
    args.add("--skill", skill.skill)
    args.add("--expected-name", target.label.package.split("/")[-1])
    args.add("--output", output)
    args.add_all(skill.files, before_each = "--file")
    if skill.openai_yaml != None:
        args.add("--openai-yaml", skill.openai_yaml)

    ctx.actions.run(
        executable = ctx.executable._validator,
        arguments = [args],
        inputs = skill.files,
        outputs = [output],
        mnemonic = "SkillValidation",
        progress_message = "Validating skill %{label}",
    )
    return [
        OutputGroupInfo(skill_validation = depset([output])),
    ]

def skill_validation(validator = None):
    """Creates an aspect that validates skill_library targets.

    Args:
        validator: Optional replacement validation executable.

    Returns:
        A configured skill validation aspect.
    """
    return aspect(
        implementation = _skill_validation_impl,
        attrs = {
            "_validator": attr.label(
                cfg = "exec",
                default = validator or Label("//main/go:skill_validator"),
                executable = True,
            ),
        },
        attr_aspects = [],
        required_providers = [[SkillInfo]],
    )

skill_validation_aspect = skill_validation()

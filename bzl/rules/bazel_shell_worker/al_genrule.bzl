load(":al_genrule_rule.bzl", "al_genrule_executable", "al_genrule_regular", "al_genrule_test")

def al_genrule(test = False, executable = False, **kwargs):
    """
    Generate al_genrule target

    Args:
        test: If set, use al_genrule_test
        executable: if set, use al_genrule_executable
        **kwargs: kwargs for the rule
    """
    if test:
        al_genrule_test(**kwargs)
    elif executable:
        al_genrule_executable(**kwargs)
    else:
        al_genrule_regular(**kwargs)

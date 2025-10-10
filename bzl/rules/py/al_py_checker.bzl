load("//bzl/rules/run_tool:al_run_tool.bzl", "al_run_tool")

def al_py_checker(
        name = None,
        tool = None,
        args_bin = None,
        args_test = None,
        test_size = "small",
        disable_fix = False,
        **kwargs):
    """
    Create -fix and -test targets for a python checker

    Args:
        name: Name prefix
        tool: Tool label
        args_bin: Args for the binary target
        args_test: Args for the test
        disable_fix: If set, do not create fix target
        **kwargs: Kwargs for both targets
    """
    if not disable_fix:
        al_run_tool(
            name = "{}.fix".format(name),
            executable = True,
            args = args_bin,
            tool = tool,
            **kwargs
        )
    al_run_tool(
        name = "{}.test".format(name),
        test = True,
        args = args_test,
        tool = tool,
        size = test_size,
        **kwargs
    )

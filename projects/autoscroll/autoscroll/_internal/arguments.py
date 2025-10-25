import typing
import argparse


def parse_arguments(**arguments: typing.Any) ->dict[str, typing.Any]:
    result = {}
    for key, value in arguments.items():
        key_split = key.split("_")
        group = key_split[0]
        name = "_".join(key_split[1:])
        if group not in result:
            result[group] = {}
        result[group][name] = value
    return result


class ArgparseFormatter(argparse.HelpFormatter):
    def _split_lines(self, text: str, width: int) -> list[str]:
        """
        Do not split lines that start with 'R|'

        https://stackoverflow.com/questions/3853722/how-to-insert-newlines-on-argparse-help-text
        """
        if not text.startswith("R|"):
            return super()._split_lines(text, width)
        result = []
        for line in text[2:].splitlines(keepends=True):
            result.extend(super()._split_lines(line, width))
        return result

    def _fill_text(self, text: str, width: int, indent: str) -> str:
        """
        Do not format the description

        Built-in argparse class argparse.RawDescriptionHelpFormatter
        """
        return "\n".join(
            indent + line for line in text.splitlines(keepends=True)
        )

    def _format_action_invocation(self, action: argparse.Action) -> str:
        """
        Change how 'metavar' is displayed

        https://stackoverflow.com/questions/23936145/python-argparse-help-message-disable-metavar-for-short-options
        """
        if not action.option_strings:
            return self._metavar_formatter(action, action.dest)(1)[0]
        # option takes no arguments -> -s, --long
        # option takes arguments:
        #    default output -> -s ARGS, --long ARGS
        #    changed output -> -s, --long type
        if action.type is None:
            raise Exception(f"action type is None for some reason: {action}")
        if isinstance(action.type, (str, argparse.FileType)):
            raise Exception(f"unsupported action type: {action}")
        return ", ".join(action.option_strings) + (
            f" {action.type.__name__}" if action.nargs != 0 else ""
        )

    # add default value to the end
    # built-in argparse class argparse.ArgumentDefaultsHelpFormatter
    # def _get_help_string(self, action):
    #    if action.nargs == 0 or action.default is SUPPRESS or not action.default:
    #        return action.help
    #    return f'{action.help}\n[default: %(default)s]'


class _ArgparseArgumentGroup(argparse._ArgumentGroup):

    def add_arguments(
        self, **arguments: dict[str, typing.Any]
    ) -> "_ArgparseArgumentGroup":
        if self.title is None:
            raise Exception("missing title for some reason")
        group = self.title.split()[0].lower()
        for name, kwargs in arguments.items():
            flags = f"-{group[0]}{name[0]}", f"--{group}-{name}"
            self.add_argument(*flags, **kwargs)
        return self


class ArgparseParser(argparse.ArgumentParser):

    def add_argument_group(
        self, *args, parameters: typing.Optional[ dict[str,dict[str, typing.Any]]] = None,
        **kwargs
    ):
        """
        Override
        """
        group = _ArgparseArgumentGroup(self, *args, **kwargs)
        self._action_groups.append(group)
        if parameters is not None:
            group.add_arguments(**parameters)
        return group

    def add_arguments(self, **groups:dict[str, typing.Any]) -> "ArgparseParser":
        """
        Add a bunch of arguments and argument groups in one go
        """
        for name, parameters in groups.items():
            self.add_argument_group(
                title=name, description="", parameters=parameters
            )
        return self

    def convert_arg_line_to_args(self, arg_line: str) -> list[str]:
        """
        Change how '@' files are parsed

        Allows to have several arguments on one line

        https://docs.python.org/3/library/argparse.html#argparse.ArgumentParser.convert_arg_line_to_args
        """
        return arg_line.split()

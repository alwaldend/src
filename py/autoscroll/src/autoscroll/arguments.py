from argparse import SUPPRESS, ArgumentParser, HelpFormatter, _ArgumentGroup
from typing import Any, Dict, List


def parse_arguments(**arguments: Any) -> Dict[str, Any]:
    result = {}
    for key, value in arguments.items():
        key_split = key.split("_")
        group = key_split[0]
        name = "_".join(key_split[1:])
        if group not in result:
            result[group] = {}
        result[group][name] = value
    return result


class ArgparseFormatter(HelpFormatter):

    # do not split lines that start with 'R|'
    # https://stackoverflow.com/questions/3853722/how-to-insert-newlines-on-argparse-help-text
    def _split_lines(self, text: str, width: int):
        if not text.startswith("R|"):
            return super()._split_lines(text, width)
        result = []
        for line in text[2:].splitlines(keepends=True):
            result.extend(super()._split_lines(line, width))
        return result

    # do not format the description
    # built-in argparse class argparse.RawDescriptionHelpFormatter
    def _fill_text(self, text: str, width: int, indent: str):
        return "\n".join(
            indent + line for line in text.splitlines(keepends=True)
        )

    # change how 'metavar' is displayed
    # https://stackoverflow.com/questions/23936145/python-argparse-help-message-disable-metavar-for-short-options
    def _format_action_invocation(self, action):
        if not action.option_strings:
            return self._metavar_formatter(action, action.dest)(1)[0]
        # option takes no arguments -> -s, --long
        # option takes arguments:
        #    default output -> -s ARGS, --long ARGS
        #    changed output -> -s, --long type
        return ", ".join(action.option_strings) + (
            f" {action.type.__name__}" if action.nargs != 0 else ""
        )

    # add default value to the end
    # built-in argparse class argparse.ArgumentDefaultsHelpFormatter
    # def _get_help_string(self, action):
    #    if action.nargs == 0 or action.default is SUPPRESS or not action.default:
    #        return action.help
    #    return f'{action.help}\n[default: %(default)s]'


class _ArgparseArgumentGroup(_ArgumentGroup):

    def add_arguments(
        self, **arguments: Dict[str, Any]
    ) -> "_ArgparseArgumentGroup":
        group = self.title.split()[0].lower()
        for name, kwargs in arguments.items():
            flags = f"-{group[0]}{name[0]}", f"--{group}-{name}"
            self.add_argument(*flags, **kwargs)
        return self


class ArgparseParser(ArgumentParser):

    # override
    def add_argument_group(
        self, *args, parameters: Dict[str, Dict[str, Any]] = None, **kwargs
    ):
        group = _ArgparseArgumentGroup(self, *args, **kwargs)
        self._action_groups.append(group)
        if parameters is not None:
            group.add_arguments(**parameters)
        return group

    # add a bunch of arguments and argument groups in one go
    def add_arguments(self, **groups: Dict[str, Any]) -> "ArgparseParser":
        for name, parameters in groups.items():
            self.add_argument_group(
                title=name, description="", parameters=parameters
            )
        return self

    # change how '@' files are parsed
    # allows to have several arguments on one line
    # https://docs.python.org/3/library/argparse.html#argparse.ArgumentParser.convert_arg_line_to_args
    def convert_arg_line_to_args(self, arg_line: str) -> List[str]:
        return arg_line.split()

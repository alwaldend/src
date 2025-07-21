import collections
import importlib.resources
import os.path
import typing


def construct_argument(argument: tuple[str, dict[str, str]]):
    help = "\n".join(
        [
            f"{key} ({value[0]}, default - {value[1]})"
            for key, value in argument[1].items()
        ]
    )
    return (
        (f"-{argument[0][0]}", f"--{argument[0]}"),
        {"help": f"R|{help}", "nargs": 2, "action": "append"},
    )


def convert_bool(value: typing.Optional[typing.Union[str, bool]] = None) -> bool:
    if value is None:
        return False
    return value if isinstance(value, bool) else value.lower() in ("true", "+")


def return_none(*args, **kwargs) -> None:
    return


def return_kwargs(**kwargs) -> dict[typing.Any, typing.Any]:
    return kwargs


def has_dict(*values: typing.Any) -> bool:
    return False in (hasattr(value, "__dict__") for value in values)


def convert(value: typing.Any) -> dict[str, typing.Any]:
    if value is None:
        return {}
    if isinstance(value, dict) and not isinstance(value, str):
        return value
    return raise_type_error(value, (type(None), typing.Iterable))


def construct(item: typing.Union[str, tuple[str, typing.Any]]) -> dict[str, typing.Any]:
    return {item: True} if isinstance(item, str) else {item[0]: item[1]}


def parse_argument(argument: list[tuple[str, typing.Any]]) -> dict[str, typing.Any]:
    return dict(collections.ChainMap(*[construct(item) for item in convert(argument)]))


def parse_arguments_old(arguments: dict[str, typing.Any]) -> dict[str, typing.Any]:
    return {name: parse_argument(value) for name, value in arguments.items()}


def construct_coordinates(x: int, y: int) -> str:
    return f"{{{x}, {y}}}"


def documented_by(original):
    def wrapper(target):
        target.__doc__ = original.__doc__
        return target

    return wrapper


def raise_type_error(value: typing.Any, target_type: typing.Any) -> typing.Any:
    if isinstance(value, target_type):
        return value
    message = f"value {str(value)} is not {target_type}, it is {type(value)}"
    raise TypeError(message)


def check_iterable(value: typing.Sequence, length: int = 2) -> typing.Sequence:
    raise_type_error(value, typing.Sequence)
    if len(value) >= length:
        return value
    message = f'value "{str(value)}" does not have the length of {str(length)}'
    raise ValueError(message)


def get_path(path: str) -> str:
    return path if os.path.isfile(path) else get_resource_path(path)


def get_resource_path(resource: typing.Optional[str] = None) -> str:
    if resource is None:
        file = ""
        addition = ""
    else:
        split = resource.split("/")
        addition = ".".join(split[0:-1])
        addition = f".{addition}" if addition else ""
        file = split[-1]
        file = file if file != addition else ""
    with importlib.resources.path(f"autoscroll{addition}", file) as path:
        return str(path)


def get_resource_content(resource: str) -> str:
    with open(get_resource_path(resource), "r") as file:
        return file.read()

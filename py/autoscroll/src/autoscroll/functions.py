from collections import ChainMap
from importlib.resources import path as importlib_path
from os.path import isfile as os_isfile
from typing import Any, Dict, Iterable, List, Tuple, Union


def construct_argument(argument: Tuple[str, Dict[str, str]]):
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


def convert_bool(value: Union[str, bool] = None) -> bool:
    return value if isinstance(value, bool) else value.lower() in ("true", "+")


def return_none(*args, **kwargs) -> None:
    return


def return_kwargs(**kwargs) -> Dict[Any, Any]:
    return kwargs


def has_dict(*values: Any) -> bool:
    return False in (hasattr(value, "__dict__") for value in values)


def convert(value: Any) -> Dict[str, Any]:
    if value is None:
        return {}
    if isinstance(value, Iterable) and not isinstance(value, str):
        return value
    return raise_type_error(value, (type(None), Iterable))


def construct(item: Union[str, Tuple[str, Any]]) -> Dict[str, Any]:
    return {item: True} if isinstance(item, str) else {item[0]: item[1]}


def parse_argument(argument: List[Tuple[str, Any]]) -> Dict[str, Any]:
    return dict(ChainMap(*[construct(item) for item in convert(argument)]))


def parse_arguments_old(arguments: Dict[str, Any]) -> Dict[str, Any]:
    return {name: parse_argument(value) for name, value in arguments.items()}


def construct_coordinates(x: int, y: int) -> str:
    return f"{{{x}, {y}}}"


def documented_by(original):
    def wrapper(target):
        target.__doc__ = original.__doc__
        return target

    return wrapper


def raise_type_error(value: Any, target_type: Any) -> Any:
    if isinstance(value, target_type):
        return value
    message = f"value {str(value)} is not {target_type}, it is {type(value)}"
    raise TypeError(message)


def check_iterable(value: Iterable, length: int = 2) -> Iterable:
    raise_type_error(value, Iterable)
    if len(value) >= length:
        return value
    message = f'value "{str(value)}" does not have the length of {str(length)}'
    raise ValueError(message)


def get_path(path: str) -> str:
    return path if os_isfile(path) else get_resource_path(path)


def get_resource_path(resource: str = None) -> str:
    if resource is None:
        file = ""
        addition = ""
    else:
        split = resource.split("/")
        addition = ".".join(split[0:-1])
        addition = f".{addition}" if addition else ""
        file = split[-1]
        file = file if file != addition else ""
    return str(importlib_path(f"autoscroll{addition}", file))


def get_resource_content(resource: str) -> str:
    with open(get_resource_path(resource), "r") as file:
        return file.read()

import argparse
import pathlib
import shutil
import sys
import tempfile
import typing

SECTION_START = "START"
SECTION_END = "END"

def replace_section(lines: typing.Iterable[str], section: str, replacement: str) -> typing.Generator[str, None, None]:
    start = f"{section} {SECTION_START}"
    end = f"{section} {SECTION_END}"
    found_start = False
    for line in lines:
        start_present = start in line
        end_present = end in line
        if start_present and end_present:
            raise Exception(f"both start '{start}' and end '{end}' are present")
        elif start_present and found_start:
            raise Exception(f"duplicate start: {start}")
        elif end_present and not found_start:
            raise Exception(f"end '{end}' before the start '{start}'")
        elif start_present:
            found_start = True
            yield line
            continue
        elif end_present:
            yield replacement
            yield "\n"
            yield line
            found_start = False
            continue
        elif found_start:
            continue
        else:
            yield line

def replace_section_in_file(file_path: str, section: str, replacement: str) -> None:
    file = pathlib.Path(file_path)
    with tempfile.TemporaryDirectory() as temp_dir:
        temp_file = pathlib.Path(temp_dir, "temp.txt")
        with temp_file.open("a") as temp_file_obj:
            with file.open("r") as file_obj:
                for line in replace_section(file_obj, section, replacement):
                    temp_file_obj.write(line) 
        shutil.move(temp_file, file)

def output_replacement_in_file(file_path: str, section: str, replacement: str) -> None:
    with open(file_path, "r") as file_obj:
        for line in replace_section(file_obj, section, replacement):
            print(line, end="")

def main(args: typing.Sequence[str] = ()) -> int:
    parser = argparse.ArgumentParser(
        description="Replace a section in a file",
    )
    parser.add_argument("-s", "--section", required=True, help="section name")
    parser.add_argument("-i", "--in-place", action=argparse.BooleanOptionalAction, help="edit files in-place")
    parser.add_argument("-r", "--replacement", required=True, help="section replacement")
    parser.add_argument("files", nargs="+", help="file paths")
    args_parsed = parser.parse_args(args)
    for file in args_parsed.files:
        if args_parsed.in_place:
            replace_section_in_file(file, args_parsed.section, args_parsed.replacement)
        else:
            output_replacement_in_file(file, args_parsed.section, args_parsed.replacement)
    return 0

if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

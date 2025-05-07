import sys
import shutil
import os
import tempfile
import argparse
import typing
import pathlib
import tarfile

def install_directory(source_dir: pathlib.Path, target_dir: pathlib.Path) -> None:
    for cur_path_str, _, files in  os.walk(source_dir):
        cur_path = pathlib.Path(cur_path_str)
        rel_path = cur_path.relative_to(source_dir)
        for file in files:
            source_path = cur_path / file
            target_path = target_dir / rel_path / file
            print(f"{rel_path / file} -> {target_path}")
            target_path.parent.mkdir(parents=True, exist_ok=True)
            target_path.unlink(missing_ok=True)
            shutil.copyfile(source_path ,target_path)


def install_archive(archive: pathlib.Path, target_dir: pathlib.Path) -> None:
    target_dir.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory() as temp_dir:
        with tarfile.open(archive) as archive_obj:
            archive_obj.extractall(temp_dir)
        install_directory(pathlib.Path(temp_dir), target_dir)

def main(argv: typing.Sequence[str] = ()) -> int:
    parser = argparse.ArgumentParser(
        description="Unpack an archive into a directory",
    )
    parser.add_argument("--type", default="auto", help="installation type",)
    parser.add_argument("--target", default=os.environ["HOME"], help="installation target",)
    parser.add_argument("archives", nargs="+", help="Files to install")
    args = parser.parse_args(argv)
    for archive in args.archives:
        install_archive(pathlib.Path(archive), pathlib.Path(args.target))
    return 0

if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

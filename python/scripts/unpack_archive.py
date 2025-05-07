import sys
import argparse
import typing
import pathlib
import tarfile

def unpack_archive(archive: pathlib.Path, directory: pathlib.Path) -> None:
    directory.mkdir(parents=True, exist_ok=True)
    with tarfile.open(archive) as archive_obj:
        archive_obj.extractall(directory)


def main(argv: typing.Sequence[str] = ()) -> int:
    parser = argparse.ArgumentParser(
        description="Unpack an archive into a directory",
    )
    parser.add_argument("archives", nargs="+", help="archive paths")
    parser.add_argument("target", nargs=1, help="target directory")
    args = parser.parse_args(argv)
    for archive in args.archives:
        unpack_archive(pathlib.Path(archive), pathlib.Path(args.target[0]))
    return 0

if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

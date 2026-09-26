"""Create a NoCloud seed ISO using the pinned pycdlib implementation."""

import sys
from importlib.metadata import version

import pycdlib


def main() -> None:
    if sys.argv[1:] == ["--version"]:
        print(f"pycdlib {version('pycdlib')}")
        return
    image = pycdlib.PyCdlib()
    image.new(
        interchange_level=3, joliet=3, rock_ridge="1.09", vol_ident="cidata"
    )
    for source, name in zip(
        sys.argv[2:], ["user-data", "meta-data"], strict=True
    ):
        image.add_file(
            source,
            iso_path=f"/{name.upper().replace('-', '_')};1",
            rr_name=name,
            joliet_path=f"/{name}",
        )
    image.write(sys.argv[1])
    image.close()


if __name__ == "__main__":
    main()

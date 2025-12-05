import argparse
import subprocess
import sys
import typing


def build(ctx: argparse.Namespace) -> None:
    subprocess.run(
        (
            ctx.cmake,
            "-G",
            "Unix Makefiles",
            "-S",
            f"{ctx.src_dir}",
            "-B",
            f"{ctx.build_dir}",
            f"-DCMAKE_BUILD_TYPE={ctx.build_type}",
            f"-DARM_NONE_EABI_TOOLCHAIN_PATH={ctx.gcc_dir}",
            f"-DNRF5_SDK_PATH={ctx.nrfsdk_dir}",
            "-DBUILD_DFU=1",
            "-DBUILD_RESOURCES=1",
        ),
        check=True,
    )
    subprocess.run(
        (
            ctx.cmake,
            "--build",
            f"{ctx.build_dir}",
            "--config",
            f"{ctx.build_type}",
            "--target",
            "pinetime-mcuboot-app",
        ),
        check=True,
    )


def main(argv: typing.Sequence[str]) -> None:
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers()
    build = subparsers.add_parser(name="build", help="Build")
    build.add_argument(
        "cmake", required=True, help="Path to the cmake executable"
    )
    build.add_argument("src_dir", required=True, help="Source directory")
    build.add_argument("build_dir", required=True, help="Build directory")
    build.add_argument("gcc_dir", required=True, help="Gcc directory")
    build.add_argument("build_type", default="Release", help="Build type")
    build.set_defaults(func=build)
    args = parser.parse_args(argv)
    args.func(args)


if __name__ == "__main__":
    main(sys.argv)

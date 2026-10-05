import argparse

import uvicorn

from projects.mlops_2026.internal.recommender import application


def main() -> None:
    parser: argparse.ArgumentParser = argparse.ArgumentParser(
        description="Run recommender"
    )
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=8000)
    args: argparse.Namespace = parser.parse_args()
    application.course_data_paths()
    uvicorn.run(
        "projects.mlops_2026.internal.recommender.application:create_app",
        factory=True,
        host=args.host,
        port=args.port,
    )


if __name__ == "__main__":
    main()

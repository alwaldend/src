"""Serve independent app instances for the routing isolation E2E scenario."""

import argparse
import asyncio
import gc
import json
import os
import pathlib

import fastapi
import uvicorn

from projects.mlops_2026.internal.recommender import application, service


class ObservedService(service.RecommendationService):
    def __init__(self, model_path: pathlib.Path, events: list[str]) -> None:
        self._events: list[str] = events
        super().__init__(model_path)

    def close(self) -> None:
        super().close()
        self._events.append("closed")


async def verify_service_lifecycle() -> None:
    events: list[str] = []
    recommendations: ObservedService = ObservedService(
        application.model_path(), events
    )
    async with recommendations:
        assert len(recommendations.recommend("Time", 1).recommendations) == 1
    assert events == ["closed"]
    try:
        recommendations.recommend("Time", 1)
    except RuntimeError:
        pass
    else:
        raise AssertionError("The exited service retained its model")
    recommendations.close()
    assert events == ["closed", "closed"]
    del recommendations
    gc.collect()
    assert events == ["closed", "closed", "closed"]
    output: pathlib.Path = pathlib.Path(
        os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
    )
    output.joinpath("service-lifecycle.json").write_text(
        json.dumps(
            {
                "context_shutdown": "passed",
                "repeat_close": "passed",
                "deletion_cleanup": "passed",
            }
        )
        + "\n"
    )


def only_first() -> dict[str, str]:
    return {"instance": "first"}


def main() -> None:
    asyncio.run(verify_service_lifecycle())
    parser: argparse.ArgumentParser = argparse.ArgumentParser()
    parser.add_argument("--host", required=True)
    parser.add_argument("--port", required=True, type=int)
    args: argparse.Namespace = parser.parse_args()
    first: fastapi.FastAPI = application.create_app()
    second: fastapi.FastAPI = application.create_app()
    first.add_api_route("/only-first", only_first)
    parent: fastapi.FastAPI = application.create_app()
    parent.mount("/first", first)
    parent.mount("/second", second)
    uvicorn.run(parent, host=args.host, port=args.port)


if __name__ == "__main__":
    main()

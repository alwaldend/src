"""Retain task outcomes without module arguments, values, or command output."""

import json
import os

from ansible.plugins.callback import CallbackBase


class CallbackModule(CallbackBase):
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = "aggregate"
    CALLBACK_NAME = "events"
    CALLBACK_NEEDS_ENABLED = True

    def _record(self, result, outcome):
        event = {
            "step": int(os.environ.get("MOLECULE_RUNNER_STEP", "0")),
            "phase": os.environ.get("MOLECULE_RUNNER_PHASE", "unknown"),
            "task_id": str(result._task._uuid),
            "action": result._task.action,
            "host": result._host.get_name(),
            "outcome": outcome,
            "changed": bool(result._result.get("changed", False)),
        }
        if isinstance(result._result.get("rc"), int):
            event["return_code"] = result._result["rc"]
        if outcome in {"failed", "ignored_failure", "unreachable"}:
            module = result._task.action.rsplit(".", 1)[-1]
            event["failure_kind"] = (
                "connection"
                if outcome == "unreachable"
                else "assertion"
                if module in {"assert", "fail"}
                else "package_operation"
                if module in {"package", "dnf", "dnf5", "apt"}
                else "task"
            )
        with open(
            os.environ["MOLECULE_EVENTS_FILE"], "a", encoding="utf-8"
        ) as stream:
            stream.write(json.dumps(event) + "\n")

    def v2_runner_on_ok(self, result):
        self._record(result, "ok")

    def v2_runner_on_failed(self, result, ignore_errors=False):
        self._record(result, "ignored_failure" if ignore_errors else "failed")

    def v2_runner_on_unreachable(self, result):
        self._record(result, "unreachable")

    def v2_runner_on_skipped(self, result):
        self._record(result, "skipped")

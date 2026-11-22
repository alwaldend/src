"""Public MLflow Projects backend for Vast.ai execution."""

from .backend import VastAIBackend

__all__: tuple[str, ...] = ("VastAIBackend",)
"""Public API for MLflow plugin discovery."""

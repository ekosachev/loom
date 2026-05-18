from pathlib import Path
from typing import Any
import yaml

from loom.models.model import Model


MODELS_PATH = Path.home() / ".loom" / "models.yaml"


class ModelManager:
    def __init__(self):
        MODELS_PATH.parent.mkdir(exist_ok=True)

        if not MODELS_PATH.exists():
            with open(MODELS_PATH, "w") as f:
                yaml.dump({}, f)

    def load_models(self) -> list[Model]:
        with open(MODELS_PATH, "r") as f:
            models: list[dict[str, Any]] = yaml.safe_load(f) or []
            return [Model.model_validate(m) for m in models]

    def save_models(self, models: list[Model]):
        with open(MODELS_PATH, "w") as f:
            yaml.safe_dump(
                [m.dict() for m in models], f, allow_unicode=True, sort_keys=True
            )

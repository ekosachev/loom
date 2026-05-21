from typing import Any


class BaseTool:
    name: str
    description: str

    @property
    def schema(self) -> dict[str, Any]:
        raise NotImplementedError

from abc import ABC, abstractmethod
from typing import AsyncGenerator
from loom.models.message import Message
from loom.models.model import Model
from loom.models.stream import StreamChunk


class BaseProvider(ABC):
    @abstractmethod
    def chat_completion(
        self, messages: list[Message], model_id: str
    ) -> AsyncGenerator[StreamChunk, None]:
        pass

    @abstractmethod
    async def get_model_meta(self, slug: str) -> Model:
        pass

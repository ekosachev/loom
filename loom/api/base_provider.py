from abc import ABC, abstractmethod
from typing import AsyncGenerator
from loom.models.message import Message
from loom.models.stream import StreamChunk


class BaseProvider(ABC):
    @abstractmethod
    def chat_completion(
        self, messages: list[Message]
    ) -> AsyncGenerator[StreamChunk, None]:
        pass

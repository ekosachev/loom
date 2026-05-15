from abc import ABC, abstractmethod
from typing import AsyncIterator
from loom.models.message import Message
from loom.models.stream import StreamChunk


class BaseProvider(ABC):
    @abstractmethod
    def chat_completion(self, messages: list[Message]) -> AsyncIterator[StreamChunk]:
        pass

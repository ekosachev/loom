from abc import ABC, abstractmethod
from loom.models.message import Message


class BaseProvider(ABC):
    @abstractmethod
    async def chat_completion(self, messages: list[Message]) -> str:
        pass

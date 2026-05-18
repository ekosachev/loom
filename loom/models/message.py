from pydantic import BaseModel
from typing import Optional


class Message(BaseModel):
    role: str
    content: str


class UserMessage(Message):
    role: str = "user"
    name: str


class AssistantMessage(Message):
    role: str = "assistant"
    reasoning: Optional[str] = None


class SystemMessage(Message):
    role: str = "system"

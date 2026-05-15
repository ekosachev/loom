from pydantic import BaseModel


class Message(BaseModel):
    role: str
    content: str


class UserMessage(Message):
    role: str = "user"
    name: str


class AssistantMessage(Message):
    role: str = "assistant"
    name: str


class SystemMessage(Message):
    role: str = "system"

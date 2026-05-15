from pydantic import BaseModel

from loom.models.message import Message


class StreamChunk(BaseModel):
    choices: list["CompletionChoice"]


class CompletionChoice(BaseModel):
    delta: Message

from pydantic import BaseModel

from loom.models.message import AssistantMessage


class StreamChunk(BaseModel):
    choices: list["CompletionChoice"]


class CompletionChoice(BaseModel):
    delta: AssistantMessage

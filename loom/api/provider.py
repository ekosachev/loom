from typing import AsyncIterator
from loom.api.base_provider import BaseProvider
import httpx
from loom.models.message import Message
from loom.models.stream import StreamChunk

OPEN_ROUTER_BASE_URL = "https://openrouter.ai/api/v1"

STOP_CHUNK = "data: [DONE]"
PROCESSING_CHUNK = ": OPENROUTER PROCESSING"


class Provider(BaseProvider):
    client: httpx.AsyncClient

    def __init__(self, api_key) -> None:
        self.client = httpx.AsyncClient(
            headers={"Authorization": f"Bearer {api_key}"},
            base_url=OPEN_ROUTER_BASE_URL,
        )

    async def chat_completion(
        self, messages: list[Message]
    ) -> AsyncIterator[StreamChunk]:
        body = {
            "model": "deepseek/deepseek-v4-flash:free",
            "stream": True,
            "messages": [msg.model_dump() for msg in messages],
        }

        async with self.client.stream(
            "POST", "/chat/completions", json=body
        ) as response:
            async for chunk in response.aiter_lines():
                if not chunk or chunk == PROCESSING_CHUNK:
                    continue
                if chunk == STOP_CHUNK:
                    break

                data = chunk.removeprefix("data: ")
                yield StreamChunk.model_validate_json(json_data=data)

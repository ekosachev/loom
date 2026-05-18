from typing import AsyncGenerator

from loom.api.base_provider import BaseProvider
import httpx
from loom.models.message import Message
from loom.models.model import Model
from loom.models.stream import StreamChunk

OPEN_ROUTER_BASE_URL = "https://openrouter.ai/api/v1"

STOP_CHUNK = "data: [DONE]"
PROCESSING_CHUNK = ": OPENROUTER PROCESSING"


class Provider(BaseProvider):
    _api_key: str
    client: httpx.AsyncClient

    def __init__(self, key: str) -> None:
        self._api_key = key
        self.client = httpx.AsyncClient(
            headers={"Authorization": f"Bearer {self._api_key}"},
            base_url=OPEN_ROUTER_BASE_URL,
            timeout=None,
        )

    async def chat_completion(
        self, messages: list[Message], model_id: str
    ) -> AsyncGenerator[StreamChunk, None]:
        body = {
            "model": model_id,
            "stream": True,
            "messages": [msg.model_dump() for msg in messages],
        }

        async with self.client.stream(
            "POST", "/chat/completions", json=body
        ) as response:
            async for chunk in response.aiter_lines():
                if not chunk or chunk == PROCESSING_CHUNK or chunk.startswith(":"):
                    continue
                if chunk == STOP_CHUNK:
                    break

                data = chunk.removeprefix("data: ")
                yield StreamChunk.model_validate_json(json_data=data)

    async def get_model_meta(self, slug: str) -> Model:
        response = await self.client.get(f"/models/{slug}/endpoints")
        response.raise_for_status()
        data = response.json().get("data", {})
        return Model(
            name=data["name"],
            slug=data["id"],
        )

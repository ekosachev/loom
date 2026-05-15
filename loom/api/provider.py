from loom.api.base_provider import BaseProvider
import httpx
from loom.models.message import Message

OPEN_ROUTER_BASE_URL = "https://openrouter.ai/api/v1"


class Provider(BaseProvider):
    client: httpx.AsyncClient

    def __init__(self, api_key) -> None:
        self.client = httpx.AsyncClient(
            headers={"Authorization": f"Bearer {api_key}"},
            base_url=OPEN_ROUTER_BASE_URL,
        )

    async def chat_completion(self, messages: list[Message]) -> str:
        body = {
            "model": "deepseek/deepseek-v4-flash:free",
            "stream": True,
            "messages": [msg.model_dump() for msg in messages],
        }

        total = ""

        async with self.client.stream(
            "POST", "/chat/completions", json=body
        ) as response:
            async for chunk in response.aiter_text():
                total += chunk

        return total

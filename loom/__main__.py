from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
from asyncio import run

from loom.models.message import SystemMessage, UserMessage

load_dotenv()

provider = Provider(getenv("OPEN_ROUTER_KEY"))

messages = [
    SystemMessage(content="You are a helpful assistant"),
    UserMessage(
        content="This is a test request. Answer with 2-3 sentences of your choosing",
        name="Egor",
    ),
]


async def run_query():
    async for chunk in provider.chat_completion(messages):
        print(chunk.choices[0].delta.content, end="")


run(run_query())

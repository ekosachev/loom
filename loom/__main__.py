from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
from asyncio import run

from loom.models.message import Message, SystemMessage, UserMessage

load_dotenv()

provider = Provider(getenv("OPEN_ROUTER_KEY"))

messages = [
    SystemMessage(content="You are a helpful assistant"),
    UserMessage(content="What is the capital of France?", name="egor"),
]

text = run(provider.chat_completion(messages))

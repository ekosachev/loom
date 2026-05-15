from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
from asyncio import run

from loom.models.message import SystemMessage, UserMessage, Message
from loom.ui.loom_ui import LoomUI

load_dotenv()

provider = Provider(getenv("OPEN_ROUTER_KEY"))

messages: list[Message] = [
    SystemMessage(content="You are a helpful assistant"),
]

ui = LoomUI()

messages.append(UserMessage(name="egor", content=ui.get_user_input()))

run(ui.consume_stream(provider.chat_completion(messages)))

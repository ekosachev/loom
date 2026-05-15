from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
from asyncio import run

from loom.models.message import SystemMessage, UserMessage
from loom.ui.loom_ui import LoomUI

load_dotenv()

provider = Provider(getenv("OPEN_ROUTER_KEY"))

messages = [
    SystemMessage(content="You are a helpful assistant"),
    UserMessage(
        content="This is a test request. Answer with 2-3 sentences of your choosing",
        name="Egor",
    ),
]

ui = LoomUI()

run(ui.consume_stream(provider.chat_completion(messages)))

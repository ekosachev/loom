from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
from asyncio import run

from loom.database.chat_storage import ChatStorage
from loom.models.message import SystemMessage, UserMessage, Message
from loom.ui.loom_ui import LoomUI


async def main():
    load_dotenv()
    provider = Provider(getenv("OPEN_ROUTER_KEY"))
    ui = LoomUI()
    storage = ChatStorage()
    storage.clear_history()
    storage.add_message(role="system", content="You are a helpful assistant.")

    while True:
        user_text = ui.get_user_input()
        if user_text.lower() == "exit":
            break

        if user_text.lower() == "clear":
            storage.clear_history()
            continue

        storage.add_message(role="user", content=user_text, name="Developer")

        full_context = storage.get_history()

        stream = provider.chat_completion(full_context)
        await ui.consume_stream(stream)

        storage.add_message(role="assistant", content=ui._buffer)


if __name__ == "__main__":
    run(main())

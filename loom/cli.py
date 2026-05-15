import asyncio
from typing import Optional
from rich.console import Console
from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
import typer
from pathlib import Path

from loom.database.chat_storage import ChatStorage
from loom.ui.loom_ui import LoomUI
import yaml


CONFIG_PATH = Path.home() / ".loom" / "config.yaml"

console = Console()


def _get_api_key() -> str:
    load_dotenv()
    key = getenv("OPEN_ROUTER_API_KEY")

    if key is not None:
        return key

    console.print("[dim]Key not in environment, searching in config file...[/dim]")

    if not CONFIG_PATH.exists():
        console.print("[bold red]Error:[/bold red] config file does not exist")
        console.print(
            "Use [bold cyan]loom config[/bold cyan] to set api key, or define OPEN_ROUTER_API_KEY environment variable"
        )
        raise typer.Exit(-1)

    with open(CONFIG_PATH, "r") as f:
        data = yaml.safe_load(f)
        key = data.get("api_key")
        if key is not None:
            return key

    console.print("[bold red]Error:[/bold red] config file does not contain api key")
    console.print(
        "Use [bold cyan]loom config[/bold cyan] to set api key, or define OPEN_ROUTER_API_KEY environment variable"
    )

    raise typer.Exit(-1)


app = typer.Typer(help="Loom: CLI interface for LLMs with git-like context management")
provider: Optional[Provider] = None


def _create_provider():
    global provider

    key = _get_api_key()
    provider = Provider(key)


@app.command(help="Configure loom-cli")
def config(key: str = typer.Option(None, "--key", help="OpenRouter API Key")):
    if not key:
        key = typer.prompt("Enter your OpenRouter API Key", hide_input=True)

    config_data = {"api_key": key}

    with open(CONFIG_PATH, "w") as f:
        yaml.dump(config_data, f)

    console.print("[bold green]Configuration saved to ~/.loom/config.yaml[/bold green]")


@app.command(help="Start a new conversation")
def init():
    storage = ChatStorage()
    storage.clear_history()
    console.print("[bold green]New chat initialized. History cleared.[/bold green]")


@app.command(help="Request a new completion using active conversation")
def send(message: str):
    asyncio.run(_async_send(message))


async def _async_send(message: str):
    global provider

    storage = ChatStorage()
    ui = LoomUI()

    if provider is None:
        _create_provider()

    assert provider is not None

    storage.add_message(role="user", content=message, name="Developer")
    history = storage.get_history()

    stream = provider.chat_completion(history)
    await ui.consume_stream(stream)
    storage.add_message(role="assistant", content=ui._buffer)


def main():
    app()

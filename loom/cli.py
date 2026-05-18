import asyncio
from typing import Optional
import httpx
from rich.console import Console
from loom.api.provider import Provider
from dotenv import load_dotenv
from os import getenv
import typer
from pathlib import Path

from loom.database.chat_storage import ChatStorage
from loom.database.model_manager import ModelManager
from loom.ui.log_viewer import LoomLogViewer
from loom.ui.loom_ui import LoomUI
import yaml


CONFIG_PATH = Path.home() / ".loom" / "config.yaml"

console = Console()
model_manager = ModelManager()
app = typer.Typer(help="Loom: CLI interface for LLMs with git-like context management")

models_app = typer.Typer(help="Managing available LLMs")
app.add_typer(models_app, name="models")

provider: Optional[Provider] = None


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


def _create_provider():
    key = _get_api_key()
    provider = Provider(key)

    return provider


@app.command(help="Configure loom-cli")
def config(key: str = typer.Option(None, "--key", help="OpenRouter API Key")):
    if not key:
        key = typer.prompt("Enter your OpenRouter API Key", hide_input=True)

    config_data = {"api_key": key}

    with open(CONFIG_PATH, "w") as f:
        yaml.dump(config_data, f)

    console.print("[bold green]Configuration saved to ~/.loom/config.yaml[/bold green]")


@app.command(help="Create a new workspace")
def init(name: str):
    storage = ChatStorage()
    storage.create_workspace(name)
    storage.switch_to_workspace(name)
    console.print(f"[bold magenta]WORKSPACE[/bold magenta] New workspace: {name}")
    console.print(
        f"[bold magenta]WORKSPACE[/bold magenta] Switched to workspace: {name}"
    )


@app.command(name="ws", help="Switch to an existing workspace")
def change_workspace(name: str):
    storage = ChatStorage()
    storage.switch_to_workspace(name)
    console.print(
        f"[bold magenta]WORKSPACE[/bold magenta] Switched to workspace: {name}"
    )


@app.command(help="Display current workspace and branch")
def status():
    storage = ChatStorage()
    ws = storage.get_current_workspace()
    console.print(f"Workspace:\t{ws.name}")
    branch = storage.get_current_branch()
    console.print(f"Branch:\t{branch.name}")


@app.command(help="Request a new completion using active conversation")
def send(message: str):
    storage = ChatStorage()
    active_model = storage.get_active_model()
    asyncio.run(_async_send(message, active_model))


@app.command(help="Create or switch to another branch")
def checkout(
    branch_name: str, b: bool = typer.Option(False, "-b", help="Create a new branch")
):
    storage = ChatStorage()
    current_branch = storage.get_current_branch()

    if b:
        storage.create_branch(branch_name)
        console.print(f"Created new branch {branch_name} from {current_branch.name}")

    storage.switch_to_branch(branch_name)
    console.print(f"Switched to branch {branch_name}")


@app.command(help="List available branches")
def branch():
    storage = ChatStorage()
    branches = storage.get_all_branches()
    current = storage.get_current_branch()

    for b in branches:
        prefix = "CUR" if str(b.name) == str(current.name) else ""
        style = "green bold" if str(b.name) == str(current.name) else "white"
        console.print(f"[{style}]{prefix}\t{b.name}[/{style}]")


@app.command(help="List all workspaces")
def workspace():
    storage = ChatStorage()
    workspaces = storage.get_all_workspaces()
    current = storage.get_current_workspace()

    for ws in workspaces:
        prefix = "CUR" if str(ws.name) == str(current.name) else ""
        style = "green bold" if str(ws.name) == str(current.name) else "white"
        console.print(f"[{style}]{prefix}\t{ws.name}[/{style}]")


@app.command(help="Display chat log for current branch")
def log():
    storage = ChatStorage()
    workspace = storage.get_current_workspace()
    branch = storage.get_current_branch()
    messages = storage.get_history()

    viewer = LoomLogViewer(messages, workspace.name, branch.name)
    viewer.run()


@models_app.command(name="list")
def models_list():
    models = model_manager.load_models()

    if not models:
        console.print(
            "[yellow]No models added.[/yellow] Add a new model using [cyan]loom model add [italic]<id>[/italic][/cyan]"
        )
        return

    for model in models:
        console.print(f"[dim]{model.slug}[/dim]\t{model.name}")


@models_app.command(name="add")
def models_add(
    model_id: str = typer.Argument(
        ..., help="OpenRouter model id, e. g. google/gemini-2.0-flash-001"
    ),
):
    provider = _create_provider()
    models = {m.slug: m for m in model_manager.load_models()}

    if model_id in models:
        console.print(
            f"[yellow]Model with id {model_id} is alreay present in config[/yellow]. Configuration will be overwritten"
        )

    try:
        with console.status(
            f"[bold green]Requesting metadata for {model_id}...[/bold green]",
            spinner="dots",
        ):
            model = asyncio.run(provider.get_model_meta(model_id))
    except httpx.HTTPStatusError as e:
        console.print(f"[bold red]Error:[/bold red] {e.__repr__()}")
        return

    models[model_id] = model
    model_manager.save_models(list(models.values()))

    console.print(f"[green bold]{model.name}[/green bold] added successfully!")


@models_app.command(name="set")
def models_set_active(
    model_id: str = typer.Argument(
        ..., help="OpenRouter model id, e. g. google/gemini-2.0-flash-001"
    ),
):
    storage = ChatStorage()
    models = {m.slug: m for m in model_manager.load_models()}

    if model_id not in models:
        console.print(
            f"[yellow]Model not found.[/yellow] Add a new model using [cyan]loom model add [italic]{model_id}[/italic][/cyan]"
        )
        return

    storage.set_active_model(model_id)


async def _async_send(message: str, model_id: str):
    storage = ChatStorage()
    ui = LoomUI()

    provider = _create_provider()

    storage.add_message(role="user", content=message, name="Developer")
    history = storage.get_history()

    stream = provider.chat_completion(history, model_id)
    await ui.consume_stream(stream)
    storage.add_message(
        role="assistant", content=ui._buffer, reasoning=ui._reasoning_buffer
    )


def main():
    app()

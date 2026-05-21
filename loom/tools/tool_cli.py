from rich.console import Console
import typer

from loom.tools.tool_manager import ToolManager


tools_app = typer.Typer()
manager = ToolManager()
console = Console()


@tools_app.command(name="list", help="Display list of available tools")
def list_tools():
    for name, tool in manager.tools.items():
        short_description = (
            tool.description[:57] + "..."
            if len(tool.description) > 60
            else tool.description
        )
        console.print(f"[bold]{name}[/bold]\t[dim]{short_description}[/dim]")

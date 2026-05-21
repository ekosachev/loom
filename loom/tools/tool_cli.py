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


@tools_app.command(name="info")
def tool_info(name: str):
    tool = manager.get_tool(name)
    params = manager.get_tool_props(name)

    console.print(f"[bold]{tool.name}[/bold]")
    console.print(tool.description)

    if not params["properties"]:
        console.print("[dim]This tool uses no parameters[/]")
        return

    for name, p in params["properties"].items():
        p_color = "red" if name in params["required"] else "white"
        console.print(
            f"[{p_color}]{name}[/{p_color}]\t{p['type']}\t[dim]{p['description']}[/]"
        )

import click
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.text import Text
from rich.markdown import Markdown

from loom.models.message import Message


console = Console()

ROLE_COLORS = {
    "user": "bright_cyan",
    "assistant": "bright_magenta",
    "system": "yellow",
}


class LoomLogViewer:
    messages: list[Message]
    workspace: str
    branch: str
    selected_index: int = 0
    expanded_indices: set[int] = set()

    def __init__(self, messages: list[Message], workspace: str, branch: str):
        self.messages = messages
        self.workspace = workspace
        self.branch = branch

    def _get_preview(self, content: str) -> str:
        first_line = content.split("\n")[0].strip()
        if len(first_line) > 60:
            return first_line[:57] + "..."
        return first_line

    def _render_screen(self):
        console.clear()

        console.print(
            f"[bold cyan]Workspace:[/bold cyan] {self.workspace} | [bold magenta]Branch:[/bold magenta] {self.branch}"
        )

        console.print("[dim]" + "─" * console.width + "[/dim]\n")

        if not self.messages:
            console.print("[yellow]No messages[/yellow]")
            return

        grid = Table.grid(padding=(0, 1))
        grid.add_column(no_wrap=True)
        grid.add_column(no_wrap=True)
        grid.add_column()

        for idx, msg in enumerate(self.messages):
            is_selected = idx == self.selected_index
            is_expanded = idx in self.expanded_indices

            role_color = ROLE_COLORS.get(msg.role, "white")

            pointer_str = "⇒ " if is_selected else "  "
            pointer = Text(
                pointer_str, style="bold blink white" if is_selected else "white"
            )

            role_text = Text(f"{msg.role.upper():^10}", style=role_color + " bold")

            if is_expanded:
                content_renderable = Panel(
                    Markdown(msg.content),
                    border_style=role_color,
                    padding=(1, 2),
                    expand=True,
                )
            else:
                preview_text = self._get_preview(msg.content)
                content_renderable = Text(preview_text)

            if is_expanded:
                grid.add_row("", "", "")

            grid.add_row(pointer, role_text, content_renderable)

            if is_expanded:
                grid.add_row("", "", "")

        console.print(grid)

        console.print("\n[dim]" + "─" * console.width + "[/dim]")
        console.print(
            "[dim][^/v] to scroll | [Enter] to unfold/fold | [q] to exit[/dim]"
        )

    def run(self):
        if not self.messages:
            console.print("[yellow]No messages in current branch.[/yellow]")
            return

        with console.screen():
            while True:
                self._render_screen()

                key = click.getchar()

                if key in ("q", "Q", "\x1b"):
                    break
                elif key == "\x1b[A":
                    self.selected_index = max(0, self.selected_index - 1)
                elif key == "\x1b[B":
                    self.selected_index = min(
                        len(self.messages) - 1, self.selected_index + 1
                    )
                elif key in ("\r", "\n"):
                    if self.selected_index in self.expanded_indices:
                        self.expanded_indices.remove(self.selected_index)
                    else:
                        self.expanded_indices.add(self.selected_index)

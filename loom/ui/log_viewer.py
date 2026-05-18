import click
from rich.console import Console, Group
from rich.panel import Panel
from rich.table import Table
from rich.text import Span, Text
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
    internal_offsets: dict[int, int] = dict()

    def __init__(self, messages: list[Message], workspace: str, branch: str):
        self.messages = messages
        self.workspace = workspace
        self.branch = branch

    def _get_preview(self, content: str) -> str:
        first_line = content.split("\n")[0].strip()
        if len(first_line) > 60:
            return first_line[:57] + "..."
        return first_line

    def _get_max_panel_lines(self) -> int:
        # terminal height - header (2) - footer (2) - 4 msgs
        return max(5, console.height - 20)

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

        start_window = max(0, self.selected_index - 2)
        end_window = min(len(self.messages), self.selected_index + 3)

        for idx in range(start_window, end_window):
            msg = self.messages[idx]
            is_selected = idx == self.selected_index
            is_expanded = idx in self.expanded_indices

            role_color = ROLE_COLORS.get(msg.role, "white")

            pointer_str = "⇒ " if is_selected else "  "
            pointer = Text(
                pointer_str, style="bold blink white" if is_selected else "white"
            )

            role_text = Text(f"{msg.role.upper():^10}", style=role_color + " bold")

            if is_expanded:
                max_lines = self._get_max_panel_lines()
                width = console.width - 16

                options = console.options.update_width(width)
                all_lines = console.render_lines(Markdown(msg.content), options)
                total_lines = len(all_lines)

                offset = self.internal_offsets.get(idx, 0)
                visible_lines = all_lines[offset : offset + max_lines]

                visible_text = Text()
                for line in visible_lines:
                    for segment in line:
                        visible_text.append(segment.text, style=segment.style)
                    visible_text.append("\n", style=line[-1].style)
                content_group = Group(visible_text)

                content_renderable = Panel(
                    content_group,
                    border_style=role_color,
                    padding=(1, 2),
                    subtitle_align="left",
                    subtitle=f"[dim]Lines {offset + 1}-{min(offset + max_lines, total_lines)} of {total_lines}[/dim]",
                )
            else:
                preview_text = self._get_preview(msg.content)
                content_renderable = Text(preview_text)

            grid.add_row(pointer, role_text, content_renderable)

            if is_expanded:
                grid.add_row("", "", "")

        console.print(grid)

        console.print("\n[dim]" + "─" * console.width + "[/dim]")
        console.print(
            "[dim] j/k to scroll | "
            "Enter to unfold/fold | "
            "q to exit | "
            f"selected {self.selected_index} | "
            f"offset {self.internal_offsets.get(self.selected_index, 0)} | "
            f"expanded {self.selected_index in self.expanded_indices}"
            "[/dim]"
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

                is_expanded = self.selected_index in self.expanded_indices
                offset = self.internal_offsets.get(self.selected_index, 0)

                total_lines = 0
                max_lines = self._get_max_panel_lines()

                if is_expanded:
                    msg = self.messages[self.selected_index]
                    options = console.options.update_width(console.width - 16)
                    total_lines = len(
                        console.render_lines(Markdown(msg.content), options)
                    )

                if key == "k":
                    if is_expanded and offset > 0:
                        self.internal_offsets[self.selected_index] = offset - 1
                    else:
                        self.selected_index = max(0, self.selected_index - 1)
                elif key == "j":
                    if is_expanded and offset + max_lines < total_lines:
                        self.internal_offsets[self.selected_index] = offset + 1
                    else:
                        self.selected_index = min(
                            len(self.messages) - 1, self.selected_index + 1
                        )
                elif key in ("\r", "\n"):
                    if self.selected_index in self.expanded_indices:
                        self.expanded_indices.remove(self.selected_index)
                    else:
                        self.expanded_indices.add(self.selected_index)

                    self.internal_offsets[self.selected_index] = 0

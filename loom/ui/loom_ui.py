from typing import AsyncGenerator
from rich.console import Console
from rich.live import Live
from rich.panel import Panel
from rich.markdown import Markdown

from loom.models.stream import StreamChunk


class LoomUI:
    console = Console()
    _buffer = ""

    async def display_stream(self, stream: AsyncGenerator[StreamChunk, None]):
        with Live(
            self._generate_panel(), console=self.console, refresh_per_second=20
        ) as live:
            async for chunk in stream:
                self._buffer += chunk.choices[0].delta.content
                live.update(self._generate_panel())

        self.console.print("\n[dim]Generation finished[dim]")

    def _generate_panel(self) -> Panel:
        return Panel(
            Markdown(self._buffer),
            title="[bold blue]OpenRouter response[/bold blue]",
            border_style="bright_magenta",
            padding=(1, 2),
        )

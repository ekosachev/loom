import time
from typing import AsyncGenerator
from rich.console import Console
from rich.live import Live
from rich.panel import Panel
from rich.markdown import Markdown
from rich.status import Status
from rich.prompt import Prompt

from loom.models.stream import StreamChunk


class LoomUI:
    console = Console()
    _buffer = ""
    start_time: float
    response_stream: AsyncGenerator[StreamChunk, None]
    ttft: float

    def get_user_input(self) -> str:
        return Prompt.ask("\n[bold cyan]You[/bold cyan]")

    async def consume_stream(self, stream: AsyncGenerator[StreamChunk, None]):
        self.start_time = time.perf_counter()
        self.response_stream = stream
        got_first_token = await self._loading_state()

        if not got_first_token:
            return

        await self._display_stream()

        self.console.print(f"[dim]Done. TTFT: {self.ttft:.2f}s[/dim]")

    async def _loading_state(self) -> bool:
        with Status(
            "[bold green]Thinking...", console=self.console, spinner="dots"
        ) as _:
            try:
                first_chunk = await anext(self.response_stream)
                self.ttft = time.perf_counter() - self.start_time
                self._buffer += first_chunk.choices[0].delta.content
            except StopAsyncIteration:
                self.console.print("[red]Error: Model returned an empty response[/red]")
                return False
            else:
                return True

    async def _display_stream(self):
        with Live(
            self._generate_panel(), console=self.console, refresh_per_second=20
        ) as live:
            async for chunk in self.response_stream:
                self._buffer += chunk.choices[0].delta.content
                live.update(self._generate_panel())

    def _generate_panel(self) -> Panel:
        return Panel(
            Markdown(self._buffer),
            title="[bold blue]OpenRouter response[/bold blue]",
            border_style="bright_magenta",
            padding=(1, 2),
        )

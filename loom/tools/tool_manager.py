from loom.errors import DuplicateToolName
from loom.tools.base_tool import BaseTool
from loom.tools.read_file import ReadFile

DEFAULT_TOOLS: list[BaseTool] = [ReadFile()]


class ToolManager:
    tools: dict[str, BaseTool] = dict()

    def __init__(self) -> None:
        self._load_all_tools()

    def _load_all_tools(self):
        self._load_default_tools()

    def _load_default_tools(self):
        for tool in DEFAULT_TOOLS:
            self._add_tool(tool)

    def _add_tool(self, tool: BaseTool):
        if tool.name in self.tools:
            raise DuplicateToolName(tool.name)

        self.tools[tool.name] = tool

from loom.tools.base_tool import BaseTool


class ReadFile(BaseTool):
    name = "read_file"
    description = "Read contents of a text or source code file. Returns file content with line numbers"

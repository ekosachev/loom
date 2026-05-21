from typing import Any
from loom.tools.base_tool import BaseTool


class ReadFile(BaseTool):
    name = "read_file"
    description = "Read contents of a text or source code file. Returns file content with line numbers"

    @property
    def schema(self) -> dict[str, Any]:
        return {
            "type": "function",
            "funciton": {
                "name": self.name,
                "description": self.description,
                "parameters": {
                    "type": "object",
                    "properties": {
                        "path": {
                            "type": "string",
                            "desctiption": "Path to the file. Can be relative or absolute",
                        },
                        "start_line": {
                            "type": "integer",
                            "description": "Start line number (1-indexed)",
                            "default": 1,
                        },
                        "end_line": {
                            "type": "integer",
                            "description": "End line number (1-indexed). Use -1 for end of file",
                            "default": -1,
                        },
                        "show_line_numbers": {
                            "type": "boolean",
                            "description": "Include line numbers in output",
                            "default": True,
                        },
                    },
                    "required": ["path"],
                },
            },
        }

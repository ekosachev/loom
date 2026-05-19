from typing import Any
from pydantic import BaseModel


class ErrorResponse(BaseModel):
    error: "ErrorMessage"


class ErrorMessage(BaseModel):
    message: str
    code: int
    metadata: dict[str, Any]

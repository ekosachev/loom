from pydantic import BaseModel


class Model(BaseModel):
    slug: str
    name: str

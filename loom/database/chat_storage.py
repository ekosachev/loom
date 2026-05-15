from typing import Optional, overload
from loom.database.db import MessageModel, SessionLocal
from loom.models.message import AssistantMessage, Message, SystemMessage, UserMessage


class ChatStorage:
    session = SessionLocal()

    def _save_message(self, message: MessageModel):
        self.session.add(message)
        self.session.commit()

    def add_message(self, role: str, content: str, name: Optional[str] = None):
        message = MessageModel(role=role, content=content, name=name)
        self._save_message(message)

    def get_history(self, limit: int = 20) -> list[Message]:
        messages = (
            self.session.query(MessageModel)
            .order_by(MessageModel.timestamp.asc())
            .limit(limit)
            .all()
        )

        history = []

        for m in messages:
            if m.role == "user":
                history.append(UserMessage.model_validate(m, from_attributes=True))
            elif m.role == "assistant":
                history.append(AssistantMessage.model_validate(m, from_attributes=True))
            elif m.role == "system":
                history.append(SystemMessage.model_validate(m, from_attributes=True))

        return history

    def clear_history(self):
        self.session.query(MessageModel).delete()
        self.session.commit()

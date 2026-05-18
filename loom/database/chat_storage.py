from typing import Optional
from loom.database.db import MessageModel, SessionLocal
from loom.models.message import AssistantMessage, Message, SystemMessage, UserMessage


class ChatStorage:
    session = SessionLocal()

    def _save_message(self, message: MessageModel):
        self.session.add(message)
        self.session.commit()

    def _map_model_to_message(self, message: MessageModel) -> Message:
        if message.role == "user":
            return UserMessage.model_validate(message, from_attributes=True)
        elif message.role == "assistant":
            return AssistantMessage.model_validate(message, from_attributes=True)
        elif message.role == "system":
            return SystemMessage.model_validate(message, from_attributes=True)

        return Message.model_validate(message, from_attributes=True)

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
            history.append(self._map_model_to_message(m))

        return history

    def clear_history(self):
        self.session.query(MessageModel).delete()
        self.session.commit()

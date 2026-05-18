from typing import Optional
from loom.database import workspace_storage
from loom.database.db import (
    BranchModel,
    MessageModel,
    SessionLocal,
    StateModel,
    WorkspaceModel,
)
from loom.database.workspace_storage import WorkspaceStorage
from loom.database.branch_storage import BranchStorage
from loom.errors import (
    NoCurrentBranch,
)
from loom.models.message import AssistantMessage, Message, SystemMessage, UserMessage


class ChatStorage:
    session = SessionLocal()
    workspace_storage: WorkspaceStorage
    branch_storage: BranchStorage

    def __init__(self):
        self.workspace_storage = WorkspaceStorage()
        self.branch_storage = BranchStorage(self.workspace_storage)

    def _save_message(self, message: MessageModel):
        self.session.add(message)
        self.session.commit()

    def _map_model_to_message(self, message: MessageModel) -> Message:
        if str(message.role) == "user":
            return UserMessage.model_validate(message, from_attributes=True)
        elif str(message.role) == "assistant":
            return AssistantMessage.model_validate(message, from_attributes=True)
        elif str(message.role) == "system":
            return SystemMessage.model_validate(message, from_attributes=True)

        return Message.model_validate(message, from_attributes=True)

    def add_message(self, role: str, content: str, name: Optional[str] = None):
        branch_name = self.branch_storage.get_current_branch()
        if branch_name is None:
            raise NoCurrentBranch()

        branch = self.branch_storage.get_current_branch()
        parent_id = branch.current_message_id

        message = MessageModel(
            role=role, content=content, name=name, parent_id=parent_id
        )
        self._save_message(message)

        branch.current_message_id = message.id
        self.session.commit()

    def get_history(self, limit: int = 20) -> list[Message]:
        branch = self.branch_storage.get_current_branch()

        history = []

        current_id = branch.current_message_id

        while current_id is not None and len(history) <= limit:
            msg = self.session.query(MessageModel).filter_by(id=current_id).first()
            if msg is None:
                raise ValueError(f"Message with id {current_id} not found")

            history.append(self._map_model_to_message(msg))
            current_id = msg.parent_id

        return history

    def switch_to_workspace(self, name: str):
        self.workspace_storage.switch_to_workspace(name)
        self.branch_storage = BranchStorage(self.workspace_storage)

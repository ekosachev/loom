from typing import Optional

from loom.database.db import (
    BranchModel,
    MessageModel,
    SessionLocal,
    WorkspaceModel,
)
from loom.database.workspace_storage import WorkspaceStorage
from loom.database.branch_storage import BranchStorage
from loom.models.message import AssistantMessage, Message, SystemMessage, UserMessage


class ChatStorage:
    workspace_storage: WorkspaceStorage
    branch_storage: BranchStorage

    def __init__(self):
        with SessionLocal() as session:
            self.workspace_storage = WorkspaceStorage(session)
            self.branch_storage = BranchStorage(self.workspace_storage, session)

    def _map_model_to_message(self, message: MessageModel) -> Message:
        if message.role == "user":
            return UserMessage.model_validate(message, from_attributes=True)
        elif message.role == "assistant":
            return AssistantMessage.model_validate(message, from_attributes=True)
        elif message.role == "system":
            return SystemMessage.model_validate(message, from_attributes=True)

        return Message.model_validate(message, from_attributes=True)

    def add_message(
        self,
        role: str,
        content: str,
        name: Optional[str] = None,
        reasoning: Optional[str] = None,
    ):
        with SessionLocal() as session:
            branch = self.branch_storage.get_current_branch(session)

            parent_id = branch.current_message_id

            message = MessageModel(
                role=role,
                content=content,
                name=name,
                parent_id=parent_id,
                reasoning=reasoning,
            )

            session.add(message)
            session.commit()

            branch.current_message_id = message.id
            session.commit()

    def get_history(self, limit: int = 20) -> list[Message]:
        with SessionLocal() as session:
            branch = self.branch_storage.get_current_branch(session)

            history = []

            current_id = branch.current_message_id

            while current_id is not None and len(history) <= limit:
                msg = session.query(MessageModel).filter_by(id=current_id).first()
                if msg is None:
                    raise ValueError(f"Message with id {current_id} not found")

                history.append(self._map_model_to_message(msg))
                current_id = msg.parent_id

        return history[::-1]

    def switch_to_workspace(self, name: str):
        with SessionLocal() as session:
            self.workspace_storage.switch_to_workspace(session, name)
            self.branch_storage = BranchStorage(self.workspace_storage, session)

    def create_workspace(self, name: str):
        with SessionLocal() as session:
            self.workspace_storage.create_workspace(session, name)

    def get_current_workspace(self) -> WorkspaceModel:
        with SessionLocal() as session:
            return self.workspace_storage.get_current_workspace(session)

    def get_current_branch(self) -> BranchModel:
        with SessionLocal() as session:
            return self.branch_storage.get_current_branch(session)

    def create_branch(self, name: str):
        with SessionLocal() as session:
            self.branch_storage.create_branch(session, name)

    def switch_to_branch(self, name: str):
        with SessionLocal() as session:
            self.branch_storage.switch_to_branch(session, name)

    def get_all_workspaces(self):
        with SessionLocal() as session:
            return self.workspace_storage.get_all_workspaces(session)

    def get_all_branches(self):
        with SessionLocal() as session:
            return self.branch_storage.get_all_branches(session)

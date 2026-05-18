from os import name
from typing import Optional
from loom.database.db import BranchModel, MessageModel, SessionLocal, StateModel
from loom.errors import (
    BranchAlreadyExists,
    BranchDoesNotExist,
    NoCurrentBranch,
)
from loom.models.message import AssistantMessage, Message, SystemMessage, UserMessage


class ChatStorage:
    session = SessionLocal()

    def __init__(self):
        self._ensure_initial_state()

    def _ensure_initial_state(self):
        if not self.session.query(StateModel).filter_by(key="HEAD").first():
            self.session.add(StateModel(key="HEAD", value="main"))
            self.session.add(BranchModel(name="main", current_message_id=None))
            self.session.commit()

    def get_current_branch(self) -> Optional[str]:
        branch = self.session.query(StateModel).filter_by(key="HEAD").first()
        return branch.value if branch is not None else None

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
        branch_name = self.get_current_branch()
        if branch_name is None:
            raise NoCurrentBranch()

        branch = self.session.query(BranchModel).filter_by(name=branch_name).first()
        parent_id = branch.current_message_id

        message = MessageModel(role=role, content=content, name=name)
        self._save_message(message)

        branch.current_message_id = message.id
        self.session.commit()

    def get_history(self, limit: int = 20) -> list[Message]:
        branch_name = self.get_current_branch()
        if branch_name is None:
            raise NoCurrentBranch()

        branch = self.session.query(BranchModel).filter_by(name=branch_name).first()

        history = []

        current_id = branch.current_message_id

        while current_id is not None and len(history) <= limit:
            msg = self.session.query(MessageModel).filter_by(id=current_id).first()

            history.append(self._map_model_to_message(msg))
            current_id = msg.parent_id

        return history

    def switch_to_branch(self, branch_name: str):
        branch = self.session.query(BranchModel).filter_by(name=branch_name).first()
        if branch is None:
            raise BranchDoesNotExist(branch_name)

        head = self.session.query(StateModel).filter_by(key="HEAD").first()
        head.value = branch.name

        self.session.commit()

    def create_branch(self, branch_name: str):
        existing_branch = (
            self.session.query(BranchModel).filter_by(name=branch_name).first()
        )
        if existing_branch is not None:
            raise BranchAlreadyExists(branch_name)

        current_branch_name = self.get_current_branch()
        if current_branch_name is None:
            raise NoCurrentBranch()
        current_branch = (
            self.session.query(BranchModel).filter_by(name=current_branch_name).first()
        )
        assert current_branch is not None

        branch = BranchModel(
            name=branch_name, current_message_id=current_branch.current_message_id
        )
        self.session.add(branch)
        self.session.commit()

    def get_all_branches(self) -> list[BranchModel]:
        return self.session.query(BranchModel).all()

    def clear_history(self):
        self.session.query(BranchModel).delete()
        self.session.query(MessageModel).delete()
        self.session.commit()

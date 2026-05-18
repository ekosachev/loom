from sqlalchemy.orm import Session
from loom.database.db import BranchModel, StateModel
from loom.database.workspace_storage import WorkspaceStorage
from loom.errors import BranchAlreadyExists, BranchDoesNotExist, NoCurrentBranch


class BranchStorage:
    workspace_storage: WorkspaceStorage

    def __init__(self, workspace_storage: WorkspaceStorage, session: Session):
        self.workspace_storage = workspace_storage

        try:
            _ = self.get_current_branch(session)
            return
        except NoCurrentBranch:
            pass

        workspace = self.workspace_storage.get_current_workspace(session)
        try:
            main = self.get_branch(session, "main")
        except BranchDoesNotExist:
            main = BranchModel(name="main", workspace_id=workspace.name)
            session.add(main)
            session.commit()

        current_branch = StateModel(key=f"{workspace.name}/HEAD", value=main.name)
        session.add(current_branch)
        session.commit()

    def get_branch(self, session: Session, name: str) -> BranchModel:
        ws = self.workspace_storage.get_current_workspace(session)
        branch = (
            session.query(BranchModel)
            .filter_by(name=name, workspace_id=ws.name)
            .first()
        )
        if branch is None:
            raise BranchDoesNotExist(name)
        return branch

    def get_current_branch(self, session: Session) -> BranchModel:
        ws = self.workspace_storage.get_current_workspace(session)
        branch = session.query(StateModel).filter_by(key=f"{ws.name}/HEAD").first()

        if branch is None:
            raise NoCurrentBranch()
        return self.get_branch(session, branch.value)

    def switch_to_branch(self, session: Session, branch_name: str):
        ws = self.workspace_storage.get_current_workspace(session)
        branch = self.get_branch(session, branch_name)

        current_branch = (
            session.query(StateModel).filter_by(key=f"{ws.name}/HEAD").first()
        )
        if current_branch is None:
            raise NoCurrentBranch()

        current_branch.value = branch.name
        session.commit()

    def create_branch(self, session: Session, branch_name: str):
        try:
            self.get_branch(session, branch_name)
        except BranchDoesNotExist:
            pass
        else:
            raise BranchAlreadyExists(branch_name)

        current_branch = self.get_current_branch(session)
        ws = self.workspace_storage.get_current_workspace(session)
        branch = BranchModel(
            name=branch_name,
            current_message_id=current_branch.current_message_id,
            workspace_id=ws.name,
        )

        session.add(branch)
        session.commit()

    def get_all_branches(self, session: Session) -> list[BranchModel]:
        current_workspace = self.workspace_storage.get_current_workspace(session)
        branch = (
            session.query(BranchModel)
            .filter_by(workspace_id=current_workspace.name)
            .all()
        )
        return branch

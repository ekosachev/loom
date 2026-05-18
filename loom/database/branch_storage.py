from loom.database.db import BranchModel, SessionLocal, StateModel
from loom.database.workspace_storage import WorkspaceStorage
from loom.errors import BranchAlreadyExists, BranchDoesNotExist, NoCurrentBranch


class BranchStorage:
    session = SessionLocal()
    workspace_storage: WorkspaceStorage

    def __init__(self, workspace_storage: WorkspaceStorage):
        self.workspace_storage = workspace_storage

        try:
            _ = self.get_current_branch()
            return
        except NoCurrentBranch:
            pass

        workspace = self.workspace_storage.get_current_workspace()
        try:
            main = self.get_branch("main")
        except BranchDoesNotExist:
            main = BranchModel(name="main", workspace_id=workspace.name)
            self.session.add(main)
            self.session.commit()

        current_branch = StateModel(key=f"{workspace.name}/HEAD", value=main.name)
        self.session.add(current_branch)
        self.session.commit()

    def get_branch(self, name: str) -> BranchModel:
        ws = self.workspace_storage.get_current_workspace()
        branch = (
            self.session.query(BranchModel)
            .filter_by(name=name, workspace_id=ws.name)
            .first()
        )
        if branch is None:
            raise BranchDoesNotExist(name)
        return branch

    def get_current_branch(self) -> BranchModel:
        ws = self.workspace_storage.get_current_workspace()
        branch = self.session.query(StateModel).filter_by(key=f"{ws.name}/HEAD").first()
        if branch is None:
            raise NoCurrentBranch()
        return self.get_branch(str(branch.value))

    def switch_to_branch(self, branch_name: str):
        branch = self.get_branch(branch_name)

        head = self.get_current_branch()
        head.value = branch.name

        self.session.commit()

    def create_branch(self, branch_name: str):
        try:
            self.get_branch(branch_name)
        except BranchDoesNotExist:
            pass
        else:
            raise BranchAlreadyExists(branch_name)

        current_branch = self.get_current_branch()
        ws = self.workspace_storage.get_current_workspace()

        branch = BranchModel(
            name=branch_name,
            current_message_id=current_branch.current_message_id,
            workspace_id=ws.name,
        )
        self.session.add(branch)
        self.session.commit()

    def get_all_branches(self) -> list[BranchModel]:
        current_workspace = self.workspace_storage.get_current_workspace()
        return (
            self.session.query(BranchModel)
            .filter_by(workspace_id=current_workspace.name)
            .all()
        )

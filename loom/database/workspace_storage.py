from loom.database.db import SessionLocal, StateModel, WorkspaceModel
from loom.errors import (
    NoCurrentWorkspace,
    WorkspaceAlreadyExists,
    WorkspaceDoesNotExist,
)


class WorkspaceStorage:
    session = SessionLocal()

    def __init__(self):
        try:
            _ = self.get_current_workspace()
            return
        except NoCurrentWorkspace:
            pass

        try:
            default_workspace = self.get_workspace(name="default")
        except WorkspaceDoesNotExist:
            default_workspace = WorkspaceModel(name="default")
            self.session.add(default_workspace)
            self.session.commit()

        self.session.add(StateModel(key="WORKSPACE", value=default_workspace.name))
        self.session.commit()

    def get_workspace(self, name: str) -> WorkspaceModel:
        workspace = self.session.query(WorkspaceModel).filter_by(name=name).first()
        if workspace is None:
            raise WorkspaceDoesNotExist(name)
        return workspace

    def create_workspace(self, workspace_name: str):
        try:
            self.get_workspace(workspace_name)
        except WorkspaceDoesNotExist:
            pass
        else:
            raise WorkspaceAlreadyExists(workspace_name)

        new_workspace = WorkspaceModel(name=workspace_name)
        self.session.add(new_workspace)
        self.session.commit()

    def switch_to_workspace(self, workspace_name: str):
        workspace = self.get_workspace(workspace_name)
        current_workspace = (
            self.session.query(StateModel).filter_by(key="WORKSPACE").first()
        )
        if current_workspace is None:
            raise NoCurrentWorkspace()
        current_workspace.value = workspace.name
        self.session.commit()

    def get_current_workspace(self) -> WorkspaceModel:
        current_workspace_name = (
            self.session.query(StateModel).filter_by(key="WORKSPACE").first()
        )
        if current_workspace_name is None:
            raise NoCurrentWorkspace()
        return self.get_workspace(str(current_workspace_name.value))

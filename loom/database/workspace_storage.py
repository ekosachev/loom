from sqlalchemy.orm import Session
from loom.database.db import StateModel, WorkspaceModel
from loom.errors import (
    NoCurrentWorkspace,
    WorkspaceAlreadyExists,
    WorkspaceDoesNotExist,
)


class WorkspaceStorage:
    def __init__(self, session: Session):
        try:
            _ = self.get_current_workspace(session)
            return
        except NoCurrentWorkspace:
            pass

        try:
            default_workspace = self.get_workspace(session, name="default")
        except WorkspaceDoesNotExist:
            default_workspace = WorkspaceModel(name="default")
            session.add(default_workspace)
            session.commit()

        session.add(StateModel(key="WORKSPACE", value=default_workspace.name))
        session.commit()

    def get_workspace(self, session: Session, name: str) -> WorkspaceModel:
        workspace = session.query(WorkspaceModel).filter_by(name=name).first()
        if workspace is None:
            raise WorkspaceDoesNotExist(name)
        return workspace

    def create_workspace(self, session: Session, workspace_name: str):
        try:
            self.get_workspace(session, workspace_name)
        except WorkspaceDoesNotExist:
            pass
        else:
            raise WorkspaceAlreadyExists(workspace_name)

        new_workspace = WorkspaceModel(name=workspace_name)
        session.add(new_workspace)
        session.commit()

    def switch_to_workspace(self, session: Session, workspace_name: str):
        workspace = self.get_workspace(session, workspace_name)
        current_workspace = session.query(StateModel).filter_by(key="WORKSPACE").first()
        if current_workspace is None:
            raise NoCurrentWorkspace()
        current_workspace.value = workspace.name
        session.commit()

    def get_current_workspace(self, session: Session) -> WorkspaceModel:
        current_workspace_name = (
            session.query(StateModel).filter_by(key="WORKSPACE").first()
        )
        if current_workspace_name is None:
            raise NoCurrentWorkspace()
        return self.get_workspace(session, current_workspace_name.value)

    def get_all_workspaces(self, session: Session) -> list[WorkspaceModel]:
        return session.query(WorkspaceModel).all()

from datetime import UTC, datetime
from sqlalchemy import (
    DateTime,
    ForeignKey,
    Integer,
    Column,
    String,
    Text,
    create_engine,
)
from sqlalchemy.orm import declarative_base, sessionmaker
from pathlib import Path

Base = declarative_base()

LOOM_HOME = Path.home() / ".loom"
LOOM_HOME.mkdir(exist_ok=True)

DB_PATH = LOOM_HOME / "loom.db"


class MessageModel(Base):
    __tablename__ = "messages"

    id = Column(Integer, primary_key=True)
    session_id = Column(String, default="default")
    role = Column(String, nullable=False)
    name = Column(String, nullable=True)
    content = Column(Text, nullable=False)
    timestamp = Column(DateTime, default=datetime.now(tz=UTC))
    parent_id = Column(Integer, ForeignKey("messages.id"))


class BranchModel(Base):
    __tablename__ = "branches"

    id = Column(Integer, primary_key=True)
    workspace_id = Column(String, ForeignKey("workspaces.name"), nullable=False)
    name = Column(String, nullable=False)
    current_message_id = Column(Integer, ForeignKey("messages.id"), nullable=True)


class StateModel(Base):
    __tablename__ = "state"

    key = Column(String, primary_key=True)
    value = Column(String, nullable=False)


class WorkspaceModel(Base):
    __tablename__ = "workspaces"
    name = Column(String, primary_key=True)


engine = create_engine(f"sqlite:///{DB_PATH}")

SessionLocal = sessionmaker(bind=engine)

Base.metadata.create_all(engine)

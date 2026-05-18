from datetime import UTC, datetime
from typing import Optional
from sqlalchemy import (
    DateTime,
    ForeignKey,
    Integer,
    String,
    Text,
    create_engine,
    event,
)
from sqlalchemy.orm import Mapped, declarative_base, mapped_column, sessionmaker
from pathlib import Path

Base = declarative_base()

LOOM_HOME = Path.home() / ".loom"
LOOM_HOME.mkdir(exist_ok=True)

DB_PATH = LOOM_HOME / "loom.db"


class MessageModel(Base):
    __tablename__ = "messages"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    role: Mapped[str] = mapped_column(String, nullable=False)
    name: Mapped[Optional[str]] = mapped_column(String, nullable=True)
    content: Mapped[str] = mapped_column(Text, nullable=False)
    reasoning: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    timestamp: Mapped[datetime] = mapped_column(DateTime, default=datetime.now(tz=UTC))
    parent_id: Mapped[Optional[int]] = mapped_column(Integer, ForeignKey("messages.id"))


class BranchModel(Base):
    __tablename__ = "branches"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    workspace_id: Mapped[str] = mapped_column(
        String, ForeignKey("workspaces.name"), nullable=False
    )
    name: Mapped[str] = mapped_column(String, nullable=False)
    current_message_id: Mapped[Optional[int]] = mapped_column(
        Integer, ForeignKey("messages.id"), nullable=True
    )


class StateModel(Base):
    __tablename__ = "state"

    key: Mapped[str] = mapped_column(String, primary_key=True)
    value: Mapped[str] = mapped_column(String, nullable=False)


class WorkspaceModel(Base):
    __tablename__ = "workspaces"
    name: Mapped[str] = mapped_column(String, primary_key=True)


engine = create_engine(f"sqlite:///{DB_PATH}", connect_args={"timeout": 30})


@event.listens_for(engine, "connect")
def set_sqlite_pragma(dbapi_connection, connection_record):
    cursor = dbapi_connection.cursor()
    cursor.execute("PRAGMA journal_mode=WAL")
    cursor.execute("PRAGMA synchronous=NORMAL")
    cursor.close()


SessionLocal = sessionmaker(bind=engine)

Base.metadata.create_all(engine)

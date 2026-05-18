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
from sqlalchemy.orm import declarative_base, relationship, sessionmaker
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

    name = Column(String, primary_key=True)
    current_message_id = Column(Integer, ForeignKey("messages.id"), nullable=True)


class StateModel(Base):
    __tablename__ = "state"

    key = Column(String, primary_key=True)
    value = Column(String, nullable=False)


engine = create_engine(f"sqlite:///{DB_PATH}")

SessionLocal = sessionmaker(bind=engine)

Base.metadata.create_all(engine)

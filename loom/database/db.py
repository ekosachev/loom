from datetime import UTC, datetime
from httpx._transports import default
from sqlalchemy import DateTime, Integer, Column, String, Text, create_engine
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


engine = create_engine(f"sqlite:///{DB_PATH}")

SessionLocal = sessionmaker(bind=engine)

Base.metadata.create_all(engine)

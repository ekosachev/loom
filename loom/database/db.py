from datetime import UTC, datetime
from httpx._transports import default
from sqlalchemy import DateTime, Integer, Column, String, Text, create_engine
from sqlalchemy.orm import declarative_base, sessionmaker

Base = declarative_base()


class MessageModel(Base):
    __tablename__ = "messages"

    id = Column(Integer, primary_key=True)
    session_id = Column(String, default="default")
    role = Column(String, nullable=False)
    name = Column(String, nullable=True)
    content = Column(Text, nullable=False)
    timestamp = Column(DateTime, default=datetime.now(tz=UTC))


engine = create_engine("sqlite:///loom.db")

SessionLocal = sessionmaker(bind=engine)

Base.metadata.create_all(engine)

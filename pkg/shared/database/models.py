# X404X — Shared Database Models
# SQLAlchemy ORM models shared between Python modules

from datetime import datetime
from enum import Enum
from typing import Optional

from sqlalchemy import (
    Column, String, Integer, Float, Boolean, DateTime, Text,
    ForeignKey, JSON, Enum as SAEnum, create_engine
)
from sqlalchemy.orm import DeclarativeBase, relationship, Session


class Base(DeclarativeBase):
    pass


class AgentStatus(str, Enum):
    ONLINE = "online"
    IDLE = "idle"
    ACTIVE = "active"
    DEAD = "dead"


class CampaignStatus(str, Enum):
    DRAFT = "draft"
    RUNNING = "running"
    PAUSED = "paused"
    COMPLETED = "completed"
    FAILED = "failed"


class KillChainPhase(str, Enum):
    RECON = "recon"
    WEAPONIZATION = "weaponization"
    DELIVERY = "delivery"
    EXPLOITATION = "exploitation"
    INSTALLATION = "installation"
    C2 = "c2"
    ACTIONS_ON_OBJECTIVE = "actions_on_objective"
    EXFILTRATION = "exfiltration"


class RiskLevel(str, Enum):
    SAFE = "safe"
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    DANGER = "danger"


class Campaign(Base):
    __tablename__ = "campaigns"

    id = Column(String, primary_key=True)
    name = Column(String, nullable=False)
    target_scope = Column(String, nullable=False)
    goal = Column(String, nullable=False)
    profile = Column(String, default="balanced")
    status = Column(SAEnum(CampaignStatus), default=CampaignStatus.DRAFT)
    phase = Column(SAEnum(KillChainPhase))
    created_at = Column(DateTime, default=datetime.utcnow)
    started_at = Column(DateTime, nullable=True)
    completed_at = Column(DateTime, nullable=True)
    auto_approval = Column(Boolean, default=False)

    agents = relationship("Agent", back_populates="campaign", cascade="all, delete-orphan")
    missions = relationship("Mission", back_populates="campaign", cascade="all, delete-orphan")
    targets = relationship("Target", back_populates="campaign", cascade="all, delete-orphan")
    kill_chain = relationship("KillChainEntry", back_populates="campaign", cascade="all, delete-orphan")
    audit_log = relationship("AuditEntry", back_populates="campaign", cascade="all, delete-orphan")


class Agent(Base):
    __tablename__ = "agents"

    id = Column(String, primary_key=True)
    campaign_id = Column(String, ForeignKey("campaigns.id"), nullable=True)
    session_id = Column(String, nullable=True)
    hostname = Column(String, nullable=False)
    os = Column(String, nullable=False)
    arch = Column(String)
    username = Column(String)
    local_ip = Column(String)
    privileges = Column(JSON, default=list)
    status = Column(SAEnum(AgentStatus), default=AgentStatus.ONLINE)
    last_checkin = Column(DateTime, default=datetime.utcnow)
    first_seen = Column(DateTime, default=datetime.utcnow)
    uptime = Column(Integer, default=0)
    metadata = Column(JSON, default=dict)
    public_key = Column(Text, nullable=True)

    campaign = relationship("Campaign", back_populates="agents")
    credentials = relationship("Credential", back_populates="agent", cascade="all, delete-orphan")
    decisions = relationship("Decision", back_populates="agent")


class Target(Base):
    __tablename__ = "targets"

    id = Column(Integer, primary_key=True, autoincrement=True)
    campaign_id = Column(String, ForeignKey("campaigns.id"))
    ip = Column(String, nullable=False)
    hostname = Column(String)
    os = Column(String)
    open_ports = Column(JSON, default=list)
    services = Column(JSON, default=list)
    asset_value = Column(Integer, default=10)

    campaign = relationship("Campaign", back_populates="targets")
    vulnerabilities = relationship("Vulnerability", back_populates="target", cascade="all, delete-orphan")


class Vulnerability(Base):
    __tablename__ = "vulnerabilities"

    id = Column(Integer, primary_key=True, autoincrement=True)
    target_id = Column(Integer, ForeignKey("targets.id"))
    cve = Column(String)
    description = Column(Text)
    severity = Column(String)
    service = Column(String)
    port = Column(Integer)
    discovered_by = Column(String)
    discovered_at = Column(DateTime, default=datetime.utcnow)

    target = relationship("Target", back_populates="vulnerabilities")


class Mission(Base):
    __tablename__ = "missions"

    id = Column(String, primary_key=True)
    campaign_id = Column(String, ForeignKey("campaigns.id"))
    name = Column(String, nullable=False)
    description = Column(Text)
    tactic = Column(String)
    technique = Column(String)
    mitre_id = Column(String)
    status = Column(String, default="pending")
    order = Column(Integer, default=0)
    dependencies = Column(JSON, default=list)
    created_at = Column(DateTime, default=datetime.utcnow)
    completed_at = Column(DateTime, nullable=True)

    campaign = relationship("Campaign", back_populates="missions")


class KillChainEntry(Base):
    __tablename__ = "kill_chain"

    id = Column(Integer, primary_key=True, autoincrement=True)
    campaign_id = Column(String, ForeignKey("campaigns.id"))
    agent_id = Column(String)
    phase = Column(SAEnum(KillChainPhase))
    tactic = Column(String)
    technique = Column(String)
    mitre_id = Column(String)
    success = Column(Boolean, default=False)
    detail = Column(Text)
    timestamp = Column(DateTime, default=datetime.utcnow)

    campaign = relationship("Campaign", back_populates="kill_chain")


class Credential(Base):
    __tablename__ = "credentials"

    id = Column(Integer, primary_key=True, autoincrement=True)
    agent_id = Column(String, ForeignKey("agents.id"))
    username = Column(String, nullable=False)
    password = Column(String)
    hash = Column(Text)
    hash_type = Column(String)
    domain = Column(String)
    source = Column(String)
    captured_at = Column(DateTime, default=datetime.utcnow)

    agent = relationship("Agent", back_populates="credentials")


class Decision(Base):
    __tablename__ = "decisions"

    id = Column(String, primary_key=True)
    agent_id = Column(String, ForeignKey("agents.id"))
    campaign_id = Column(String)
    tactic = Column(String)
    technique = Column(String)
    mitre_id = Column(String)
    target = Column(String)
    confidence = Column(Float, default=0.0)
    reasoning = Column(Text)
    requires_approval = Column(Boolean, default=False)
    approved = Column(Boolean, nullable=True)
    source = Column(String, default="ai")
    timestamp = Column(DateTime, default=datetime.utcnow)

    agent = relationship("Agent", back_populates="decisions")


class AuditEntry(Base):
    __tablename__ = "audit_log"

    id = Column(Integer, primary_key=True, autoincrement=True)
    campaign_id = Column(String, ForeignKey("campaigns.id"))
    agent_id = Column(String)
    action = Column(String, nullable=False)
    result = Column(String)
    detail = Column(Text)
    timestamp = Column(DateTime, default=datetime.utcnow)

    campaign = relationship("Campaign", back_populates="audit_log")


class AIAnalysis(Base):
    __tablename__ = "ai_analysis"

    id = Column(Integer, primary_key=True, autoincrement=True)
    campaign_id = Column(String)
    agent_id = Column(String)
    model = Column(String)
    prompt = Column(Text)
    response = Column(Text)
    decision_taken = Column(Boolean)
    confidence = Column(Float)
    timestamp = Column(DateTime, default=datetime.utcnow)


class BlueMetric(Base):
    __tablename__ = "blue_metrics"

    id = Column(Integer, primary_key=True, autoincrement=True)
    campaign_id = Column(String)
    tool = Column(String)
    detected = Column(Boolean, default=False)
    alert_type = Column(String)
    agent_id = Column(String)
    timestamp = Column(DateTime, default=datetime.utcnow)


def create_session(database_url: str) -> Session:
    engine = create_engine(database_url)
    Base.metadata.create_all(engine)
    return Session(engine)

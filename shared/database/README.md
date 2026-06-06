# RBYHACK Framework — Database Setup
# Run: alembic init -t async shared/database/migrations
# Then: alembic upgrade head

# This file documents the migration commands.
# Database: PostgreSQL (production) / SQLite (development/agent)

"""
Supported databases:
  - SQLite (default, zero-config)
  - PostgreSQL (production)

Quickstart:
  1. pip install alembic sqlalchemy
  2. alembic init -t async shared/database/migrations
  3. Edit migrations/env.py to import models.Base metadata
  4. alembic revision --autogenerate -m "init"
  5. alembic upgrade head
"""

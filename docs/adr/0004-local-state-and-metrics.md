# ADR 0004: local SQLite state and observations

Status: accepted

SQLite stores local operational state and bounded aggregate observations.
Human-readable JSON operation records accompany it. Secrets and full SQL are
excluded. No telemetry leaves the host.

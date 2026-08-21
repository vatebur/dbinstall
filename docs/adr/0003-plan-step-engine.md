# ADR 0003: plan and step engine

Status: accepted

Planning and applying use the same stable step graph. Steps inspect, describe,
apply, verify, and compensate bounded state changes. The first executor is
sequential and journaled; parallelism is reserved for independent read-only or
content-addressed work.

# Roadmap

## Milestone 0: foundation (this implementation)

- Product design, threat model, support policy, and ADRs.
- Versioned YAML specification and validation.
- Environment/file secret references with redaction boundaries.
- Host and OS discovery for the certified Linux families.
- Compiled provider registry and descriptors for MySQL, GreatSQL, PostgreSQL.
- Deterministic step plan, JSON/human output, local state, audit IDs, and locks.
- `plan`, `status`, and `version` commands.
- Unit and provider-contract tests plus CI.
- MySQL plan generation without enabling mutating installation steps.

## Milestone 1: trusted distributions

- Signed artifact manifest format and trust store.
- Online cache with resume and bounded retry.
- Offline artifact verification.
- MySQL 8.4 archive distribution installation.
- Typed filesystem, account, systemd, and process operations.
- Failure compensation and integration tests.

## Milestone 2: real single instances

- MySQL initialization, configuration publication, startup, and verification.
- GreatSQL stable-series artifact and instance implementation.
- PostgreSQL thin vertical slice.
- Package-manager backends and their single-instance constraints.
- VM certification for the initial OS matrix.

## Milestone 3: configuration advisor

- Host resource allocator for multiple instances.
- Versioned static rule sets and recommendation explanations.
- Bounded `observe` task with privacy-preserving SQLite samples.
- Runtime recommendation generation and explicit `tune apply` workflow.
- Atomic configuration publication and tested automatic restoration.

## Later milestones

- Explicit adoption workflow.
- Upgrade, downgrade constraints, and complete uninstall.
- MySQL-family replication and MGR topology layer.
- PostgreSQL HA topology.
- External RPC Provider protocol and stable Go SDK.
- Prometheus metric import and optional central orchestration.

Containers, Kubernetes Operators, storage provisioning, and business-data
backup remain separate products/backends and are not implicit extensions of the
host installer.

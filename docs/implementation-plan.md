# Implementation plan and acceptance gates

This document turns the product roadmap into independently reviewable delivery
stages. A stage is complete only when its code, tests, documentation, and
failure behavior are committed together.

## Package boundaries

| Package | Responsibility | Must not do |
|---|---|---|
| `internal/spec` | Strict schema, defaults, validation, normalization | Resolve secrets or inspect hosts |
| `internal/secrets` | Resolve explicit local secret references | Log or persist values |
| `internal/host` | Read-only host discovery | Apply remediation |
| `internal/provider` | Provider contracts and registry | Switch on product names in core |
| `internal/plan` | Stable JSON plan types | Perform I/O |
| `internal/planner` | Combine desired state, host state, and Provider steps | Mutate the host |
| `internal/executor` | Sequential checks, apply, verify, compensation | Invent Provider actions |
| `internal/state` | SQLite instance and operation catalog | Store resolved secrets |
| `internal/audit` | Atomic human-readable operation artifacts | Replace authoritative state |
| `internal/lock` | Host mutation exclusion | Coordinate remote hosts |
| `internal/providers/*` | Product/version behavior and capabilities | Modify unowned global resources |

## Delivery stages

### Stage A: contracts and read-only planning

Deliver strict YAML, host inspection, three built-in Provider descriptors,
stable plans, state primitives, audit files, locks, CLI output, and unit tests.
No installation step is executable. This is the foundation milestone.

Acceptance:

- Unknown YAML fields and unpinned versions fail.
- JSON output contains no resolved credentials.
- MySQL plans show every intended lifecycle step and `executable: false`.
- Binding `0.0.0.0` produces a high-risk exposure warning.
- Provider contracts reject duplicate or forward step dependencies.
- Race tests, vet, build, examples, and CLI smoke tests pass on Go 1.26.5.

### Stage B: artifact trust and typed host operations

Implement signed manifest parsing, trust roots, content-addressed cache, offline
verification, typed filesystem/user/systemd/process operations, and a fake-host
contract suite. Enable only artifact resolution after its fault-injection tests
pass.

Acceptance includes corrupted downloads, interrupted resume, wrong
architecture, signature failure, cache races, filesystem-full behavior, and no
shell interpolation.

### Stage C: MySQL 8.4 archive vertical slice

Implement one exact MySQL 8.4 patch release on one certified OS/architecture
combination first. Add immutable distribution installation, empty-directory
initialization, configuration validation, local administrator creation,
dedicated systemd unit, TLS, startup verification, and compensation.

Expand to the full OS matrix only after the tracer path is green. An `apply`
command remains hidden until every emitted action is implemented and the plan
is executable.

### Stage D: GreatSQL and PostgreSQL proof

Select official maintained stable series and manifests. Reuse MySQL-family
components for GreatSQL through composition. Implement a narrow PostgreSQL path
to force validation of port, configuration, account, initialization, metrics,
and service abstractions.

### Stage E: static configuration advisor

Implement host resource reservation, multi-instance budget accounting,
workload profiles, typed parameter metadata, versioned formulas, provenance,
risk, and configuration validation. Golden tests cover resource boundaries and
version-specific rules.

### Stage F: bounded runtime observation

Implement the local observation job, privacy filter, sample retention and size
limits, MySQL aggregate collectors, rule evidence, recommendation confidence,
and explicit tune application. Fault tests cover missing schemas, restarts,
clock changes, full disks, and partial samples.

### Stage G: certification and release

Run real systemd VM tests for every advertised OS/product/architecture tuple.
Generate signed binaries, checksums, SBOM, third-party notices, and upgrade notes
for schema, state, rules, and artifact manifests.

## Commit strategy

Each stage uses small coherent commits in this order when applicable:

1. Contract or ADR and failing tests.
2. Pure domain implementation.
3. Host integration implementation.
4. Fault-injection and VM coverage.
5. Documentation and release metadata.

Never combine new destructive behavior with an unrelated refactor. A commit
that makes a new Step executable must include its verification and compensation
tests.

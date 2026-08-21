# Threat model

## Protected assets

- Existing database data, services, system configuration, and package policy.
- Administrative credentials and runtime query privacy.
- Integrity of downloaded database binaries and generated configuration.
- Accurate audit history and ownership records.

## Principal threats and controls

| Threat | Primary controls |
|---|---|
| Artifact substitution | Pinned version, SHA-256, vendor signature, explicit trust |
| Command injection | Typed operations, argv execution, no shell by default |
| Accidental host takeover | Plan-first, owned resources, conflict blocking |
| Excess privilege | Caller-controlled elevation, no internal sudo, DB user isolation |
| Secret disclosure | Reference resolution, redaction, no secrets in state or plans |
| Database exposure | Exposure plan, local admin, TLS, explicit remote accounts |
| Configuration outage | Native validation, atomic publish, health check, restore |
| Concurrent mutation | Host and instance locks, operation journal |
| Sensitive query collection | Aggregate metrics, no SQL text or identities by default |
| Supply-chain drift | Versioned manifests, cached content address, SBOM |

dbinstall cannot prove external security-group policy, recover business data,
or make unsafe application parameters safe. Those limitations must remain
visible in plans and documentation.

# dbinstall product and architecture design

## 1. Scope

`dbinstall` is a host-local CLI and reusable internal engine for native Linux
database installation. Its first product providers are GreatSQL, Oracle MySQL,
and PostgreSQL. GreatSQL and MySQL receive complete support; PostgreSQL provides
a thin but real vertical slice to prove that core abstractions are not tied to
the MySQL family.

The first release supports online official artifacts and offline local
artifacts. It does not support containers, Kubernetes, source compilation,
cross-host clusters, upgrades, downgrades, full uninstall, storage provisioning,
or business-data backup and restore.

## 2. Control and public contracts

The executable runs on the target host. There is no central controller and no
privileged resident agent. Core packages remain reusable so a future SSH or RPC
controller can call the same engine.

The stable first-release contracts are:

1. CLI behavior and documented exit codes.
2. Versioned YAML installation specifications.
3. JSON output from `plan`, `status`, and `verify`.

Internal Go interfaces may evolve until all three initial providers have proven
them. External plugins and a stable Go SDK are deferred.

## 3. Layered architecture

```text
CLI / JSON presentation
        |
Spec loading, validation, defaults, secret references
        |
Planner ---- host discovery ---- state and drift
        |
Step graph and execution journal
        |
Provider registry and capability components
        |
Typed host operations (files, packages, systemd, processes)
```

The core never switches on a database name. A compiled-in provider registers a
descriptor and creates steps. `mysqlfamily` capability components share MySQL
and GreatSQL semantics without pretending the products are identical.

## 4. Desired-state and execution model

YAML is the source of truth and has an `apiVersion`. Interactive UX may only
generate YAML; the engine never prompts. Each provider emits stable steps with:

- `Check`: inspect current state without mutation.
- `Plan`: describe the exact change and risk.
- `Apply`: perform one bounded mutation.
- `Verify`: prove the postcondition.
- `Rollback`: compensate for an uncommitted mutation where practical.

Steps have stable IDs, dependencies, timeouts, retry classification, required
privilege, and owned resources. First-release execution is sequential. Failed
operations roll back only resources created or modified by that operation; they
never promise to reverse application data writes.

## 5. Distribution and instance model

A `Distribution` is an immutable database/version/OS/architecture installation.
An `Instance` references one distribution and owns configuration, port, socket,
data, logs, temporary files, and a systemd unit.

Archive distributions can serve several instances. Vendor RPM/DEB packages
usually own fixed paths and services, so the first release limits a
package-managed product to one managed instance per host. True multi-instance
deployment uses the archive backend.

Only a nonexistent or verified-empty data directory may be used. Existing
installations are inspected and reported, never automatically adopted. Any
port, path, user, package, or systemd conflict blocks the plan.

## 6. Host support and mutation policy

Certified first-release targets are Rocky Linux and AlmaLinux 8/9, Debian 12,
and Ubuntu 22.04/24.04, on amd64 and arm64 when the vendor publishes matching
artifacts. Similar systems can be detected but remain unsupported until their
VM end-to-end matrix passes.

The engine does not call `sudo`, replace global package mirrors, disable a
firewall, or disable SELinux. Applying system tuning requires explicit
`--apply-system-tuning` authorization. Changes retain their old values and are
recorded for restoration. The first release supports systemd only.

Providers use typed host operations. Arbitrary shell execution is a restricted
escape hatch for signed, bundled vendor helpers with declared inputs, timeout,
privilege, and side effects. Commands otherwise execute directly with argument
arrays, not through a shell.

## 7. Artifact trust

Production specifications pin a complete version. A provider artifact manifest
maps product, version, OS, architecture, and installation method to URL, size,
SHA-256, and vendor-signature metadata. `latest` is only an interactive lookup;
planning resolves it to a locked version.

Downloads use a content-addressed cache, timeouts, bounded retries, and resume.
Official URLs may be replaced by a configured corporate mirror, but identity,
hash, and signature requirements remain unchanged. Offline artifacts undergo
the same checks. Custom manifests must be explicitly trusted.

## 8. Configuration and tuning

Initial configuration combines:

1. Vendor defaults.
2. A workload profile (`development`, `oltp`, `olap`, or `custom`).
3. A host-level resource budget divided among managed instances.
4. Deterministic provider/version tuning rules.
5. Typed user overrides.
6. Explicit unverified extra configuration.

Every value retains provenance. Instance-owned parameters such as data path,
port, socket, PID file, and service user cannot be overridden through free-form
configuration.

Before publication, configuration is generated in a temporary location,
validated with the database, atomically installed, and health checked. Failure
restores the prior configuration and retries startup once before preserving the
failure site for diagnosis.

Runtime advice comes from a bounded local observation job, not a monitoring
platform. It stores host metrics and aggregate database counters in SQLite.
Full SQL, parameters, table names, users, and client addresses are not collected
by default. Advice includes current value, proposed value, formula or rule,
evidence, confidence, risk, dynamic applicability, and restart requirement.
Runtime advice is never applied without an explicit `tune apply` command.

## 9. Configuration ownership and drift

Package installations use a clearly named include fragment and never overwrite
the global MySQL or PostgreSQL configuration when includes are available.
Archive instances receive dedicated configuration directories. Generated files
include the tool version, specification digest, and content hash.

Manual changes to owned files produce drift. Planning reports drift and refuses
to overwrite it without explicit approval. Files not owned by dbinstall remain
outside its rollback and cleanup scope.

## 10. State, audit, locks, and secrets

Operational state lives under `/var/lib/dbinstall` by default:

- `state.db`: SQLite catalog of instances, artifacts, owned resources, and
  operations.
- `operations/`: human-readable plan, specification snapshot, and JSON event
  stream for every operation.
- `backups/`: bounded retention of modified configuration and service files.

A host file lock prevents concurrent mutation. State contains no secrets.
Specifications reference secrets with `env://NAME` or `file:///path`; generated
credentials are stored in a mode-0600 file and reported once. Logs and errors
redact secret values.

## 11. Network and accounts

Per product decision, a new database defaults to `0.0.0.0`. The plan therefore
lists every interface and exposed database port and emits a high-risk warning
when host access controls cannot be confirmed. dbinstall never silently changes
the firewall.

Administrative accounts remain local-only. Anonymous users, empty passwords,
and test databases are forbidden. Runtime collection uses a separate
least-privilege account. Remote application accounts require an explicit source
CIDR. TLS is enabled by default when supported, and remote plaintext
administrative login is forbidden.

## 12. Privacy, logging, and output

dbinstall uploads no telemetry. Network access is limited to explicit artifact
or catalog operations. Metrics and diagnostic material remain local.

Terminal output is human-oriented. `--output json` is a clean stable document
without progress messages. Every operation has an ID shared by CLI output,
SQLite state, and the JSON event log. Exit codes distinguish validation,
preflight, drift, execution, verification, and partial-rollback failures.

## 13. Testing and support claims

Support requires four layers: unit tests, provider contract tests, real database
integration tests, and systemd VM end-to-end tests. Containers cannot certify a
host OS. Only combinations passing the entire matrix are advertised as
supported; all others are experimental.

Each product starts with one stable/LTS major series and the current and prior
patch releases. MySQL starts with 8.4 LTS. GreatSQL and PostgreSQL select their
still-maintained stable series when artifact manifests are implemented.

## 14. Licensing and release integrity

The project uses Apache-2.0. GPL implementations can inform externally visible
behavior but are not copied. Dependencies require license review, an SBOM, and
third-party notices.

Releases contain signed amd64 and arm64 binaries, SHA-256 files, and an SBOM.
The tool version, schema version, provider version, tuning-rule version, and
artifact-manifest version remain independently traceable.

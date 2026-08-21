# dbinstall

`dbinstall` is a database-neutral, host-local installer and configuration
advisor written in Go.

The first supported providers are GreatSQL, Oracle MySQL, and PostgreSQL. The
project is deliberately conservative: it produces an auditable plan before it
changes a host, owns only resources it creates, and never uploads host or
database telemetry.

> Status: foundation milestone. `plan` and `status` are implemented before any
> production database installation step is enabled.

## Product principles

- A versioned YAML document is the source of truth.
- `plan` and `apply` share the same step graph.
- Host changes are explicit, minimal, reversible where practical, and never
  invoke `sudo` internally.
- Exact database artifacts are pinned and integrity checked.
- Static tuning is deterministic. Runtime tuning is explainable and opt-in.
- Existing installations are discovered but never adopted implicitly.
- Human output and stable JSON output are separate contracts.

## Initial commands

```text
dbinstall plan --file instance.yaml
dbinstall status [--instance INSTANCE_ID]
dbinstall inspect
dbinstall providers
dbinstall validate --file instance.yaml
dbinstall version
```

See [docs/design.md](docs/design.md), [docs/implementation-plan.md](docs/implementation-plan.md),
[docs/roadmap.md](docs/roadmap.md), and [docs/specification.md](docs/specification.md)
for the complete agreed design.

## Development toolchain

- Go 1.26.5
- Rust 1.90.0 (reserved for tooling; the product remains Go-first)
- Node.js 26.7.0 (documentation and repository tooling)

Tool versions are managed with `vfox`.

```text
vfox use -p rust@1.90.0
vfox use -p nodejs@26.7.0
vfox use -p golang@1.26.5
make check
```

## License

Apache License 2.0. GPL projects may inform behavior and design, but their
implementation code is not copied into this repository.

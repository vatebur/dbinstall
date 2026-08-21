# ADR 0006: default listen address is all IPv4 interfaces

Status: accepted by product owner

New instances default to `0.0.0.0`. This intentionally differs from the
security-conservative recommendation of localhost-only binding.

The decision requires compensating controls: the plan always emits a high-risk
warning and lists exposed interfaces and ports; administrative accounts remain
local-only; anonymous and empty-password accounts are forbidden; remote
application accounts require explicit source policy; TLS is enabled when the
database supports it; and dbinstall never silently opens a firewall.

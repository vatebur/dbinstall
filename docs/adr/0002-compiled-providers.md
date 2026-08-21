# ADR 0002: compiled providers

Status: accepted

Initial Providers compile into the binary and register capabilities. Go shared
object plugins are rejected because toolchain and dependency identity make them
fragile across distributions. A future extension protocol will use process
isolation and RPC after internal interfaces stabilize.

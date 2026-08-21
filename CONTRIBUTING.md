# Contributing

Use the project-pinned toolchain through vfox:

```text
vfox use -p rust@1.90.0
vfox use -p nodejs@26.7.0
vfox use -p golang@1.26.5
```

Code must pass `make check`, `test -z "$(gofmt -l .)"`, and
`git diff --check`.

Public behavior starts with a test. Providers must pass the common contract
suite. A platform is not called supported until its real-VM end-to-end test
passes. Shell commands are not accepted where a typed host operation can express
the same behavior.

Commits should be small, coherent milestones. Do not mix generated artifacts or
unrelated formatting with behavior changes.

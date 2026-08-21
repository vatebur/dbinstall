# Installation specification

The desired state uses `dbinstall.vatebur.dev/v1alpha1`. Unknown fields are
rejected. Defaults are materialized before a plan digest is calculated.

```yaml
apiVersion: dbinstall.vatebur.dev/v1alpha1
kind: DatabaseInstance
metadata:
  name: orders-mysql
spec:
  provider: mysql
  version: 8.4.6
  installation:
    method: archive
    source: online
  instance:
    port: 3306
    bindAddress: 0.0.0.0
    paths:
      data: /var/lib/dbinstall/instances/orders-mysql/data
      logs: /var/log/dbinstall/orders-mysql
      temp: /var/lib/dbinstall/instances/orders-mysql/tmp
  workload:
    profile: oltp
    memoryLimit: 8GiB
    expectedConnections: 300
    durability: strict
  credentials:
    adminPassword: file:///run/secrets/orders-mysql-admin
  configOverrides:
    max_connections: 300
  extraConfig: {}
```

`provider` is one of `mysql`, `greatsql`, or `postgresql` in the foundation
release. `method` is `archive` or `package`; `source` is `online` or `offline`.
Offline specifications add a local artifact path.

Instance names are DNS-label-like stable IDs. Paths must be absolute and must
not alias one another. Ports must be non-privileged and unique among discovered
and managed services. Memory quantities use IEC suffixes.

Secrets use `env://NAME` or an absolute `file:///path`. Values are resolved only
at execution boundaries and are never returned in JSON plans.

The initial default bind address is `0.0.0.0` by explicit product decision.
Plans attach a security warning and providers retain local-only administrative
accounts.

# ADR 0001: host-local CLI first

Status: accepted

dbinstall runs on the target host as a CLI. It has no central controller or
privileged resident daemon. Internal boundaries preserve a future controller
option without making networking part of the first trust boundary.

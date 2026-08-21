# ADR 0005: conservative host mutation

Status: accepted

dbinstall never disables SELinux or firewalls, replaces global mirrors, or calls
sudo. System tuning is explicit, recorded, and recoverable. Existing database
resources are discovered but not adopted, and conflicts block execution.

# Integration tests

This directory is reserved for tests that exercise interactions between real project
components (for example Registry nodes, the Load Balancer, and mock services).

No domain integration test is included in the bootstrap because the corresponding
protocols and application behavior have not been implemented yet.

Integration tests added here should be runnable through `make test-integration`.

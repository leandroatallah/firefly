# Contracts

This package defines **interfaces (contracts)** used throughout the engine.

## Why No Tests?

This package contains **only interface definitions** — no implementation logic. Interfaces are validated at compile-time when concrete types implement them. Testing happens at the implementation level (e.g., `../physics/body/`, `../physics/space/`).

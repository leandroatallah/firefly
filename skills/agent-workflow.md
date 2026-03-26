# Chained Agent Workflow

A multi-agent pipeline for automating test coverage improvement.

## Agents

**1. Coverage Analyzer**
Runs `go test -coverprofile` on all target packages, parses output, and produces a prioritized list of files/functions with low or zero coverage.

**2. Code Reader**
Takes the low-coverage file list, reads source files, and extracts function signatures, state machine branches, and interface dependencies. Outputs a structured gap report.

**3. Mock Generator**
Inspects `internal/engine/contracts/` and `internal/engine/mocks/`, then generates or updates mock implementations needed by the test writer. Decides shared vs package-local placement.

**4. Test Writer**
Consumes the gap report and available mocks to write `_test.go` files. Follows table-driven patterns, avoids GPU/window dependencies, and respects code style rules.

**5. Test Validator**
Runs the newly written tests, checks for failures or flakiness, and re-runs coverage to confirm improvement. Reports delta coverage per package.

**6. Orchestrator**
Coordinates the full pipeline: triggers agents in order, passes outputs between them, retries failed steps, and produces a final coverage summary.

## Chain

```
Analyzer → Code Reader → Mock Generator → Test Writer → Validator
                                                              ↑
                                                        Orchestrator
```

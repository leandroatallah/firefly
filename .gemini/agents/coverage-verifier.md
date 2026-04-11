---
name: coverage-verifier
description: Runs tests, detects flakiness, measures coverage improvement, and reports delta
kind: local
tools:
  - run_shell_command
  - read_file
---


# Coverage Verifier

## Purpose

Runs newly written tests, checks for failures/flakiness, re-runs coverage analysis, and reports delta coverage per package.

## Responsibilities

- Execute tests written by Test Writer
- Detect test failures and provide diagnostic information
- Run tests multiple times to detect flakiness
- Re-run coverage analysis on updated packages
- Calculate coverage delta (before vs. after)
- Validate that coverage improvement meets expectations
- Report results with actionable feedback

## Inputs

- Newly written test files from Test Writer
- Baseline coverage report from Coverage Analyzer

## Outputs

- Test execution report:
  - Pass/fail status
  - Flakiness detection results
  - Coverage delta per package
  - Overall progress toward 80%+ goal
- Failure diagnostics (if any)
- Recommendations for Test Writer if improvements needed

## Verification Steps

1. Run `go test ./path/to/package -v`
2. Check for failures or panics
3. Run tests 3-5 times to detect flakiness
4. Run `go test -coverprofile=coverage_new.out`
5. Compare with baseline coverage
6. Generate delta report

## Success Criteria

- All tests pass
- No flaky tests detected
- Coverage increased by expected amount
- No regression in existing coverage

## Integration

Receives tests from **Test Writer**, reports results to **Workflow Orchestrator**.

---
name: gap-detector
description: Analyzes source files to extract function signatures, branches, and interface dependencies for gap reporting
tools:
  - read_file
  - grep_search
  - glob
---


# Gap Detector

## Purpose

Reads source files from low-coverage list, extracts function signatures, state machine branches, and interface dependencies. Outputs structured gap report identifying what's untested.

## Responsibilities

- Read source files identified by Coverage Analyzer
- Parse Go code to extract:
  - Function signatures and parameters
  - Conditional branches (if/else, switch/case)
  - State machine transitions (especially in `entity/actors`)
  - Interface dependencies from `internal/engine/contracts/`
- Identify untested code paths and edge cases
- Determine required mocks based on interface usage
- Generate actionable gap report

## Inputs

- Coverage report from Coverage Analyzer
- Source file paths

## Outputs

- Structured gap report containing:
  - Function name and signature
  - Untested branches/paths
  - Required mock interfaces
  - Suggested test scenarios (happy path, edge cases, error cases)
  - Complexity assessment

## Analysis Focus

- State machine logic in `handleState` methods
- Physics calculations with fixed-point arithmetic
- Scene lifecycle methods
- i18n edge cases (missing keys, formatting)
- Collision detection edge cases

## Integration

Receives input from **Coverage Analyzer**, feeds report to **Mock Generator** and **Test Writer**.

# Extension Hooks

Place hook files here to inject custom steps into the pipeline.

## Supported hooks

| File | Triggered by | Purpose |
|---|---|---|
| `before_spec.md` | Story Architect, Spec Engineer | Run before any spec is written |

## How to use

Create a hook file with instructions in plain text. Agents will read and follow it before proceeding.

Example `.agents/hooks/before_spec.md`:
```
Check if a related spec already exists in .agents/work/active/ or .agents/work/done/
before creating a new one. If a duplicate is found, report it and stop.
```

Delete the file to disable the hook.

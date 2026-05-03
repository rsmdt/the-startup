---
title: "Activity Export — Solution"
status: completed
---

# Solution Design

## Components

### ActivityQuery (unit A)
A new helper at `src/queries/activity.ts` that returns an async iterator over user-activity rows. Uses the existing `db` connection. Must exist before the controller can use it.

### ExportController (unit B)
A new controller at `src/controllers/export.ts` exposing GET /export/activity. Imports ActivityQuery, pipes its iterator through a CSV serializer, and streams the response. **Depends on ActivityQuery being implemented first** — the controller imports its type and iterator interface.

## Key Decisions

- **ADR-1**: Stream rather than buffer — large activity tables would otherwise OOM the server.
- **ADR-2**: Sequential implementation order — ActivityQuery first, ExportController second. The dependency is structural (import + type), not just stylistic.

# Append-only audit log — REST only, action + ref, per-project scope

An `AuditEntry` table records every mutation to Issues and Epics. Each row stores: `project_id`, `entity_type` ("issue"
or "epic"), `ref` (e.g. `PROJ-42`, `PROJ-E1`), `title` (the entity's name at the time of the action), `action`
("created", "updated", "moved", "deleted"), an optional `detail` summary string, and `created_at`. No `updated_at`, no
`deleted_at` — rows are never modified or removed.

**Human-only:** The audit feed is intentionally excluded from MCP tools. Agents act on the current state; the audit
trail is for human review only.

**Action + ref + detail:** No full payload snapshot. A lightweight, optional `detail` string carries a human-readable
summary of the change — "Backlog → In Progress" for a move, the changed field names ("name, priority") for an update,
the entity name for create/delete. This keeps the log useful at a glance (ZDZ-17) without storing before/after copies; if
"what did it look like" is needed, `created_at` can still be correlated with git history or comments.

**Per-project scope:** Exposed at `GET /projects/{slug}/audit`. A global feed adds no value in a single-user tool and
complicates queries.

**Issues and Epics only:** Comments and Links have no domain ref (no `PROJ-42`-style identifier), so they are excluded
to keep the log human-readable without a lookup table.

**Hook point:** Service layer. Each mutating method on `IssueService` and `EpicService` calls `stores.Audit.Record`
after a successful `write`. This is the only layer that knows semantic action names (e.g. "moved" vs. "updated").

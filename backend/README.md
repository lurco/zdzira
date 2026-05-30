# Backend

The Go service behind Zdzira. It exposes a REST API for humans and an MCP server
for AI agents over a single SQLite database, and streams board updates via SSE.
See the [root README](../README.md) for setup and the full endpoint list.

## Layout

```
api/      HTTP layer — Chi router + Huma (OpenAPI) handlers, SSE, middleware
mcp/      MCP server (streamable HTTP) — agent-facing tools
model/    GORM models and enums (the source of truth for the schema)
service/  Business logic — all mutations and validation live here
store/    Persistence — GORM over SQLite, one store per aggregate
```

The dependency direction is one-way: `api` and `mcp` call `service`, `service`
calls `store`. Handlers stay thin; rules and reference parsing (e.g. `PROJ-42`)
belong in `service`. Both the REST and MCP layers share the same services, so a
mutation behaves identically whoever triggers it — and both fire the SSE
broadcast on success.

## Develop

```sh
go test ./...                       # all backend tests
go build -o bin/zdzira ./cmd/zdzira # build the binary
go run ./cmd/zdzira -addr :8080     # run against ./zdzira.db
make openapi                        # regenerate docs/openapi.json from the routes
```

The OpenAPI spec is generated from the Huma route definitions — there is no
hand-maintained spec. After adding or changing a route, run `make openapi`.

## Conventions

- New endpoint: register it in `api/`, put the logic in a `service` method,
  add persistence to the relevant `store`, and cover the service method with a
  test (`service/*_test.go` use an in-memory SQLite).
- Architectural decisions are recorded in [`../docs/adr`](../docs/adr).
- Domain vocabulary is fixed in [`../CONTEXT.md`](../CONTEXT.md); the data model
  is in [`../docs/erd.md`](../docs/erd.md).

## REST endpoints

Base path `/api/v1`. Issues are referenced as `PROJ-42`, epics as `PROJ-E1`.

| Method         | Path                                               | Description             |
|----------------|----------------------------------------------------|-------------------------|
| GET / POST     | `/projects`                                        | List / create projects  |
| GET / DELETE   | `/projects/{slug}`                                 | Get / delete project    |
| GET            | `/projects/{slug}/board`                           | Board view              |
| GET / POST     | `/projects/{slug}/epics`                           | List / create epics     |
| GET/PUT/DELETE | `/projects/{slug}/epics/{epicRef}`                 | Epic CRUD               |
| GET / POST     | `/projects/{slug}/issues`                          | List / create issues    |
| GET/PUT/DELETE | `/projects/{slug}/issues/{issueRef}`               | Issue CRUD              |
| POST           | `/projects/{slug}/issues/{issueRef}/move`          | Move to swimlane        |
| GET / POST     | `/projects/{slug}/issues/{issueRef}/comments`      | List / add comments     |
| PUT / DELETE   | `/projects/{slug}/issues/{issueRef}/comments/{id}` | Edit / delete comment   |
| GET / POST     | `/projects/{slug}/issues/{issueRef}/links`         | List / create links     |
| DELETE         | `/projects/{slug}/issues/{issueRef}/links/{id}`    | Delete link             |
| GET / POST     | `/projects/{slug}/swimlanes`                       | List / create swimlanes |
| PATCH / DELETE | `/projects/{slug}/swimlanes/{id}`                  | Update / delete swimlane|
| POST           | `/projects/{slug}/swimlanes/reorder`               | Reorder swimlanes       |
| GET            | `/projects/{slug}/audit`                           | Audit log               |
| GET            | `/api/v1/events`                                   | SSE board update stream |
| GET            | `/health`, `/ready`                                | Liveness / readiness    |

The OpenAPI spec ([`../docs/openapi.json`](../docs/openapi.json), live at `/docs`)
and the Postman collection in [`../docs`](../docs) are the authoritative reference.

## MCP tools

Agent-facing tools exposed at `/mcp`:

| Tool | Description |
|------|-------------|
| `list_projects` / `get_project` / `create_project` | Project read & create |
| `list_epics` / `get_epic` / `create_epic` / `update_epic` | Epic read, create, update |
| `list_issues` / `get_issue` / `create_issue` / `update_issue` | Issue read, create, update (filterable list) |
| `move_issue` / `delete_issue` | Move between swimlanes / delete |
| `get_board` | Full board (swimlanes + issues) |
| `add_comment` / `list_comments` | Comments on an issue or epic (one of `issue_ref` / `epic_ref`) |
| `link_issues` / `list_links` | Issue links |
| `list_swimlanes` | Swimlanes in a project |

### Read model (agent-facing)

MCP read tools return projections shaped for an agent's attention and context budget, not the raw models — see
[ADR 0008](../docs/adr/0008-mcp-agent-read-model.md):

- **Detail calls inline feedback.** `get_issue` returns the description plus the issue's `comments` and `links`;
  `get_epic` inlines its `comments`. An agent acting on an entity always sees the feedback attached to it.
- **List/board calls are a scannable index.** `list_issues`, `get_board`, and `list_epics` omit descriptions and carry a
  `comment_count` per row — the cue to drill in with `get_issue` / `get_epic`. Counts are batched (one grouped query),
  so there is no N+1 on the board.

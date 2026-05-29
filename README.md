# Zdzira

A local issue tracker for personal software development in the AI age. No accounts, no auth — access is direct.
Exposes a **REST API** for human access and an **MCP server** for AI agents (Claude Code, Claude Desktop, etc.).
Runs as a single Go binary backed by SQLite.

![Zdzira board](docs/screenshot1b.png)

---

## Features

- Kanban board with drag-and-drop lanes and issues
- Issues: types (Task / Bug / Story), priorities (Low / High / Immediate), epic grouping
- Epics manager with issue listings
- Comments on issues (visible to both humans and agents)
- Issue links (Blocks, Duplicates, Relates To, Is Part Of)
- Filtering by type, priority, epic, and free-text search
- Real-time board sync via Server-Sent Events — agent changes appear instantly
- Light, dark, and high-contrast themes
- MCP server for AI agent integration
- OpenAPI docs at `/docs`
- Postman collection in `docs/`

---

## Quick Start

### Binary (simplest)

```sh
go build -o bin/zdzira ./cmd/zdzira
./bin/zdzira
```

Opens on `:8080`. SQLite database is created at `./zdzira.db`.

Options:

```sh
./bin/zdzira -addr :9000 -db ~/my-projects.db
```

Then open `http://localhost:8080` in your browser.

---

### Docker Compose — Development

The dev stack runs the Go backend and a Vite dev server with hot-module reload. Nginx is disabled.

```sh
docker compose up
```

| Service  | URL                   |
|----------|-----------------------|
| Frontend | http://localhost:5173 |
| Backend  | http://localhost:8080 |

The frontend Vite server proxies `/api/v1` and `/mcp` to the backend container.

> **Backend changes** require a container rebuild: `docker compose up --build backend`

---

### Docker Compose — Production

```sh
docker compose -f docker-compose.yml up --build
```

| Service | URL                   |
|---------|-----------------------|
| App     | http://localhost:3400 |

Nginx serves the built frontend and reverse-proxies API and MCP traffic to the backend.

---

## MCP Setup

The MCP server runs at `/mcp` over **streamable HTTP**.

### Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "zdzira": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### Claude Code (CLI)

Add to your project's `.claude/settings.json`:

```json
{
  "mcpServers": {
    "zdzira": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

### Claude Code running inside the Docker dev container

When Claude Code itself runs inside the `backend` container, `localhost` refers to the container — not the host.
Use the **internal Docker network hostname** instead:

```json
{
  "mcpServers": {
    "zdzira": {
      "type": "http",
      "url": "http://backend:8080/mcp"
    }
  }
}
```

### Available MCP Tools

| Tool             | Description                               |
|------------------|-------------------------------------------|
| `list_projects`  | List all projects                         |
| `get_project`    | Get project details                       |
| `create_project` | Create a new project                      |
| `list_epics`     | List all epics in a project               |
| `get_epic`       | Get epic details with linked issues       |
| `create_epic`    | Create an epic                            |
| `update_epic`    | Update an epic                            |
| `list_issues`    | List issues (filterable by type/priority) |
| `get_issue`      | Get a single issue                        |
| `create_issue`   | Create an issue                           |
| `update_issue`   | Update issue fields                       |
| `move_issue`     | Move an issue to a different swimlane     |
| `delete_issue`   | Delete an issue                           |
| `add_comment`    | Add a comment to an issue                 |
| `list_comments`  | List comments on an issue                 |
| `link_issues`    | Link two issues                           |
| `list_links`     | List links for an issue                   |
| `list_swimlanes` | List swimlanes in a project               |

Issues are referenced as `PROJ-42`, epics as `PROJ-E1`.

---

## REST API

Full OpenAPI spec at `http://localhost:8080/docs`. A committed snapshot lives at
[`docs/openapi.json`](docs/openapi.json); regenerate it with `make openapi`.

Base path: `/api/v1`

| Method         | Path                                               | Description             |
|----------------|----------------------------------------------------|-------------------------|
| GET            | `/projects`                                        | List projects           |
| POST           | `/projects`                                        | Create project          |
| GET            | `/projects/{slug}`                                 | Get project             |
| DELETE         | `/projects/{slug}`                                 | Delete project          |
| GET            | `/projects/{slug}/board`                           | Get board view          |
| GET            | `/projects/{slug}/epics`                           | List epics              |
| POST           | `/projects/{slug}/epics`                           | Create epic             |
| GET/PUT/DELETE | `/projects/{slug}/epics/{epicRef}`                 | Epic CRUD               |
| GET            | `/projects/{slug}/issues`                          | List issues             |
| POST           | `/projects/{slug}/issues`                          | Create issue            |
| GET/PUT/DELETE | `/projects/{slug}/issues/{issueRef}`               | Issue CRUD              |
| POST           | `/projects/{slug}/issues/{issueRef}/move`          | Move to swimlane        |
| GET/POST       | `/projects/{slug}/issues/{issueRef}/comments`      | List / add comments     |
| DELETE         | `/projects/{slug}/issues/{issueRef}/comments/{id}` | Delete comment          |
| GET/POST       | `/projects/{slug}/issues/{issueRef}/links`         | List / create links     |
| GET            | `/projects/{slug}/swimlanes`                       | List swimlanes          |
| POST           | `/projects/{slug}/swimlanes`                       | Create swimlane         |
| PATCH          | `/projects/{slug}/swimlanes/{id}`                  | Update swimlane         |
| POST           | `/projects/{slug}/swimlanes/reorder`               | Reorder swimlanes       |
| DELETE         | `/projects/{slug}/swimlanes/{id}`                  | Delete swimlane         |
| GET            | `/projects/{slug}/audit`                           | List audit log          |
| GET            | `/api/v1/events`                                   | SSE board update stream |
| GET            | `/health`                                          | Liveness check          |
| GET            | `/ready`                                           | Readiness check (DB)    |

A **Postman collection** with DEV and Nginx environments is in `docs/`.

---

## Development

```sh
make hooks           # wire up git hooks (run once after cloning)
make install-tools   # install golangci-lint
go test ./...        # run all backend tests
make build           # build the binary to bin/zdzira
make build-frontend  # build the frontend bundle to frontend/dist/
```

Commits follow [Conventional Commits](https://www.conventionalcommits.org). The commit-msg hook enforces this.

### Project structure

```
backend/
  api/       HTTP handlers and router (Chi + Huma)
  mcp/       MCP server (streamable HTTP)
  model/     GORM models
  service/   Business logic
  store/     Database layer (SQLite via GORM)
cmd/zdzira/  Binary entry point
frontend/
  src/
    js/      Vanilla JS (htmx, Handlebars, custom DnD)
    styles/  SASS — neo-brutalist design system
    templates/ Handlebars templates (rendered client-side)
    includes/  Pug partials
docs/        Screenshots, ERD, Postman collection, ADRs
```

---

## Screenshots

![Issue panel](docs/screenshot1b.png)

![Issue details view](docs/screenshot2.png)

![Epic details view](docs/screenshot3.png)

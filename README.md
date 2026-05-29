# Zdzira

A local issue tracker for personal software development. No accounts or auth - access is direct and unguarded for local
development and project management in the AI age: Rest API with a simple UI for human users and an MCP server for agents.
Runs as a single binary on your machine.

## Running

```sh
go build -o bin/zdzira ./cmd/zdzira
./bin/zdzira
```

By default, it listens on `:8080` and stores data in `./zdzira.db` (`SQLite`). Both are configurable:

```sh
./bin/zdzira -addr :9000 -db ~/my-projects.db
```

## MCP

The MCP server runs at `http://localhost:8080/mcp` over streamable HTTP. Add it to your Claude Desktop config:

```json
{
  "mcpServers": {
    "zdzira": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

Available tools: `list_projects`, `get_project`, `create_project`, `list_epics`, `get_epic`, `create_epic`,
`list_issues`, `get_issue`, `create_issue`, `move_issue`, `delete_issue`, `add_comment`, `list_comments`, `link_issues`, `list_links`.

Issues are referenced as `PROJ-42`, epics as `PROJ-E1`.

## REST

The REST API mirrors the same operations under `/projects/{slug}/`. See the router in `internal/api/router.go` for the
full route list.

## Development

```sh
make hooks          # wire up git hooks (run once after cloning)
make install-tools  # install golangci-lint
go test ./...
```

Commits follow [Conventional Commits](https://www.conventionalcommits.org). The commit-msg hook enforces this.

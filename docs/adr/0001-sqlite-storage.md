# SQLite as the storage engine

This is a local, single-user dev tool with no server infrastructure. We chose SQLite over PostgreSQL and use `modernc.org/sqlite` (pure Go, no CGO) so the entire service ships as a single binary with zero external dependencies.

## Considered options

- **PostgreSQL** — better tooling and query capabilities, but requires a running server and adds operational overhead that is unjustified for a local tool
- **CGO-based SQLite (`mattn/go-sqlite3`)** — more mature driver, but requires a C compiler at build time; `modernc.org/sqlite` provides equivalent functionality as pure Go

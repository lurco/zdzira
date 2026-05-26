# Polymorphic Comment with check constraint

Comments can be attached to an Issue, an Epic, or a Project. Rather than three separate tables (`issue_comments`, `epic_comments`, `project_comments`), we use a single `comments` table with three nullable foreign keys and a check constraint enforcing that exactly one is non-null.

This keeps comment logic unified (one GORM model, one set of service methods) at the cost of nullable FKs. The check constraint makes the invariant explicit at the database level.

## Consequences

Every query against `comments` that joins to a parent must handle the three-FK pattern. The check constraint prevents invalid states but is enforced at the application layer for SQLite (SQLite parses but does not enforce check constraints referencing multiple columns in all versions — enforce in the GORM `BeforeCreate`/`BeforeSave` hooks as well).

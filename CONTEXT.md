# Zdzira - a local issue tracker

A local, single-user service for tracking software development work. No authentication. Exposes a REST API and an MCP
server for agent access. The name is Zdzira.

## Language

**Project**:
Top-level container for issues, epics, and swimlanes.
_Avoid_: workspace, board, repository

**Swimlane**:
An ordered column within a project that represents issue status. Projects are seeded with `Backlog → In Progress → Done`
on creation.
_Avoid_: Swimline, column, status, stage

**Issue**:
A discrete unit of work within a project, classified by type (`TASK`, `BUG`, `STORY`) and priority (`LOW`, `HIGH`,
`IMMEDIATE`).
_Avoid_: ticket, card, item, task (when used generically)

**Epic**:
A named grouping of related issues within a project.
_Avoid_: milestone, theme, initiative

**Issue Reference**:
The canonical identifier for an issue: project shortcut + sequential number, e.g. `PROJ-42`. Shortcut is always
uppercase.
_Avoid_: issue ID, issue key, ticket number

**Epic Reference**:
The canonical identifier for an epic: project shortcut + `E` + sequential number, e.g. `PROJ-E1`.
_Avoid_: epic ID, epic key

**Shortcut**:
An uppercase, user-defined string that prefixes all issue and epic references for a project (e.g. `PROJ`, `API`, `FE`).
_Avoid_: prefix, key, abbreviation

**Slug**:
A lowercase, hyphenated string auto-derived from the project name. Used only in URI paths (e.g. `/projects/my-project`).
_Avoid_: handle, URL key, identifier

**Link**:
A directed connection between two issues with a semantic type (`DUPLICATES`, `IS_PART_OF`, `BLOCKS`, `RELATES_TO`). The
source issue is `issue_a`; the target issue is `issue_b`.
_Avoid_: relation, dependency, association

**Comment**:
A text note attached to exactly one of: Issue, Epic, or Project.

## Relationships

- A **Project** contains many **Swimlanes**, **Issues**, and **Epics**
- An **Issue** belongs to exactly one **Project** and exactly one **Swimlane**
- An **Issue** belongs to at most one **Epic**
- An **Epic** belongs to exactly one **Project**
- A **Link** connects exactly two **Issues**: a source (`issue_a`) and a target (`issue_b`)
- A **Comment** belongs to exactly one of: **Issue**, **Epic**, or **Project**

## Example dialogue

> **Agent:** "Move PROJ-42 to In Progress and add a comment saying it's blocked by PROJ-38."
> **Domain expert:** "`PROJ` is the project shortcut, `42` is the issue number. Moving it means changing its swimlane.
> The block relationship should be a `Link` with type `BLOCKS` where `issue_a` is PROJ-38 and `issue_b` is PROJ-42."

## Flagged ambiguities

- "Swimline" appeared in early ERD drafts — resolved: the correct term is **Swimlane**
- "status" is informally used to mean current swimlane — resolved: the canonical term for the column is **Swimlane**; an
  issue's current status is which swimlane it is in

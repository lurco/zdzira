# Entity-relations diagram

```mermaid
classDiagram
    class Project {
        +id: integer, PK, autoincrement
        +name: string, unique, non-nullable
        +slug: string, unique, non-nullable
        +shortcut: string, unique, non-nullable
        +description: text, nullable
        +issue_counter: integer, non-nullable, default 0
        +epic_counter: integer, non-nullable, default 0
        +created_at: timestamp, non-nullable
        +updated_at: timestamp, non-nullable
        +deleted_at: timestamp, nullable
    }

    class Epic {
        +id: integer, PK, autoincrement
        +number: integer, non-nullable
        +name: string, non-nullable
        +description: text, nullable
        +project_id: integer, FK, non-nullable
        +created_at: timestamp, non-nullable
        +updated_at: timestamp, non-nullable
        +deleted_at: timestamp, nullable
    }

    class Issue {
        +id: integer, PK, autoincrement
        +number: integer, non-nullable
        +type: IssueType, non-nullable
        +priority: Priority, non-nullable
        +name: string, non-nullable
        +description: text, nullable
        +project_id: integer, FK, non-nullable
        +epic_id: integer, FK, nullable
        +swimlane_id: integer, FK, non-nullable
        +position: integer, non-nullable
        +created_at: timestamp, non-nullable
        +updated_at: timestamp, non-nullable
        +deleted_at: timestamp, nullable
    }

    class Swimlane {
        +id: integer, PK, autoincrement
        +project_id: integer, FK, non-nullable
        +name: string, non-nullable
        +position: integer, non-nullable
        +deleted_at: timestamp, nullable
    }

    class Link {
        +id: integer, PK, autoincrement
        +type: LinkType, non-nullable
        +issue_a: integer, FK, non-nullable
        +issue_b: integer, FK, non-nullable
    }

    class Comment {
        +id: integer, PK, autoincrement
        +contents: text, non-nullable
        +issue_id: integer, FK, nullable
        +epic_id: integer, FK, nullable
        +project_id: integer, FK, nullable
        +created_at: timestamp, non-nullable
        +updated_at: timestamp, non-nullable
        +deleted_at: timestamp, nullable
    }

    class AuditEntry {
        +id: integer, PK, autoincrement
        +project_id: integer, FK, non-nullable
        +entity_type: string, non-nullable
        +ref: string, non-nullable
        +action: string, non-nullable
        +created_at: timestamp, non-nullable
    }

    class IssueType {
        <<enumeration>>
        TASK
        BUG
        STORY
    }

    class Priority {
        <<enumeration>>
        LOW
        HIGH
        IMMEDIATE
    }

    class LinkType {
        <<enumeration>>
        DUPLICATES
        IS_PART_OF
        BLOCKS
        RELATES_TO
    }

    Project *-- Epic : contains
    Project *-- Issue : contains
    Project *-- Swimlane : contains
    Epic o-- Issue : can collect
    Issue --> Swimlane : has status
    Link --> Issue : connects (a=source)
    Link --> Issue : connects (b=target)
    Issue o-- Comment : on
    Epic o-- Comment : on
    Project o-- Comment : on
    Project *-- AuditEntry : logs
```

## Constraints

- `Comment`: exactly one of `issue_id`, `epic_id`, `project_id` is non-null (check constraint)
- `Link`: `issue_a` = source, `issue_b` = target for directed types (`BLOCKS`, `IS_PART_OF`)
- `Issue.position`: creation order within swimlane, not manually reorderable
- `Swimlane`: seeded with `Backlog → In Progress → Done` on project creation
- Soft delete (`deleted_at`) on all entities except `Link` and `AuditEntry`; project deletion cascades to all owned entities
- `AuditEntry`: append-only, no `updated_at` or `deleted_at`; `entity_type` ∈ {"issue", "epic"}, `action` ∈ {"created", "updated", "moved", "deleted"}
- `AuditEntry` is REST-only — not exposed via MCP (see ADR 0004)

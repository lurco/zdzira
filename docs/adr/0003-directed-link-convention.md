# `issue_a` as source in directed Links

The `Link` table uses generic column names `issue_a` and `issue_b` rather than `source_id` / `target_id`. By convention, `issue_a` is always the source and `issue_b` is always the target for directed link types (`BLOCKS`, `IS_PART_OF`). For symmetric types (`DUPLICATES`, `RELATES_TO`) the order is irrelevant and ignored at query time.

This convention is enforced in the service layer and documented in CONTEXT.md and the ERD. The column names were kept generic to match the original schema design; renaming them would require a migration with no functional benefit.

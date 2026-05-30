# MCP read model is a deliberate agent-facing projection

The MCP tools do not return the same shapes as the REST API or the GORM models. They return purpose-built projections
(`issueSummary`, `issueDetail`, `epicSummary`, `epicDetail`, `commentView`) shaped for an LLM agent's attention and
context budget rather than for database fidelity.

Two agent failure modes drove this:

- **Missing comments.** Comments carry human feedback and instructions. When `get_issue` returned a bare summary, an
  agent could pick up an issue and act on it without ever seeing its comments — nothing signalled they existed.
- **Context overwhelm.** `get_board` and `list_issues` inlined every issue's full description. On a large board this
  floods the context window and buries the signal, degrading the agent's reasoning.

## Decision

Shape MCP payloads around what an agent needs at each step, not around the storage model:

- **Detail calls inline feedback.** `get_issue` returns the description plus the issue's `comments` and `links` in one
  call; `get_epic` inlines its `comments`. An agent acting on an entity cannot miss the feedback attached to it.
- **List/board calls are a scannable index.** `list_issues` and `get_board` omit descriptions and instead carry a
  `comment_count` per issue — a cue that feedback exists and a prompt to drill in with `get_issue`. `list_epics`
  mirrors this.
- **Comment counts are batched.** `CommentStore.CountByIssueIDs` / `CountByEpicIDs` compute counts with a single
  grouped query, so adding the cue does not introduce an N+1 on the board's hot path.
- **Tool descriptions steer the workflow.** Descriptions tell the agent to read comments before acting and to fetch
  full text via the detail call.

## Consequences

- The MCP read shape intentionally diverges from REST and from the GORM models. A reader comparing the two should not
  expect them to match; this is the reason why.
- The contract is load-bearing for agent behaviour (comments inlined on detail, descriptions absent from lists,
  `comment_count` present), so it is guarded by tests in `backend/mcp/tools_test.go`. Changing a projection shape means
  updating those tests deliberately.
- Mutation tools (`create_*`, `update_*`) still return their existing shapes; only the read paths were reshaped. Aligning
  those is possible future work.

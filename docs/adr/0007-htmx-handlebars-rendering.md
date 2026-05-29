# HTMX + client-side Handlebars rendering against the JSON API

All UI interactions are wired declaratively via HTMX attributes in pug
(`hx-get`, `hx-post`, `hx-target`, `hx-swap`, `hx-push-url`). The backend API
remains JSON-only. HTMX consumes JSON via the `client-side-templates` extension,
which feeds the response to a named Handlebars template; HTMX swaps the rendered
HTML into the target.

**Why not server-rendered HTML fragments:** would force a second markup language
in the repo (Go `html/template`) alongside pug, and split the rendering surface
between Go and the frontend. Keeping rendering client-side preserves the
pug-everywhere rule (ADR 0005) and keeps the API a single contract reusable by
MCP and external clients.

**Why not raw `fetch` + manual DOM:** that is what `board.js` does today and is
exactly the cognitive complexity we want gone. HTMX gives declarative wiring
(action, target, swap, history push) directly in markup next to the element
that triggers it, instead of imperative JS coordinating fetches and DOM writes.

**Template authoring:** templates are declared inside pug as
`script(type="text/x-handlebars-template")` blocks under `frontend/src/templates/`,
one partial per renderable shape (card, lane, board, issue-detail, etc.). All
templates are included once into each entry HTML shell so they are available at
page load — no on-demand template fetching.

**JS surface:** restricted to three concerns: HTMX config (`Content-Type:
application/json` on requests, `popstate` re-fetch for back/forward), drag-drop
on the board (fires `htmx.ajax` on drop), and a shared `<dialog>` helper that
renders a template into the modal body and calls `showModal()`. No imperative
DOM-from-JSON code anywhere.

**Trade-off accepted:** Handlebars adds ~20 KB and a build-time dependency. The
client must hold the full template bundle in memory. Both are acceptable for a
small single-tenant tool; if template count grows past ~30 or the bundle past
~100 KB, switch to on-demand template loading per route.

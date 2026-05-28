# HTML entry points with `<pug src>` partials — not `.pug` as Vite entries

Vite entry points are `.html` files (`board.html`, `index.html`, `issue.html`). Pug templates are partials injected via
`<pug src="includes/topbar.pug" />` tags, compiled by `vite-plugin-pug` during the HTML transform step.

**Why not `.pug` as entries directly:** `vite-plugin-pug` does not transform `.pug` files — it only replaces
`<pug src="…">` tags inside HTML files that Vite already knows about. Passing a `.pug` file as a Rollup input causes a
parse error ("Expression expected") because Rollup treats it as JavaScript. The plugin's job is HTML post-processing,
not entry resolution.

**Structure:** Each page has a thin HTML shell (doctype, `<head>`, font links, `<script type="module">`) and one or more
`<pug src="…" />` tags for the body regions. Pug files are pure partials — no `doctype`, `extends` or layout
inheritance as the "layout" is the HTML shell itself.

**Trade-off accepted:** Pug's `extends`/`block` inheritance is unavailable. Shared structure (head metadata, font links)
is duplicated across the three HTML shells. For three pages this is acceptable; if the page count grows significantly, a
pre-compilation step (e.g. `pug-cli` generating HTML that Vite then picks up) would recover inheritance without changing
the runtime output.

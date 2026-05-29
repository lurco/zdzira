# Frontend

The Zdzira web UI: a static, build-time-rendered board with no SPA framework.
See the [root README](../README.md) for running the full stack.

## Stack

- **Vite** multi-page build — two entries: `src/index.html` (project list) and `src/board.html` (the board).
- **Pug** partials (`src/includes`) compose each page at build time; **Handlebars** templates (`src/templates`) render
  dynamic fragments client-side.
- **htmx** for declarative request/swap, plus hand-written JS (`src/js`) for the board, drag-and-drop, filtering, and
  the SSE live-update subscription.
- **SASS** (`src/styles`) — a neo-brutalist design system driven by CSS custom properties (`abstracts/_tokens.sass`).

## Layout

```
src/
  index.html, board.html   page entries (Vite rollup inputs)
  includes/   Pug partials (topbar, toolbar, board shell, modals)
  templates/  Handlebars templates rendered at runtime
  js/         board logic, drag-and-drop, filtering, SSE client
  styles/     SASS — abstracts (tokens/mixins), components, themes
  mixins/     Pug mixins
  public/     static assets copied verbatim (e.g. favicon.svg)
```

## Develop

```sh
npm install
npm run dev     # Vite dev server on :5173 with HMR
npm run build   # production bundle to ../dist (served by nginx in prod)
```

In dev, Vite proxies `/api` and `/mcp` to the backend (`VITE_BACKEND_URL`, default `http://localhost:8080`); the
`/api/v1/events` SSE stream is proxied unbuffered so board updates arrive live. Styling rule: selectors are class-only
(see [`../docs/adr/0006-class-only-sass-selectors.md`](../docs/adr/0006-class-only-sass-selectors.md)).

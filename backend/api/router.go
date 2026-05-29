package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"zdzira/backend/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const docsHTML = `<!doctype html>
<html>
<head>
  <title>Zdzira API</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script id="api-reference" data-url="/api/v1/openapi.json"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

func NewRouter(svcs *service.Services, logger *slog.Logger, ready func(context.Context) error) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(SlogMiddleware(logger))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := ready(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{"status": "unavailable", "error": err.Error()})
			return
		}
		w.Write([]byte(`{"status":"ready"}`))
	})

	config := huma.DefaultConfig("Zdzira API", "1.0.0")
	config.Info.Description = "Local issue tracker for AI-assisted software development. " +
		"REST API for human access; MCP server at /mcp for agent access."
	config.DocsPath = "" // disabled — we serve our own docs below

	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(docsHTML))
	})

	r.Route("/api/v1", func(sub chi.Router) {
		api := humachi.New(sub, config)
		registerProjectRoutes(api, svcs)
		registerEpicRoutes(api, svcs)
		registerIssueRoutes(api, svcs)
		registerSwimlaneRoutes(api, svcs)
		registerBoardRoutes(api, svcs)
		registerCommentRoutes(api, svcs)
		registerLinkRoutes(api, svcs)
		registerAuditRoutes(api, svcs)
	})

	return r
}

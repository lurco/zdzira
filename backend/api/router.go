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

func NewRouter(svcs *service.Services, logger *slog.Logger, ready func(context.Context) error, b *Broadcaster) http.Handler {
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

	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(docsHTML))
	})

	r.Route("/api/v1", func(sub chi.Router) {
		sub.Use(NotifyMiddleware(b))
		sub.Get("/events", eventsHandler(b))
		api := humachi.New(sub, openAPIConfig())
		registerAllRoutes(api, svcs)
	})

	return r
}

// openAPIConfig holds the single source of truth for API metadata, shared by
// the live server and the offline spec exporter.
func openAPIConfig() huma.Config {
	config := huma.DefaultConfig("Zdzira API", "1.0.0")
	config.Info.Description = "Local issue tracker for AI-assisted software development. " +
		"REST API for human access; MCP server at /mcp for agent access."
	config.DocsPath = "" // disabled — we serve our own docs at /docs
	return config
}

// registerAllRoutes wires every resource onto the given API. Kept separate so
// the same registration drives both the running server and OpenAPISpec.
func registerAllRoutes(api huma.API, svcs *service.Services) {
	registerProjectRoutes(api, svcs)
	registerEpicRoutes(api, svcs)
	registerIssueRoutes(api, svcs)
	registerSwimlaneRoutes(api, svcs)
	registerBoardRoutes(api, svcs)
	registerCommentRoutes(api, svcs)
	registerLinkRoutes(api, svcs)
	registerAuditRoutes(api, svcs)
}

// OpenAPISpec builds the OpenAPI document offline (no listening server) and
// returns it as indented JSON, matching what /api/v1/openapi.json serves.
func OpenAPISpec(svcs *service.Services) ([]byte, error) {
	api := humachi.New(chi.NewRouter(), openAPIConfig())
	registerAllRoutes(api, svcs)
	return json.MarshalIndent(api.OpenAPI(), "", "  ")
}

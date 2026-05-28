package api

import (
	"log/slog"
	"net/http"
	"zdzira/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(svcs *service.Services, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(SlogMiddleware(logger))

	config := huma.DefaultConfig("Zdzira API", "1.0.0")
	config.Info.Description = "Local issue tracker for AI-assisted software development. " +
		"REST API for human access; MCP server at /mcp for agent access."

	api := humachi.New(r, config)

	registerProjectRoutes(api, svcs)
	registerEpicRoutes(api, svcs)
	registerIssueRoutes(api, svcs)
	registerSwimlaneRoutes(api, svcs)
	registerCommentRoutes(api, svcs)
	registerLinkRoutes(api, svcs)
	registerAuditRoutes(api, svcs)

	return r
}

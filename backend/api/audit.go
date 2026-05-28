package api

import (
	"context"
	"net/http"
	"zdzira/backend/model"
	"zdzira/backend/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerAuditRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-audit-log",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/audit",
		Summary:     "Audit log for a project (human-only, not available via MCP)",
		Description: "Returns all recorded mutations for issues and epics in reverse-chronological order. Append-only — entries are never modified or deleted.",
		Tags:        []string{"Audit"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{ Body []model.AuditEntry }, error) {
		entries, err := svcs.Audit.ListForProject(ctx, input.Slug)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body []model.AuditEntry }{entries}, nil
	})
}

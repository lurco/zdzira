package api

import (
	"context"
	"net/http"
	"zdzira/backend/model"
	"zdzira/backend/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerLinkRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-links-for-issue",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/issues/{issueRef}/links",
		Summary:     "List all links for an issue (as source or target)",
		Tags:        []string{"Links"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
	}) (*struct{ Body []model.Link }, error) {
		links, err := svcs.Links.ListForIssue(ctx, input.Slug, input.IssueRef)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body []model.Link }{links}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-link",
		Method:        http.MethodPost,
		Path:          "/projects/{slug}/issues/{issueRef}/links",
		Summary:       "Create a directed link from this issue to another",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Links"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Source issue reference"          example:"PROJ-42"`
		Body     struct {
			TargetRef string         `json:"target_ref" doc:"Target issue reference"         example:"PROJ-38"`
			Type      model.LinkType `json:"type"       doc:"Link type"                      example:"BLOCKS"`
		}
	}) (*struct{ Body *model.Link }, error) {
		l, err := svcs.Links.Create(ctx, service.CreateLinkInput{
			ProjectSlug: input.Slug,
			SourceRef:   input.IssueRef,
			TargetRef:   input.Body.TargetRef,
			Type:        input.Body.Type,
		})
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Link }{l}, nil
	})
}

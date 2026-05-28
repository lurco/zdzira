package api

import (
	"context"
	"net/http"
	"zdzira/internal/model"
	"zdzira/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerIssueRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-issues",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/issues",
		Summary:     "List issues in a project, with optional filters",
		Tags:        []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug       string `path:"slug"        doc:"Project slug"       example:"my-project"`
		Type       string `query:"type"       doc:"Filter by type"     example:"BUG"`
		Priority   string `query:"priority"   doc:"Filter by priority" example:"HIGH"`
		SwimlaneID uint   `query:"swimlane_id" doc:"Filter by swimlane ID"`
		EpicID     uint   `query:"epic_id"    doc:"Filter by epic ID"`
	}) (*struct{ Body []model.Issue }, error) {
		f := service.IssueFilterInput{ProjectSlug: input.Slug}
		filtered := false
		if input.Type != "" {
			t := model.IssueType(input.Type)
			f.Type = &t
			filtered = true
		}
		if input.Priority != "" {
			p := model.Priority(input.Priority)
			f.Priority = &p
			filtered = true
		}
		if input.SwimlaneID != 0 {
			f.SwimlaneID = &input.SwimlaneID
			filtered = true
		}
		if input.EpicID != 0 {
			f.EpicID = &input.EpicID
			filtered = true
		}
		if filtered {
			issues, err := svcs.Issues.Filter(ctx, f)
			if err != nil {
				return nil, huma.Error404NotFound(err.Error())
			}
			return &struct{ Body []model.Issue }{issues}, nil
		}
		issues, err := svcs.Issues.List(ctx, input.Slug)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body []model.Issue }{issues}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-issue",
		Method:        http.MethodPost,
		Path:          "/projects/{slug}/issues",
		Summary:       "Create an issue",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
		Body service.CreateIssueInput
	}) (*struct{ Body *model.Issue }, error) {
		input.Body.ProjectSlug = input.Slug
		issue, err := svcs.Issues.Create(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Issue }{issue}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-issue",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/issues/{issueRef}",
		Summary:     "Get an issue by reference",
		Tags:        []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
	}) (*struct{ Body *model.Issue }, error) {
		issue, err := svcs.Issues.Get(ctx, input.Slug, input.IssueRef)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body *model.Issue }{issue}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-issue",
		Method:      http.MethodPut,
		Path:        "/projects/{slug}/issues/{issueRef}",
		Summary:     "Update issue fields",
		Tags:        []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
		Body     service.UpdateIssueInput
	}) (*struct{ Body *model.Issue }, error) {
		input.Body.ProjectSlug = input.Slug
		input.Body.IssueRef = input.IssueRef
		issue, err := svcs.Issues.Update(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Issue }{issue}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "move-issue",
		Method:      http.MethodPost,
		Path:        "/projects/{slug}/issues/{issueRef}/move",
		Summary:     "Move an issue to a different swimlane",
		Tags:        []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
		Body     service.MoveIssueInput
	}) (*struct{ Body *model.Issue }, error) {
		input.Body.ProjectSlug = input.Slug
		input.Body.IssueRef = input.IssueRef
		issue, err := svcs.Issues.Move(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Issue }{issue}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-issue",
		Method:        http.MethodDelete,
		Path:          "/projects/{slug}/issues/{issueRef}",
		Summary:       "Delete an issue",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
	}) (*struct{}, error) {
		if err := svcs.Issues.Delete(ctx, input.Slug, input.IssueRef); err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return nil, nil
	})
}

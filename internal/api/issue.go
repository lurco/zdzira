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
		Summary:     "List all issues in a project",
		Tags:        []string{"Issues"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{ Body []model.Issue }, error) {
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

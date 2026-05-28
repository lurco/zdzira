package api

import (
	"context"
	"net/http"
	"zdzira/internal/model"
	"zdzira/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerCommentRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-comments-for-issue",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/issues/{issueRef}/comments",
		Summary:     "List comments on an issue",
		Tags:        []string{"Comments"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
	}) (*struct{ Body []model.Comment }, error) {
		comments, err := svcs.Comments.ListForIssue(ctx, input.Slug, input.IssueRef)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body []model.Comment }{comments}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "add-comment-to-issue",
		Method:        http.MethodPost,
		Path:          "/projects/{slug}/issues/{issueRef}/comments",
		Summary:       "Add a comment to an issue",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Comments"},
	}, func(ctx context.Context, input *struct {
		Slug     string `path:"slug"     doc:"Project slug"                  example:"my-project"`
		IssueRef string `path:"issueRef" doc:"Issue reference, e.g. PROJ-42"  example:"PROJ-42"`
		Body     struct {
			Contents string `json:"contents" doc:"Comment text" example:"Looks good, merging."`
		}
	}) (*struct{ Body *model.Comment }, error) {
		c, err := svcs.Comments.AddToIssue(ctx, input.Slug, input.IssueRef, input.Body.Contents)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Comment }{c}, nil
	})
}

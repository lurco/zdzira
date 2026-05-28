package api

import (
	"context"
	"net/http"
	"zdzira/backend/model"
	"zdzira/backend/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerProjectRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-projects",
		Method:      http.MethodGet,
		Path:        "/projects",
		Summary:     "List all projects",
		Tags:        []string{"Projects"},
	}, func(ctx context.Context, _ *struct{}) (*struct{ Body []model.Project }, error) {
		projects, err := svcs.Projects.List(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		return &struct{ Body []model.Project }{projects}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-project",
		Method:        http.MethodPost,
		Path:          "/projects",
		Summary:       "Create a project",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Projects"},
	}, func(ctx context.Context, input *struct {
		Body service.CreateProjectInput
	}) (*struct{ Body *model.Project }, error) {
		p, err := svcs.Projects.Create(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Project }{p}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-project",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}",
		Summary:     "Get a project by slug",
		Tags:        []string{"Projects"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{ Body *model.Project }, error) {
		p, err := svcs.Projects.Get(ctx, input.Slug)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body *model.Project }{p}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-project",
		Method:        http.MethodDelete,
		Path:          "/projects/{slug}",
		Summary:       "Delete a project and all its contents",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Projects"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{}, error) {
		if err := svcs.Projects.Delete(ctx, input.Slug); err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return nil, nil
	})
}

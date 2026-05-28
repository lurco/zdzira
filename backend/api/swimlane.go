package api

import (
	"context"
	"net/http"
	"zdzira/backend/model"
	"zdzira/backend/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerSwimlaneRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-swimlanes",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/swimlanes",
		Summary:     "List swimlanes in a project",
		Tags:        []string{"Swimlanes"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{ Body []model.Swimlane }, error) {
		swimlanes, err := svcs.Swimlanes.ListForProject(ctx, input.Slug)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body []model.Swimlane }{swimlanes}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-swimlane",
		Method:        http.MethodPost,
		Path:          "/projects/{slug}/swimlanes",
		Summary:       "Add a swimlane to a project",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Swimlanes"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
		Body service.CreateSwimlaneInput
	}) (*struct{ Body *model.Swimlane }, error) {
		input.Body.ProjectSlug = input.Slug
		sl, err := svcs.Swimlanes.Create(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Swimlane }{sl}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "rename-swimlane",
		Method:      http.MethodPut,
		Path:        "/projects/{slug}/swimlanes/{id}",
		Summary:     "Rename a swimlane",
		Tags:        []string{"Swimlanes"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug"  example:"my-project"`
		ID   uint   `path:"id"   doc:"Swimlane ID"`
		Body service.RenameSwimlaneInput
	}) (*struct{ Body *model.Swimlane }, error) {
		input.Body.ProjectSlug = input.Slug
		input.Body.ID = input.ID
		sl, err := svcs.Swimlanes.Rename(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Swimlane }{sl}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-swimlane",
		Method:        http.MethodDelete,
		Path:          "/projects/{slug}/swimlanes/{id}",
		Summary:       "Delete a swimlane",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Swimlanes"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
		ID   uint   `path:"id"   doc:"Swimlane ID"`
	}) (*struct{}, error) {
		if err := svcs.Swimlanes.Delete(ctx, input.Slug, input.ID); err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return nil, nil
	})
}

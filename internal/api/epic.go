package api

import (
	"context"
	"net/http"
	"zdzira/internal/model"
	"zdzira/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerEpicRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "list-epics",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/epics",
		Summary:     "List all epics in a project",
		Tags:        []string{"Epics"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{ Body []model.Epic }, error) {
		epics, err := svcs.Epics.List(ctx, input.Slug)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body []model.Epic }{epics}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-epic",
		Method:        http.MethodPost,
		Path:          "/projects/{slug}/epics",
		Summary:       "Create an epic",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Epics"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
		Body service.CreateEpicInput
	}) (*struct{ Body *model.Epic }, error) {
		input.Body.ProjectSlug = input.Slug
		e, err := svcs.Epics.Create(ctx, input.Body)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity(err.Error())
		}
		return &struct{ Body *model.Epic }{e}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-epic",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/epics/{epicRef}",
		Summary:     "Get an epic by reference",
		Tags:        []string{"Epics"},
	}, func(ctx context.Context, input *struct {
		Slug    string `path:"slug"    doc:"Project slug"                example:"my-project"`
		EpicRef string `path:"epicRef" doc:"Epic reference, e.g. PROJ-E1" example:"PROJ-E1"`
	}) (*struct{ Body *model.Epic }, error) {
		e, err := svcs.Epics.Get(ctx, input.Slug, input.EpicRef)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body *model.Epic }{e}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-epic",
		Method:        http.MethodDelete,
		Path:          "/projects/{slug}/epics/{epicRef}",
		Summary:       "Delete an epic",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Epics"},
	}, func(ctx context.Context, input *struct {
		Slug    string `path:"slug"    doc:"Project slug"                 example:"my-project"`
		EpicRef string `path:"epicRef" doc:"Epic reference, e.g. PROJ-E1"  example:"PROJ-E1"`
	}) (*struct{}, error) {
		if err := svcs.Epics.Delete(ctx, input.Slug, input.EpicRef); err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return nil, nil
	})
}

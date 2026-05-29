package api

import (
	"context"
	"net/http"
	"zdzira/backend/service"

	"github.com/danielgtaylor/huma/v2"
)

func registerBoardRoutes(api huma.API, svcs *service.Services) {
	huma.Register(api, huma.Operation{
		OperationID: "get-board",
		Method:      http.MethodGet,
		Path:        "/projects/{slug}/board",
		Summary:     "Aggregate board view: swimlanes with grouped issues and the project's epics",
		Tags:        []string{"Board"},
	}, func(ctx context.Context, input *struct {
		Slug string `path:"slug" doc:"Project slug" example:"my-project"`
	}) (*struct{ Body *service.BoardView }, error) {
		view, err := svcs.Board.Get(ctx, input.Slug)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}
		return &struct{ Body *service.BoardView }{view}, nil
	})
}

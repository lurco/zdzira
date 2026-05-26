package api

import (
	"net/http"
	"zdzira/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(svcs *service.Services) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	projects := &projectHandler{svcs.Projects}
	epics := &epicHandler{svcs.Epics}
	issues := &issueHandler{svcs.Issues}
	comments := &commentHandler{svcs.Comments}
	links := &linkHandler{svcs.Links}

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", projects.list)
		r.Post("/", projects.create)
		r.Route("/{slug}", func(r chi.Router) {
			r.Get("/", projects.get)
			r.Delete("/", projects.delete)

			r.Route("/epics", func(r chi.Router) {
				r.Get("/", epics.list)
				r.Post("/", epics.create)
				r.Route("/{epicRef}", func(r chi.Router) {
					r.Get("/", epics.get)
					r.Delete("/", epics.delete)
				})
			})

			r.Route("/issues", func(r chi.Router) {
				r.Get("/", issues.list)
				r.Post("/", issues.create)
				r.Route("/{issueRef}", func(r chi.Router) {
					r.Get("/", issues.get)
					r.Put("/", issues.update)
					r.Delete("/", issues.delete)
					r.Post("/move", issues.move)

					r.Route("/comments", func(r chi.Router) {
						r.Get("/", comments.listForIssue)
						r.Post("/", comments.addToIssue)
					})

					r.Route("/links", func(r chi.Router) {
						r.Get("/", links.listForIssue)
						r.Post("/", links.create)
					})
				})
			})
		})
	})

	return r
}

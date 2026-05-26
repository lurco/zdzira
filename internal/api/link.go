package api

import (
	"encoding/json"
	"net/http"
	"zdzira/internal/model"
	"zdzira/internal/service"

	"github.com/go-chi/chi/v5"
)

type linkHandler struct{ svc *service.LinkService }

func (h *linkHandler) listForIssue(w http.ResponseWriter, r *http.Request) {
	links, err := h.svc.ListForIssue(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "issueRef"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, links)
}

func (h *linkHandler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TargetRef string         `json:"target_ref"`
		Type      model.LinkType `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	l, err := h.svc.Create(r.Context(), service.CreateLinkInput{
		ProjectSlug: chi.URLParam(r, "slug"),
		SourceRef:   chi.URLParam(r, "issueRef"),
		TargetRef:   body.TargetRef,
		Type:        body.Type,
	})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, l)
}

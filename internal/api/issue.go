package api

import (
	"encoding/json"
	"net/http"
	"zdzira/internal/service"

	"github.com/go-chi/chi/v5"
)

type issueHandler struct{ svc *service.IssueService }

func (h *issueHandler) list(w http.ResponseWriter, r *http.Request) {
	issues, err := h.svc.List(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, issues)
}

func (h *issueHandler) get(w http.ResponseWriter, r *http.Request) {
	issue, err := h.svc.Get(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "issueRef"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, issue)
}

func (h *issueHandler) create(w http.ResponseWriter, r *http.Request) {
	var in service.CreateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in.ProjectSlug = chi.URLParam(r, "slug")
	issue, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, issue)
}

func (h *issueHandler) update(w http.ResponseWriter, r *http.Request) {
	var in service.UpdateIssueInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in.ProjectSlug = chi.URLParam(r, "slug")
	in.IssueRef = chi.URLParam(r, "issueRef")
	issue, err := h.svc.Update(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, issue)
}

func (h *issueHandler) move(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Swimlane string `json:"swimlane"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	issue, err := h.svc.Move(r.Context(), service.MoveIssueInput{
		ProjectSlug:  chi.URLParam(r, "slug"),
		IssueRef:     chi.URLParam(r, "issueRef"),
		SwimlaneName: body.Swimlane,
	})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, issue)
}

func (h *issueHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "issueRef")); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

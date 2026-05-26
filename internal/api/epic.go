package api

import (
	"encoding/json"
	"net/http"
	"zdzira/internal/service"

	"github.com/go-chi/chi/v5"
)

type epicHandler struct{ svc *service.EpicService }

func (h *epicHandler) list(w http.ResponseWriter, r *http.Request) {
	epics, err := h.svc.List(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, epics)
}

func (h *epicHandler) get(w http.ResponseWriter, r *http.Request) {
	e, err := h.svc.Get(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "epicRef"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *epicHandler) create(w http.ResponseWriter, r *http.Request) {
	var in service.CreateEpicInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	in.ProjectSlug = chi.URLParam(r, "slug")
	e, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *epicHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "epicRef")); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

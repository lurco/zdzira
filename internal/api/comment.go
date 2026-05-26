package api

import (
	"encoding/json"
	"net/http"
	"zdzira/internal/service"

	"github.com/go-chi/chi/v5"
)

type commentHandler struct{ svc *service.CommentService }

func (h *commentHandler) listForIssue(w http.ResponseWriter, r *http.Request) {
	comments, err := h.svc.ListForIssue(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "issueRef"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *commentHandler) addToIssue(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Contents string `json:"contents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	c, err := h.svc.AddToIssue(r.Context(), chi.URLParam(r, "slug"), chi.URLParam(r, "issueRef"), body.Contents)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

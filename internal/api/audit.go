package api

import (
	"net/http"
	"zdzira/internal/service"

	"github.com/go-chi/chi/v5"
)

type auditHandler struct{ svc *service.AuditService }

func (h *auditHandler) listForProject(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.ListForProject(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

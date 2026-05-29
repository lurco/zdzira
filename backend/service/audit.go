package service

import (
	"context"
	"fmt"
	"zdzira/backend/model"
	"zdzira/backend/store"
)

type AuditService struct {
	stores *store.Stores
}

func (s *AuditService) ListForProject(ctx context.Context, projectSlug string) ([]model.AuditEntry, error) {
	p, err := s.stores.Projects.GetBySlug(ctx, projectSlug)
	if err != nil {
		return nil, fmt.Errorf("project %q not found", projectSlug)
	}
	return s.stores.Audit.ListByProject(ctx, p.ID)
}

// record appends an audit entry. title is the entity's name at the time of the
// action (so the log is readable without a lookup, even after rename/delete);
// detail is an optional summary of the change (e.g. "Backlog → In Progress").
func (s *AuditService) record(ctx context.Context, projectID uint, entityType, ref, action, title, detail string) {
	_ = s.stores.Audit.Record(ctx, &model.AuditEntry{
		ProjectID:  projectID,
		EntityType: entityType,
		Ref:        ref,
		Title:      title,
		Action:     action,
		Detail:     detail,
	})
}

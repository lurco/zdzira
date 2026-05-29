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

// record appends an audit entry. detail is an optional human-readable summary
// of the change (e.g. "Backlog → In Progress"); pass "" when there is nothing
// useful to add beyond the action.
func (s *AuditService) record(ctx context.Context, projectID uint, entityType, ref, action, detail string) {
	_ = s.stores.Audit.Record(ctx, &model.AuditEntry{
		ProjectID:  projectID,
		EntityType: entityType,
		Ref:        ref,
		Action:     action,
		Detail:     detail,
	})
}

package service

import (
	"context"
	"fmt"
	"zdzira/internal/model"
	"zdzira/internal/store"
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

func (s *AuditService) record(ctx context.Context, projectID uint, entityType, ref, action string) {
	_ = s.stores.Audit.Record(ctx, &model.AuditEntry{
		ProjectID:  projectID,
		EntityType: entityType,
		Ref:        ref,
		Action:     action,
	})
}

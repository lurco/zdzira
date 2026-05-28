package store

import (
	"context"
	"zdzira/backend/model"

	"gorm.io/gorm"
)

type gormAuditStore struct{ db *gorm.DB }

func (s *gormAuditStore) Record(ctx context.Context, entry *model.AuditEntry) error {
	return s.db.WithContext(ctx).Create(entry).Error
}

func (s *gormAuditStore) ListByProject(ctx context.Context, projectID uint) ([]model.AuditEntry, error) {
	var entries []model.AuditEntry
	err := s.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&entries).Error
	return entries, err
}

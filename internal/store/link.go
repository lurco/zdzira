package store

import (
	"context"
	"zdzira/internal/model"

	"gorm.io/gorm"
)

type gormLinkStore struct{ db *gorm.DB }

func (s *gormLinkStore) Create(ctx context.Context, l *model.Link) error {
	return s.db.WithContext(ctx).Create(l).Error
}

func (s *gormLinkStore) ListByIssue(ctx context.Context, issueID uint) ([]model.Link, error) {
	var links []model.Link
	err := s.db.WithContext(ctx).
		Where("issue_a = ? OR issue_b = ?", issueID, issueID).
		Find(&links).Error
	return links, err
}

func (s *gormLinkStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Link{}, id).Error
}

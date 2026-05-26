package store

import (
	"context"
	"zdzira/internal/model"

	"gorm.io/gorm"
)

type gormIssueStore struct{ db *gorm.DB }

func (s *gormIssueStore) Create(ctx context.Context, i *model.Issue) error {
	return s.db.WithContext(ctx).Create(i).Error
}

func (s *gormIssueStore) GetByRef(ctx context.Context, projectID uint, number uint) (*model.Issue, error) {
	var i model.Issue
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND number = ?", projectID, number).
		First(&i).Error
	return &i, err
}

func (s *gormIssueStore) ListByProject(ctx context.Context, projectID uint) ([]model.Issue, error) {
	var issues []model.Issue
	err := s.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("position ASC").
		Find(&issues).Error
	return issues, err
}

func (s *gormIssueStore) ListBySwimlane(ctx context.Context, swimlaneID uint) ([]model.Issue, error) {
	var issues []model.Issue
	err := s.db.WithContext(ctx).
		Where("swimlane_id = ?", swimlaneID).
		Order("position ASC").
		Find(&issues).Error
	return issues, err
}

func (s *gormIssueStore) Update(ctx context.Context, i *model.Issue) error {
	return s.db.WithContext(ctx).Save(i).Error
}

func (s *gormIssueStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Issue{}, id).Error
}

func (s *gormIssueStore) DeleteByProject(ctx context.Context, projectID uint) error {
	return s.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&model.Issue{}).Error
}

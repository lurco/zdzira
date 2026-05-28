package store

import (
	"context"
	"zdzira/backend/model"

	"gorm.io/gorm"
)

type gormCommentStore struct{ db *gorm.DB }

func (s *gormCommentStore) Create(ctx context.Context, c *model.Comment) error {
	return s.db.WithContext(ctx).Create(c).Error
}

func (s *gormCommentStore) ListByIssue(ctx context.Context, issueID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := s.db.WithContext(ctx).Where("issue_id = ?", issueID).Find(&comments).Error
	return comments, err
}

func (s *gormCommentStore) ListByEpic(ctx context.Context, epicID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := s.db.WithContext(ctx).Where("epic_id = ?", epicID).Find(&comments).Error
	return comments, err
}

func (s *gormCommentStore) ListByProject(ctx context.Context, projectID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := s.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&comments).Error
	return comments, err
}

func (s *gormCommentStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

func (s *gormCommentStore) DeleteByProject(ctx context.Context, projectID uint) error {
	return s.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&model.Comment{}).Error
}

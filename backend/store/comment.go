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

func (s *gormCommentStore) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var c model.Comment
	err := s.db.WithContext(ctx).First(&c, id).Error
	return &c, err
}

func (s *gormCommentStore) Update(ctx context.Context, c *model.Comment) error {
	return s.db.WithContext(ctx).Save(c).Error
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

// CountByIssueIDs returns a map of issue ID → comment count for the given
// issues in a single grouped query, avoiding an N+1 across board/list views.
// Issues with no comments are absent from the map (callers treat missing as 0).
func (s *gormCommentStore) CountByIssueIDs(ctx context.Context, issueIDs []uint) (map[uint]uint, error) {
	return s.countByParent(ctx, "issue_id", issueIDs)
}

// CountByEpicIDs is the epic-keyed counterpart of CountByIssueIDs.
func (s *gormCommentStore) CountByEpicIDs(ctx context.Context, epicIDs []uint) (map[uint]uint, error) {
	return s.countByParent(ctx, "epic_id", epicIDs)
}

// countByParent groups comment counts by a parent foreign-key column. Soft-deleted
// comments are excluded automatically by GORM's deleted_at scope.
func (s *gormCommentStore) countByParent(ctx context.Context, column string, parentIDs []uint) (map[uint]uint, error) {
	counts := make(map[uint]uint, len(parentIDs))
	if len(parentIDs) == 0 {
		return counts, nil
	}
	type countRow struct {
		ParentID uint
		Total    uint
	}
	var rows []countRow
	err := s.db.WithContext(ctx).
		Model(&model.Comment{}).
		Select(column+" AS parent_id, COUNT(*) AS total").
		Where(column+" IN ?", parentIDs).
		Group(column).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		counts[r.ParentID] = r.Total
	}
	return counts, nil
}

func (s *gormCommentStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

func (s *gormCommentStore) DeleteByProject(ctx context.Context, projectID uint) error {
	return s.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&model.Comment{}).Error
}

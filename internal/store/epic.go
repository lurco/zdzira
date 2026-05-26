package store

import (
	"context"
	"zdzira/internal/model"

	"gorm.io/gorm"
)

type gormEpicStore struct{ db *gorm.DB }

func (s *gormEpicStore) Create(ctx context.Context, e *model.Epic) error {
	return s.db.WithContext(ctx).Create(e).Error
}

func (s *gormEpicStore) GetByRef(ctx context.Context, projectID uint, number uint) (*model.Epic, error) {
	var e model.Epic
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND number = ?", projectID, number).
		First(&e).Error
	return &e, err
}

func (s *gormEpicStore) ListByProject(ctx context.Context, projectID uint) ([]model.Epic, error) {
	var epics []model.Epic
	err := s.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&epics).Error
	return epics, err
}

func (s *gormEpicStore) Update(ctx context.Context, e *model.Epic) error {
	return s.db.WithContext(ctx).Save(e).Error
}

func (s *gormEpicStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Epic{}, id).Error
}

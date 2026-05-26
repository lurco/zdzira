package store

import (
	"context"
	"zdzira/internal/model"

	"gorm.io/gorm"
)

type gormSwimlaneStore struct{ db *gorm.DB }

func (s *gormSwimlaneStore) Create(ctx context.Context, sl *model.Swimlane) error {
	return s.db.WithContext(ctx).Create(sl).Error
}

func (s *gormSwimlaneStore) ListByProject(ctx context.Context, projectID uint) ([]model.Swimlane, error) {
	var swimlanes []model.Swimlane
	err := s.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("position ASC").
		Find(&swimlanes).Error
	return swimlanes, err
}

func (s *gormSwimlaneStore) GetByName(ctx context.Context, projectID uint, name string) (*model.Swimlane, error) {
	var sl model.Swimlane
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND name = ?", projectID, name).
		First(&sl).Error
	return &sl, err
}

func (s *gormSwimlaneStore) Update(ctx context.Context, sl *model.Swimlane) error {
	return s.db.WithContext(ctx).Save(sl).Error
}

func (s *gormSwimlaneStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Swimlane{}, id).Error
}

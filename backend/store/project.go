package store

import (
	"context"
	"zdzira/backend/model"

	"gorm.io/gorm"
)

type gormProjectStore struct{ db *gorm.DB }

func (s *gormProjectStore) Create(ctx context.Context, p *model.Project) error {
	return s.db.WithContext(ctx).Create(p).Error
}

func (s *gormProjectStore) GetBySlug(ctx context.Context, slug string) (*model.Project, error) {
	var p model.Project
	err := s.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	return &p, err
}

func (s *gormProjectStore) List(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	err := s.db.WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (s *gormProjectStore) Update(ctx context.Context, p *model.Project) error {
	return s.db.WithContext(ctx).Save(p).Error
}

func (s *gormProjectStore) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&model.Project{}, id).Error
}

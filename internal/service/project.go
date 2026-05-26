package service

import (
	"context"
	"regexp"
	"strings"
	"zdzira/internal/model"
	"zdzira/internal/store"
)

type ProjectService struct {
	stores *store.Stores
}

type CreateProjectInput struct {
	Name        string
	Shortcut    string
	Description *string
}

func (s *ProjectService) Create(ctx context.Context, in CreateProjectInput) (*model.Project, error) {
	p := &model.Project{
		Name:        in.Name,
		Slug:        slugify(in.Name),
		Shortcut:    strings.ToUpper(in.Shortcut),
		Description: in.Description,
	}
	if err := s.stores.Projects.Create(ctx, p); err != nil {
		return nil, err
	}

	if err := s.seedSwimlanes(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *ProjectService) Get(ctx context.Context, slug string) (*model.Project, error) {
	return s.stores.Projects.GetBySlug(ctx, slug)
}

func (s *ProjectService) List(ctx context.Context) ([]model.Project, error) {
	return s.stores.Projects.List(ctx)
}

func (s *ProjectService) Delete(ctx context.Context, slug string) error {
	p, err := s.stores.Projects.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	return s.stores.Projects.Delete(ctx, p.ID)
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

var defaultSwimlanes = []string{"Backlog", "In Progress", "Done"}

func (s *ProjectService) seedSwimlanes(ctx context.Context, p *model.Project) error {
	for i, name := range defaultSwimlanes {
		sl := &model.Swimlane{
			ProjectID: p.ID,
			Name:      name,
			Position:  uint(i + 1),
		}
		if err := s.stores.Swimlanes.Create(ctx, sl); err != nil {
			return err
		}
	}
	return nil
}

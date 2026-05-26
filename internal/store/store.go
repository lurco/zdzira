package store

import (
	"context"
	"zdzira/internal/model"

	"gorm.io/gorm"
)

type Stores struct {
	Projects  ProjectStore
	Epics     EpicStore
	Issues    IssueStore
	Swimlanes SwimlaneStore
	Comments  CommentStore
	Links     LinkStore
}

func New(db *gorm.DB) *Stores {
	return &Stores{
		Projects:  &gormProjectStore{db},
		Epics:     &gormEpicStore{db},
		Issues:    &gormIssueStore{db},
		Swimlanes: &gormSwimlaneStore{db},
		Comments:  &gormCommentStore{db},
		Links:     &gormLinkStore{db},
	}
}

type ProjectStore interface {
	Create(ctx context.Context, p *model.Project) error
	GetBySlug(ctx context.Context, slug string) (*model.Project, error)
	List(ctx context.Context) ([]model.Project, error)
	Update(ctx context.Context, p *model.Project) error
	Delete(ctx context.Context, id uint) error
}

type EpicStore interface {
	Create(ctx context.Context, e *model.Epic) error
	GetByRef(ctx context.Context, projectID uint, number uint) (*model.Epic, error)
	ListByProject(ctx context.Context, projectID uint) ([]model.Epic, error)
	Update(ctx context.Context, e *model.Epic) error
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type IssueStore interface {
	Create(ctx context.Context, i *model.Issue) error
	GetByRef(ctx context.Context, projectID uint, number uint) (*model.Issue, error)
	ListByProject(ctx context.Context, projectID uint) ([]model.Issue, error)
	ListBySwimlane(ctx context.Context, swimlaneID uint) ([]model.Issue, error)
	Update(ctx context.Context, i *model.Issue) error
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type SwimlaneStore interface {
	Create(ctx context.Context, s *model.Swimlane) error
	ListByProject(ctx context.Context, projectID uint) ([]model.Swimlane, error)
	GetByName(ctx context.Context, projectID uint, name string) (*model.Swimlane, error)
	Update(ctx context.Context, s *model.Swimlane) error
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type CommentStore interface {
	Create(ctx context.Context, c *model.Comment) error
	ListByIssue(ctx context.Context, issueID uint) ([]model.Comment, error)
	ListByEpic(ctx context.Context, epicID uint) ([]model.Comment, error)
	ListByProject(ctx context.Context, projectID uint) ([]model.Comment, error)
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type LinkStore interface {
	Create(ctx context.Context, l *model.Link) error
	ListByIssue(ctx context.Context, issueID uint) ([]model.Link, error)
	Delete(ctx context.Context, id uint) error
}

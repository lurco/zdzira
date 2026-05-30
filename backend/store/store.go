package store

import (
	"context"
	"zdzira/backend/model"

	"gorm.io/gorm"
)

type Stores struct {
	db        *gorm.DB
	Projects  ProjectStore
	Epics     EpicStore
	Issues    IssueStore
	Swimlanes SwimlaneStore
	Comments  CommentStore
	Links     LinkStore
	Audit     AuditStore
}

func New(db *gorm.DB) *Stores {
	return &Stores{
		db:        db,
		Projects:  &gormProjectStore{db},
		Epics:     &gormEpicStore{db},
		Issues:    &gormIssueStore{db},
		Swimlanes: &gormSwimlaneStore{db},
		Comments:  &gormCommentStore{db},
		Links:     &gormLinkStore{db},
		Audit:     &gormAuditStore{db},
	}
}

func (s *Stores) Ping(ctx context.Context) error {
	return s.db.WithContext(ctx).Exec("SELECT 1").Error
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
	GetByID(ctx context.Context, id uint) (*model.Epic, error)
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
	ListFiltered(ctx context.Context, projectID uint, f IssueStoreFilter) ([]model.Issue, error)
	Update(ctx context.Context, i *model.Issue) error
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type IssueStoreFilter struct {
	Type       *model.IssueType
	Priority   *model.Priority
	SwimlaneID *uint
	EpicID     *uint
}

type SwimlaneStore interface {
	Create(ctx context.Context, s *model.Swimlane) error
	ListByProject(ctx context.Context, projectID uint) ([]model.Swimlane, error)
	GetByName(ctx context.Context, projectID uint, name string) (*model.Swimlane, error)
	GetByID(ctx context.Context, id uint) (*model.Swimlane, error)
	Update(ctx context.Context, s *model.Swimlane) error
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type CommentStore interface {
	Create(ctx context.Context, c *model.Comment) error
	GetByID(ctx context.Context, id uint) (*model.Comment, error)
	Update(ctx context.Context, c *model.Comment) error
	ListByIssue(ctx context.Context, issueID uint) ([]model.Comment, error)
	ListByEpic(ctx context.Context, epicID uint) ([]model.Comment, error)
	ListByProject(ctx context.Context, projectID uint) ([]model.Comment, error)
	CountByIssueIDs(ctx context.Context, issueIDs []uint) (map[uint]uint, error)
	CountByEpicIDs(ctx context.Context, epicIDs []uint) (map[uint]uint, error)
	Delete(ctx context.Context, id uint) error
	DeleteByProject(ctx context.Context, projectID uint) error
}

type LinkStore interface {
	Create(ctx context.Context, l *model.Link) error
	ListByIssue(ctx context.Context, issueID uint) ([]model.Link, error)
	Delete(ctx context.Context, id uint) error
}

type AuditStore interface {
	Record(ctx context.Context, entry *model.AuditEntry) error
	ListByProject(ctx context.Context, projectID uint) ([]model.AuditEntry, error)
}

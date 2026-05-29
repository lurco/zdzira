package service

import (
	"zdzira/backend/store"
)

type Services struct {
	Projects  *ProjectService
	Epics     *EpicService
	Issues    *IssueService
	Comments  *CommentService
	Links     *LinkService
	Swimlanes *SwimlaneService
	Board     *BoardService
	Audit     *AuditService
}

func New(stores *store.Stores) *Services {
	audit := &AuditService{stores: stores}
	return &Services{
		Projects:  &ProjectService{stores: stores},
		Epics:     &EpicService{stores: stores, audit: audit},
		Issues:    &IssueService{stores: stores, audit: audit},
		Comments:  &CommentService{stores: stores},
		Links:     &LinkService{stores: stores},
		Swimlanes: &SwimlaneService{stores: stores},
		Board:     &BoardService{stores: stores},
		Audit:     audit,
	}
}

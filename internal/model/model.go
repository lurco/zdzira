package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID           uint   `gorm:"primarykey"`
	Name         string `gorm:"unique;not null"`
	Slug         string `gorm:"unique;not null"`
	Shortcut     string `gorm:"unique;not null"`
	Description  *string
	IssueCounter uint `gorm:"not null;default:0"`
	EpicCounter  uint `gorm:"not null;default:0"`
	Timestamps
	SoftDelete

	Swimlanes []Swimlane `gorm:"foreignKey:ProjectID"`
	Issues    []Issue    `gorm:"foreignKey:ProjectID"`
	Epics     []Epic     `gorm:"foreignKey:ProjectID"`
}

type Swimlane struct {
	ID        uint   `gorm:"primarykey"`
	ProjectID uint   `gorm:"not null;index"`
	Name      string `gorm:"not null"`
	Position  uint   `gorm:"not null"`
	SoftDelete
}

type Epic struct {
	ID          uint   `gorm:"primarykey"`
	Number      uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description *string
	ProjectID   uint `gorm:"not null;index"`
	Timestamps
	SoftDelete
}

type Issue struct {
	ID          uint      `gorm:"primarykey"`
	Number      uint      `gorm:"not null"`
	Type        IssueType `gorm:"not null"`
	Priority    Priority  `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Description *string
	ProjectID   uint  `gorm:"not null;index"`
	EpicID      *uint `gorm:"index"`
	SwimlaneID  uint  `gorm:"not null;index"`
	Position    uint  `gorm:"not null"`
	Timestamps
	SoftDelete
}

// Link has no soft delete — deletions are permanent.
type Link struct {
	ID     uint     `gorm:"primarykey"`
	Type   LinkType `gorm:"not null"`
	IssueA uint     `gorm:"not null;column:issue_a;index"` // source
	IssueB uint     `gorm:"not null;column:issue_b;index"` // target
}

type Comment struct {
	ID        uint   `gorm:"primarykey"`
	Contents  string `gorm:"not null"`
	IssueID   *uint  `gorm:"index"`
	EpicID    *uint  `gorm:"index"`
	ProjectID *uint  `gorm:"index"`
	Timestamps
	SoftDelete
}

// AuditEntry is append-only — no UpdatedAt, no DeletedAt.
type AuditEntry struct {
	ID         uint      `gorm:"primarykey"               json:"id"`
	ProjectID  uint      `gorm:"not null;index"           json:"project_id"`
	EntityType string    `gorm:"not null"                 json:"entity_type"`
	Ref        string    `gorm:"not null"                 json:"ref"`
	Action     string    `gorm:"not null"                 json:"action"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"  json:"created_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error { return c.validateParent() }
func (c *Comment) BeforeSave(tx *gorm.DB) error   { return c.validateParent() }

func (c *Comment) validateParent() error {
	set := 0
	if c.IssueID != nil {
		set++
	}
	if c.EpicID != nil {
		set++
	}
	if c.ProjectID != nil {
		set++
	}
	if set != 1 {
		return errors.New("comment must belong to exactly one of: issue, epic, or project")
	}
	return nil
}

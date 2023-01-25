package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Jus        int
	GroupID    uint
	MemberID   sql.NullInt64
	StartedAt  sql.NullTime
	Deadline   sql.NullTime
	FinishedAt sql.NullTime
	Member     Member
	Group      Group
}

package models

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	GroupID        int
	TelegramUserId int64
	Name           string
	Ordering       int
	Group          Group
}

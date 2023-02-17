package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	TelegramChatId int64
	Name           string
	Round          int
	Duration       int `gorm:"default:7"`
	StartedAt      sql.NullTime
	Deadline       sql.NullTime
	Timezone       string `gorm:"default:Asia/Jakarta"`
	Members        []Member
	Schedules      []Schedule
}

func (g Group) GetLastMemberOrder() (int, error) {
	member := &Member{}

	result := DB.Where("group_id = ?", g.ID).Order("ordering desc").First(&member)
	if result.Error != nil {
		return 0, result.Error
	}

	return member.Ordering, nil
}

func (g Group) GenerateSchedules() []Schedule {
	var s []Schedule
	for i := 1; i < 31; i++ {
		s = append(s, Schedule{
			Jus:     i,
			GroupID: g.ID,
		})
	}

	return s
}

func (g *Group) SetSchedule(t time.Time) error {
	loc, err := time.LoadLocation(g.Timezone)
	if err != nil {
		return err
	}

	sd := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	ed := sd.AddDate(0, 0, g.Duration)

	g.StartedAt = sql.NullTime{Time: sd, Valid: true}
	g.Deadline = sql.NullTime{Time: ed, Valid: true}

	tx := DB.Save(g)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (g *Group) AssignMembersSchedules() error {
	var schedules []Schedule
	members := []Member{}

	err := DB.Where("group_id = ?", g.ID).Order("ordering ASC").Find(&members).Error
	if err != nil {
		return err
	}

	query := DB.Where("group_id = ?", g.ID).Where("started_at IS NULL").Order("jus ASC").Limit(len(members)).Find(&schedules)
	if query.Error != nil {
		return query.Error
	}

	for i, s := range schedules {
		s.MemberID = sql.NullInt64{Int64: int64(members[i].ID), Valid: true}
		s.StartedAt = g.StartedAt
		s.Deadline = g.Deadline
		DB.Save(&s)
	}

	return nil
}

func (g *Group) IncreaseRound() error {
	g.Round += 1
	tx := DB.Save(g)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (g *Group) ResetMemberOrder() error {
	var tx *gorm.DB
	for _, m := range g.Members {
		l := len(g.Members)
		m.Ordering += 1
		if m.Ordering > l {
			m.Ordering -= l
		}
		tx = DB.Save(m)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func (g Group) JobUniqueKey() map[string]interface{} {
	return map[string]interface{}{"group_id": g.ID}
}

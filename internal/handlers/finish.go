package handlers

import (
	"database/sql"
	"time"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func Finish(c telebot.Context) error {
	group := c.Get("group").(*models.Group)
	member := c.Get("member").(*models.Member)
	sc := &models.Schedule{}

	tx := models.DB.Where("group_id = ?", group.ID).Where("member_id = ?", member.ID).Where("finished_at IS NULL").Where("started_at IS NOT NULL").First(sc)
	if tx.RowsAffected == 0 {
		return c.Reply("No active schedules")
	}

	sc.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	models.DB.Save(sc)

	err := c.Reply("Finished")
	if err != nil {
		return err
	}

	remainings := []models.Schedule{}
	models.DB.Where("group_id = ?", group.ID).Where("finished_at IS NULL").Where("started_at IS NOT NULL").Preload("Member").Find(&remainings)
	if len(remainings) == 0 {
		var free int64
		models.DB.Model(&models.Schedule{}).Where("group_id = ?", group.ID).Where("member_id IS NULL").Count(&free)
		if free == 0 {
			err = group.IncreaseRound()
			if err != nil {
				return err
			}

			err = group.ResetMemberOrder()
			if err != nil {
				return err
			}

			err = models.DB.Model(&models.Schedule{}).Where("group_id = ?", group.ID).UpdateColumns(map[string]interface{}{
				"member_id":   nil,
				"started_at":  nil,
				"deadline":    nil,
				"finished_at": nil,
			}).Error
			if err != nil {
				return err
			}
		}

		group.SetSchedule(time.Now().Add(time.Hour * 24))
		group.AssignMembersSchedules()
	}

	return Progress(c)
}

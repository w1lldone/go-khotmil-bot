package handlers

import (
	"fmt"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

type icon string

const (
	finishedIcon   icon = "🕋"
	onProgressIcon icon = "⭐"
)

func Progress(c telebot.Context) error {
	schedules := &[]models.Schedule{}
	group := c.Get("group").(*models.Group)
	models.DB.Where("group_id = ?", group.ID).Preload("Member").Find(schedules)
	text := fmt.Sprintf(`
*Khotmil Quran %s Putaran %d*
	`, group.Name, group.Round)
	text += fmt.Sprintf(`
Periode %s - %s
	`, group.StartedAt.Time.Format("02 Jan 2006"),
		group.Deadline.Time.Format("02 Jan 2006"))

	for _, s := range *schedules {
		text += fmt.Sprintf(`
Juz *%d* %s %s`, s.Jus, getIcon(s), s.Member.Name)
	}

	text += fmt.Sprintf(`
	
%s = Progres membaca
%s = Selesai
	`, onProgressIcon, finishedIcon)

	return c.Send(text, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

func getIcon(s models.Schedule) icon {
	if !s.MemberID.Valid {
		return ""
	}

	if s.FinishedAt.Valid {
		return finishedIcon
	} else {
		return onProgressIcon
	}
}

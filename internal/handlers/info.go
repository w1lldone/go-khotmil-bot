package handlers

import (
	"fmt"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func Info(c telebot.Context) error {
	group := c.Get("group").(*models.Group)

	text := fmt.Sprintf(`
Khotmil %s
Periode %d
	`, group.Name, group.Round)

	for i, m := range group.Members {
		text += fmt.Sprintf(`
%d. %s`, i+1, m.Name)
	}

	return c.Send(text, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

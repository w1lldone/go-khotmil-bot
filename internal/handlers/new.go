package handlers

import (
	"fmt"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func New(c telebot.Context) error {
	var title string

	chatId := c.Chat().ID
	group := &models.Group{}

	if c.Chat().Title == "" {
		title = c.Chat().FirstName
	} else {
		title = c.Chat().Title
	}

	models.DB.Where(models.Group{
		TelegramChatId: chatId,
	}).Attrs(models.Group{
		Name:      title,
		Schedules: group.GenerateSchedules(),
	}).Preload("Members").FirstOrCreate(group)

	return c.Reply(fmt.Sprintf("A Khotmil has been created: %s", group.Name), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}

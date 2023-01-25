package middlewares

import (
	"log"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"golang.org/x/exp/slices"
	"gopkg.in/telebot.v3"
)

func HasGroup(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		group := &models.Group{}
		chatId := c.Chat().ID

		models.DB.Where("telegram_chat_id = ?", chatId).Preload("Members").Find(group)
		if group.ID == 0 {
			return c.Reply("Group not found!")
		}

		c.Set("group", group)

		return next(c)
	}
}

func AdminOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Chat().Type == "private" {
			return next(c)
		}

		admins, err := c.Bot().AdminsOf(c.Chat())
		if err != nil {
			log.Fatalf("failed getting admins %v", err)
		}

		member := c.Sender()
		isAdmin := slices.IndexFunc(admins, func(cm telebot.ChatMember) bool {
			return cm.User.ID == member.ID
		})

		if isAdmin == -1 {
			return c.Reply("Admin only")
		}

		return next(c)
	}
}

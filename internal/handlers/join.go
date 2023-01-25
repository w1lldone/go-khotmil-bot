package handlers

import (
	"log"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func Join(c telebot.Context) error {
	group := c.Get("group").(*models.Group)
	user := &models.Member{
		TelegramUserId: c.Message().Sender.ID,
	}

	lastMemberOrder, err := group.GetLastMemberOrder()
	if err != nil {
		log.Print(err)
	}

	models.DB.Where(models.Member{TelegramUserId: user.TelegramUserId, GroupID: int(group.ID)}).Attrs(models.Member{
		Name:     c.Sender().FirstName,
		Ordering: lastMemberOrder + 1,
	}).FirstOrCreate(user)

	return c.Reply("Joined")
}

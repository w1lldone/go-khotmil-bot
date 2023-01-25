package handlers

import (
	"fmt"
	"time"

	"github.com/w1lldone/go-khotmil-bot/internal/models"
	"gopkg.in/telebot.v3"
)

func Start(c telebot.Context) error {
	group := c.Get("group").(*models.Group)
	fmt.Printf("Started AT: %v", group.StartedAt)
	if group.StartedAt.Valid {
		return c.Send("Already started")
	}

	if len(group.Members) == 0 {
		return c.Send("No member")
	}

	group.Round += 1
	err := group.SetSchedule(time.Now())
	if err != nil {
		return err
	}

	err = group.AssignMembersSchedules()
	if err != nil {
		return err
	}

	return c.Send("Hello")
}
